// +build !android
// +build !opengl2
// +build !sdl2

package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT -Wno-implicit-function-declaration
#cgo windows LDFLAGS: -Wl,--allow-multiple-definition
#include <string.h>

#define NK_IMPLEMENTATION
#define NK_GLFW_GL3_IMPLEMENTATION

#include "nuklear.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
)

func NkPlatformRender(aa AntiAliasing, maxVertexBuffer, maxElementBuffer int) {
	dev := state.ogl
	ortho := [4][4]float32{
		{2.0, 0.0, 0.0, 0.0},
		{0.0, -2.0, 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	ortho[0][0] /= float32(state.width)
	ortho[1][1] /= float32(state.height)

	// setup global state
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.ActiveTexture(gl.TEXTURE0)

	// setup program
	gl.UseProgram(dev.prog)
	gl.Uniform1i(dev.uniform_tex, 0)
	gl.UniformMatrix4fv(dev.uniform_proj, 1, false, &ortho[0][0])
	gl.Viewport(0, 0, int32(state.display_width), int32(state.display_height))

	// convert from command queue into draw list and draw to screen
	{
		// allocate vertex and element buffer
		gl.BindVertexArray(dev.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, dev.vbo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, dev.ebo)

		gl.BufferData(gl.ARRAY_BUFFER, maxVertexBuffer, nil, gl.STREAM_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, maxElementBuffer, nil, gl.STREAM_DRAW)

		// load draw vertices & elements directly into vertex + element buffer
		vertices := gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY)
		elements := gl.MapBuffer(gl.ELEMENT_ARRAY_BUFFER, gl.WRITE_ONLY)
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

		//  setup buffers to load vertices and elements
		vbuf := NewBuffer()
		ebuf := NewBuffer()
		NkBufferInitFixed(vbuf, vertices, Size(maxVertexBuffer))
		NkBufferInitFixed(ebuf, elements, Size(maxElementBuffer))
		NkConvert(state.ctx, dev.cmds, vbuf, ebuf, config)
		// vbuf.Free()
		// ebuf.Free()
		// config.Free()

		gl.UnmapBuffer(gl.ARRAY_BUFFER)
		gl.UnmapBuffer(gl.ELEMENT_ARRAY_BUFFER)

		var offset uintptr

		// iterate over and execute each draw command
		NkDrawForeach(state.ctx, dev.cmds, func(cmd *DrawCommand) {
			elemCount := cmd.ElemCount()
			if elemCount == 0 {
				return
			}
			clipRect := cmd.ClipRect()
			gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.Texture().ID()))
			gl.Scissor(
				int32(clipRect.X()*state.fbScaleX),
				int32(float32(state.height-int(clipRect.Y()+clipRect.H()))*state.fbScaleY),
				int32(clipRect.W()*state.fbScaleX),
				int32(clipRect.H()*state.fbScaleY),
			)
			gl.DrawElements(gl.TRIANGLES, int32(elemCount), gl.UNSIGNED_SHORT, unsafe.Pointer(offset))
			offset += uintptr(elemCount) * sizeofDrawIndex
		})

		NkClear(state.ctx)
	}

	// default GL state
	gl.UseProgram(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.Disable(gl.BLEND)
	gl.Disable(gl.SCISSOR_TEST)
}

