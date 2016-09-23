package nk

// #include "nuklear.h"
import "C"

func (c Color) R() Byte {
	return Byte(c.r)
}

func (c Color) G() Byte {
	return Byte(c.g)
}

func (c Color) B() Byte {
	return Byte(c.b)
}

func (c Color) A() Byte {
	return Byte(c.a)
}

func (c Color) RGBA() (Byte, Byte, Byte, Byte) {
	return Byte(c.r), Byte(c.g), Byte(c.b), Byte(c.a)
}

func (c Color) RGBAi() (int32, int32, int32, int32) {
	return int32(c.r), int32(c.g), int32(c.b), int32(c.a)
}

func (c *Color) SetR(r Byte) {
	c.r = (C.nk_byte)(r)
}

func (c *Color) SetG(g Byte) {
	c.g = (C.nk_byte)(g)
}

func (c *Color) SetB(b Byte) {
	c.b = (C.nk_byte)(b)
}

func (c *Color) SetA(a Byte) {
	c.a = (C.nk_byte)(a)
}

func (c *Color) SetRGBA(r, g, b, a Byte) {
	c.r = (C.nk_byte)(r)
	c.g = (C.nk_byte)(g)
	c.b = (C.nk_byte)(b)
	c.a = (C.nk_byte)(a)
}

func (c *Color) SetRGBAi(r, g, b, a int32) {
	c.r = (C.nk_byte)(r)
	c.g = (C.nk_byte)(g)
	c.b = (C.nk_byte)(b)
	c.a = (C.nk_byte)(a)
}
