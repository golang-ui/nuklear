package nk

// #include "nuklear.h"
import "C"
import "unsafe"

func (f Font) Handle() *UserFont {
	return NewUserFontRef(unsafe.Pointer(&f.handle))
}