func deviceCreate() {
	dev := state.ogl
	dev.cmds = NewBuffer()
	NkBufferInitDefault(dev.cmds)
	dev.prog = gl.CreateProgram()

	dev.vert_shdr = gl.CreateShader(gl.VERTEX_SHADER)
	dev.frag_shdr = gl.CreateShader(gl.FRAGMENT_SHADER)
	assignShader(dev.vert_shdr, vertexShader)
	assignShader(dev.frag_shdr, fragmentShader)

	var status int32
	gl.GetShaderiv(dev.vert_shdr, gl.COMPILE_STATUS, &status)
	if status != gl.TRUE {
		panic("vert_shdr failed to compile")
	}
	gl.GetShaderiv(dev.frag_shdr, gl.COMPILE_STATUS, &status)
	if status != gl.TRUE {
		panic("frag_shdr failed to compile")
	}
	gl.AttachShader(dev.prog, dev.vert_shdr)
	gl.AttachShader(dev.prog, dev.frag_shdr)
	gl.LinkProgram(dev.prog)
	gl.GetProgramiv(dev.prog, gl.LINK_STATUS, &status)
	if status != gl.TRUE {
		panic("gl program failed to link")
	}
	dev.uniform_tex = gl.GetUniformLocation(dev.prog, gl.Str("Texture\x00"))
	dev.uniform_proj = gl.GetUniformLocation(dev.prog, gl.Str("ProjMtx\x00"))
	dev.attrib_pos = uint32(gl.GetAttribLocation(dev.prog, gl.Str("Position\x00")))
	dev.attrib_uv = uint32(gl.GetAttribLocation(dev.prog, gl.Str("TexCoord\x00")))
	dev.attrib_col = uint32(gl.GetAttribLocation(dev.prog, gl.Str("Color\x00")))

	{
		// buffer setup
		vs := int32(platformVertexSize)
		vp := unsafe.Offsetof(emptyVertex.position)
		vt := unsafe.Offsetof(emptyVertex.uv)
		vc := unsafe.Offsetof(emptyVertex.col)
		gl.GenBuffers(1, &dev.vbo)
		gl.GenBuffers(1, &dev.ebo)
		gl.GenVertexArrays(1, &dev.vao)

		gl.BindVertexArray(dev.vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, dev.vbo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, dev.ebo)

		gl.EnableVertexAttribArray(dev.attrib_pos)
		gl.EnableVertexAttribArray(dev.attrib_uv)
		gl.EnableVertexAttribArray(dev.attrib_col)

		gl.VertexAttribPointer(dev.attrib_pos, 2, gl.FLOAT, false, vs, unsafe.Pointer(vp))
		gl.VertexAttribPointer(dev.attrib_uv, 2, gl.FLOAT, false, vs, unsafe.Pointer(vt))
		gl.VertexAttribPointer(dev.attrib_col, 4, gl.UNSIGNED_BYTE, true, vs, unsafe.Pointer(vc))
	}

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
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
	gl.DetachShader(dev.prog, dev.vert_shdr)
	gl.DetachShader(dev.prog, dev.frag_shdr)
	gl.DeleteShader(dev.vert_shdr)
	gl.DeleteShader(dev.frag_shdr)
	gl.DeleteProgram(dev.prog)
	gl.DeleteTextures(1, &dev.font_tex)
	gl.DeleteBuffers(1, &dev.vbo)
	gl.DeleteBuffers(1, &dev.ebo)
	NkBufferFree(dev.cmds)
}

func assignShader(shaderHandle uint32, shaderSource string) {
	var header = "#version 300 es"
	if runtime.GOOS == "darwin" {
		header = "#version 150"
	}
	shader := fmt.Sprintf("%s\n%s\x00", header, shaderSource)
	shaderData, free := gl.Strs(shader)
	gl.ShaderSource(shaderHandle, 1, shaderData, nil)
	gl.CompileShader(shaderHandle)
	free()
}

var vertexShader = `
uniform mat4 ProjMtx;
in vec2 Position;
in vec2 TexCoord;
in vec4 Color;
out vec2 Frag_UV;
out vec4 Frag_Color;

void main() {
   Frag_UV = TexCoord;
   Frag_Color = Color;
   gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}`

var fragmentShader = `
precision mediump float;
uniform sampler2D Texture;
in vec2 Frag_UV;
in vec4 Frag_Color;
out vec4 Out_Color;

void main(){
   Out_Color = Frag_Color * texture(Texture, Frag_UV.st);
}`

var state = &platformState{
	ogl: &platformDevice{},
}

type platformDevice struct {
	cmds *Buffer
	null DrawNullTexture

	vbo, vao, ebo uint32
	prog          uint32
	vert_shdr     uint32
	frag_shdr     uint32

	attrib_pos   uint32
	attrib_uv    uint32
	attrib_col   uint32
	uniform_tex  int32
	uniform_proj int32

	font_tex uint32
}
