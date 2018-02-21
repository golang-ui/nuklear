// +build opengl2
// +build sdl2

package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT -Wno-implicit-function-declaration
#cgo windows LDFLAGS: -Wl,--allow-multiple-definition
#include <string.h>

#define NK_IMPLEMENTATION
#define NK_SDL2_GL2_IMPLEMENTATION

#include "nuklear.h"
*/
import "C"
import (
	"unsafe"

	"github.com/go-gl/gl/v2.1/gl"
)

func NkPlatformRender(aa AntiAliasing, maxVertexBuffer, maxElementBuffer int) {
	dev := state.ogl

	// setup global state
	gl.PushAttrib(gl.ENABLE_BIT | gl.COLOR_BUFFER_BIT | gl.TRANSFORM_BIT)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.TEXTURE_2D)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// setup viewport/project
	gl.Viewport(0, 0, int32(state.display_width), int32(state.display_height))
	gl.MatrixMode(gl.PROJECTION)
	gl.PushMatrix()
	gl.LoadIdentity()
	gl.Ortho(0.0, float64(state.width), float64(state.height), 0.0, -1.0, 1.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PushMatrix()
	gl.LoadIdentity()

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	{
		// convert from command queue into draw list and draw to screen

		vs := int32(platformVertexSize)
		vp := unsafe.Offsetof(emptyVertex.position)
		vt := unsafe.Offsetof(emptyVertex.uv)
		vc := unsafe.Offsetof(emptyVertex.col)

		// fill convert configuration
		config := &ConvertConfig{
			VertexLayout: []DrawVertexLayoutElement{
				{
					Attribute: VertexPosition,
					Format:    FormatFloat,
					Offset:    Size(unsafe.Offsetof(emptyVertex.position)),
				}, {
					Attribute: VertexTexcoord,
					Format:    FormatFloat,
					Offset:    Size(unsafe.Offsetof(emptyVertex.uv)),
				}, {
					Attribute: VertexColor,
					Format:    FormatR8g8b8a8,
					Offset:    Size(unsafe.Offsetof(emptyVertex.col)),
				}, VertexLayoutEnd,
			},
			VertexSize:      Size(platformVertexSize),
			VertexAlignment: Size(platformVertexAlign),
			Null:            dev.null,

			CircleSegmentCount: 22,
			CurveSegmentCount:  22,
			ArcSegmentCount:    22,

			GlobalAlpha: 1.0,
			ShapeAa:     aa,
			LineAa:      aa,
		}

		// convert shapes into vertexes
		vbuf := NewBuffer()
		ebuf := NewBuffer()
		NkBufferInitDefault(vbuf)
		NkBufferInitDefault(ebuf)
		NkConvert(state.ctx, dev.cmds, vbuf, ebuf, config)
		// vbuf.Free()
		// ebuf.Free()
		// config.Free()

		// setup vertex buffer pointer
		vertices := uintptr(unsafe.Pointer(NkBufferMemoryConst(vbuf)))
		gl.VertexPointer(2, gl.FLOAT, vs, unsafe.Pointer(vertices+vp))
		gl.TexCoordPointer(2, gl.FLOAT, vs, unsafe.Pointer(vertices+vt))
		gl.ColorPointer(4, gl.UNSIGNED_BYTE, vs, unsafe.Pointer(vertices+vc))

		offset := uintptr(unsafe.Pointer(NkBufferMemoryConst(ebuf)))

		NkDrawForeach(state.ctx, dev.cmds, func(cmd *DrawCommand) {
			elemCount := cmd.ElemCount()
			if elemCount == 0 {
				return
			}
			clipRect := cmd.ClipRect()
			gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.Texture().ID()))
			gl.Scissor(
				int32(clipRect.X()*state.fbScaleX),
				int32(float32(state.height-int32(clipRect.Y()+clipRect.H()))*state.fbScaleY),
				int32(clipRect.W()*state.fbScaleX),
				int32(clipRect.H()*state.fbScaleY),
			)
			gl.DrawElements(gl.TRIANGLES, int32(elemCount), gl.UNSIGNED_SHORT, unsafe.Pointer(offset))
			offset += uintptr(elemCount)
		})

		NkClear(state.ctx)
		NkBufferFree(vbuf)
		NkBufferFree(ebuf)
	}

	// default OpenGL state
	gl.DisableClientState(gl.VERTEX_ARRAY)
	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.DisableClientState(gl.COLOR_ARRAY)

	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.SCISSOR_TEST)
	gl.Disable(gl.BLEND)
	gl.Disable(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.PopMatrix()
	gl.MatrixMode(gl.PROJECTION)
	gl.PopMatrix()
	gl.PopAttrib()
}

func deviceCreate() {
	dev := state.ogl
	dev.cmds = NewBuffer()
	NkBufferInitDefault(state.ogl.cmds)
}

func deviceUploadAtlas(image unsafe.Pointer, width, height int32) {
	dev := state.ogl
	gl.GenTextures(1, &dev.font_tex)
	gl.BindTexture(gl.TEXTURE_2D, dev.font_tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, image)
}

func deviceDestroy() {
	dev := state.ogl
	gl.DeleteTextures(1, &dev.font_tex)
	NkBufferFree(dev.cmds)
}

var state = &platformState{
	ogl: &platformDevice{},
}

type platformDevice struct {
	cmds *Buffer
	null DrawNullTexture

	font_tex uint32
}
