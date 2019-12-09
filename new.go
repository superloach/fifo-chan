package fifochan

func New(path string) FIFOChan {
	return NewBuf(path, BufSize)
}

func NewBuf(path string, buf int) FIFOChan {
	mutex.Lock()
	defer mutex.Unlock()

	ch, ok := chans[path]
	if ok {
		return ch
	}

	ch = MakeBuf(buf)
	chans[path] = ch

	return ch
}
