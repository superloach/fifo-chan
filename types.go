package fifochan

import "os"

type FIFOChan = chan interface{}
type node = *os.File
