package fifochan

func Stop() {
	close(stopSubChan)
	close(stopPubChan)
	wg.Wait()
	close(Done)
}
