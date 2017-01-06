package nk

func (cmd *DrawCommand) ElemCount() int {
	return (int)(cmd.elem_count)
}

func (cmd *DrawCommand) Texture() Handle {
	return (Handle)(cmd.texture)
}

func (cmd *DrawCommand) ClipRect() *Rect {
	return (*Rect)(&cmd.clip_rect)
}
