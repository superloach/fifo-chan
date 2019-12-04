package fifochan

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"syscall"
)

var chans map[string]chan interface{}
var mutex sync.Mutex
var wg sync.WaitGroup
var errChan chan error
var stopChan chan struct{}
var stopPubChan chan struct{}
var stopSubChan chan struct{}

func init() {
	chans = make(map[string]chan interface{})
	errChan = make(chan error, 10)
	stopChan = make(chan struct{})
	stopPubChan = make(chan struct{})
	stopSubChan = make(chan struct{})
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func newNode(path string) (*os.File, error) {
	fmt.Printf("newNode %s\n", path)
	var err error

	if !fileExists(path) {
		err = syscall.Mkfifo(path, 0666)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
}

func New(path string) chan interface{} {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	ch, ok := chans[path]
	if ok {
		return ch
	}

	ch = make(chan interface{}, 1024)

	err = makePublisher(&ch, path)
	if err != nil {
		errChan <- err
		return nil
	}

	err = makeSubscriber(&ch, path)
	if err != nil {
		errChan <- err
		return nil
	}

	chans[path] = ch

	return ch
}

func makePublisher(ch *chan interface{}, path string) error {
	var n *os.File

	n, err := newNode(path)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopPubChan:
				if len(*ch) != 0 {
					continue
				}

				err := n.Close()
				if err != nil {
					errChan <- err
				}

				err = os.Remove(path)
				if err != nil {
					errChan <- err
				}

				return
			case obj := <-*ch:
				data, err := json.Marshal(obj)
				data = append(data, '\n')
				if err != nil {
					errChan <- err
				} else {
					_, err := n.Write(data)
					if err != nil {
						errChan <- err
					}
				}
			}
		}
	}()

	return nil
}

func makeSubscriber(ch *chan interface{}, path string) error {
	n, err := newNode(path)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(n)

	go func() {
		for {
			var data []byte

			if !s.Scan() {
				err = s.Err()
				select {
				case <-stopSubChan:
					err := n.Close()
					if err != nil {
						errChan <- err
					}

					err = os.Remove(path)
					if err != nil {
						errChan <- err
					}

					return
				default:
					errChan <- err
					continue
				}
			}

			var obj interface{}
			data = s.Bytes()
			err := json.Unmarshal(data, &obj)
			if err != nil {
				errChan <- err
			} else {
				*ch <- obj
			}
		}
	}()

	return nil
}

func ErrChan() <-chan error {
	return errChan
}

func Stop() {
	close(stopSubChan)
	close(stopPubChan)
	wg.Wait()
	close(stopChan)
}

func Done() chan struct{} {
	return stopChan
}
