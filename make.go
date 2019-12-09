package fifochan

func Make() FIFOChan {
	return MakeBuf(BufSize)
}

func MakeBuf(buf int) FIFOChan {
	return make(FIFOChan, buf)
}
