package nk

import "unsafe"

var VertexLayoutEnd = DrawVertexLayoutElement{
	Attribute: VertexAttributeCount,
	Format:    FormatCount,
	Offset:    0,
}

func NkDrawForeach(ctx *Context, b *Buffer, fn func(cmd *DrawCommand)) {
	cmd := Nk_DrawBegin(ctx, b)
	for {
		if cmd == nil {
			break
		}
		fn(cmd)
		cmd = Nk_DrawNext(cmd, b, ctx)
	}
}

func (h Handle) ID() int {
	return int(*(*int64)(unsafe.Pointer(&h)))
}

func (h Handle) Ptr() uintptr {
	return *(*uintptr)(unsafe.Pointer(&h))
}
