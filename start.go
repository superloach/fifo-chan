package fifochan

func Start() {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	Done = make(chan struct{})
	stopPubChan = make(chan struct{})
	stopSubChan = make(chan struct{})

	for path, ch := range chans {
		err = makePublisher(&ch, path)
		if err != nil {
			Err <- err
		}

		err = makeSubscriber(&ch, path)
		if err != nil {
			Err <- err
		}
	}
}
