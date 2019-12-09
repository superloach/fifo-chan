package fifochan

import (
	"bufio"
	"encoding/json"
	"os"
)

func makeSubscriber(ch *FIFOChan, path string) error {
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
						Err <- err
					}

					err = os.Remove(path)
					if err != nil {
						Err <- err
					}

					return
				default:
					Err <- err
					continue
				}
			}

			var obj interface{}
			data = s.Bytes()
			err := json.Unmarshal(data, &obj)
			if err != nil {
				Err <- err
			} else {
				*ch <- obj
			}
		}
	}()

	return nil
}
