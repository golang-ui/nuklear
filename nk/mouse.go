package nk

import "C"

func (mouse *Mouse) Grab() bool {
	return mouse.grab == True
}

func (mouse *Mouse) Grabbed() bool {
	return mouse.grabbed == True
}

func (mouse *Mouse) Ungrab() bool {
	return mouse.ungrab == True
}

func (mouse *Mouse) ScrollDelta() float32 {
	return (float32)(mouse.scroll_delta)
}

func (mouse *Mouse) Pos() (float32, float32) {
	return (float32)(mouse.pos.x), (float32)(mouse.pos.y)
}

func (mouse *Mouse) SetPos(x, y float32) {
	mouse.pos.x = (C.float)(x)
	mouse.pos.y = (C.float)(y)
}

func (mouse *Mouse) Prev() (float32, float32) {
	return (float32)(mouse.prev.x), (float32)(mouse.prev.y)
}

func (mouse *Mouse) Delta() (float32, float32) {
	return (float32)(mouse.delta.x), (float32)(mouse.delta.y)
}
