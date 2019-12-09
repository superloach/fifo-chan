package fifochan

import (
	"os"
	"syscall"
)

func nodeExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func newNode(path string) (node, error) {
	var err error

	if !nodeExists(path) {
		err = syscall.Mkfifo(path, 0666)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
}
