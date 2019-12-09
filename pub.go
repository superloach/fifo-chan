package fifochan

import (
	"encoding/json"
	"os"
)

func makePublisher(ch *FIFOChan, path string) error {
	var n node

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
					Err <- err
				}

				err = os.Remove(path)
				if err != nil {
					Err <- err
				}

				return
			case obj := <-*ch:
				data, err := json.Marshal(obj)
				data = append(data, '\n')
				if err != nil {
					Err <- err
				} else {
					_, err := n.Write(data)
					if err != nil {
						Err <- err
					}
				}
			}
		}
	}()

	return nil
}
