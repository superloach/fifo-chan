package fifochan

func init() {
	chans = make(map[string]FIFOChan)
	Err = make(chan error, 10)
}
