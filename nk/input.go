package nk

func (input *Input) Mouse() *Mouse {
	return (*Mouse)(&input.mouse)
}
