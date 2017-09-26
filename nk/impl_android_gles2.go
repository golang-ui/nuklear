// +build android
// +build gles2
// +build !sdl2

package nk

/*
#cgo CFLAGS: -DNK_INCLUDE_FIXED_TYPES -DNK_INCLUDE_STANDARD_IO -DNK_INCLUDE_DEFAULT_ALLOCATOR -DNK_INCLUDE_FONT_BAKING -DNK_INCLUDE_DEFAULT_FONT -DNK_INCLUDE_VERTEX_BUFFER_OUTPUT -Wno-implicit-function-declaration
#cgo LDFLAGS: -lEGL -lGLESv2
#include <string.h>

#define NK_IMPLEMENTATION
#define NK_ANDROID_GLES2_IMPLEMENTATION

#include "nuklear.h"

#include <GLES2/gl2.h>
#include <GLES2/gl2ext.h>
#include <EGL/egl.h>

typedef void (*nkglGenVertexArrays)(GLsizei, GLuint*);
typedef void (*nkglBindVertexArray)(GLuint);
typedef void (*nkglDeleteVertexArrays)(GLsizei, const GLuint*);

typedef void (*nkglGenVertexArraysOES)(GLsizei, GLuint*);
typedef void (*nkglBindVertexArrayOES)(GLuint);
typedef void (*nkglDeleteVertexArraysOES)(GLsizei, const GLuint*);

static nkglGenVertexArrays glGenVertexArrays;
static nkglBindVertexArray glBindVertexArray;
static nkglDeleteVertexArrays glDeleteVertexArrays;

static nkglGenVertexArraysOES glGenVertexArraysOES;
static nkglBindVertexArrayOES glBindVertexArrayOES;
static nkglDeleteVertexArraysOES glDeleteVertexArraysOES;

static void callBindVertexArray(GLuint arg) {
	glBindVertexArray(arg);
}

static void callGenVertexArrays(GLsizei arg1, GLuint* arg2) {
	glGenVertexArrays(arg1, arg2);
}

#define GL_EXT(name) (nk##name)nk_gl_ext(#name)

#define _CHECK_GL_IMPL(line)                        \
    {                                               \
        int ret;                                    \
        while ((ret = glGetError()) != GL_NO_ERROR) \
        {                                           \
            printf("%d: glGetError(): %#x\n",       \
                (line), ret);                       \
        }                                           \
    }

#define CHECK_GL _CHECK_GL_IMPL(__LINE__)

typedef void ( *__GLXextFuncPtr)(void);

static __GLXextFuncPtr nk_gl_ext(const char *name);

static
__GLXextFuncPtr nk_gl_ext(const char *name)
{
    __GLXextFuncPtr func = eglGetProcAddress(name);
    if (!func) {
        fprintf(stdout, "[GL]: failed to load extension: %s\n", name);
        return NULL;
    }
    return func;
}

static int init_ext() {
	glGenVertexArrays = GL_EXT(glGenVertexArraysOES);
	if (!glGenVertexArrays) {
		return -1;
	}
    glBindVertexArray = GL_EXT(glBindVertexArrayOES);
    if (!glBindVertexArray) {
		return -1;
	}
    glDeleteVertexArrays = GL_EXT(glDeleteVertexArraysOES);
    if (!glDeleteVertexArrays) {
		return -1;
	}
    return 0;
}
*/
import "C"
import (
	"log"
	"unsafe"

	"github.com/xlab/android-go/android"
	"github.com/xlab/android-go/egl"
	gl "github.com/xlab/android-go/gles2"
)

type PlatformInitOption int

const (
	PlatformDefault PlatformInitOption = iota
	PlatformInstallCallbacks
)

func NkPlatformInit(win *android.NativeWindow, opt PlatformInitOption) *Context {
	display, err := egl.NewDisplayHandle(win, map[int32]int32{
		egl.SurfaceType:          egl.WindowBit,
		egl.ContextClientVersion: 2.0, // OpenGL ES 2.0

		egl.RedSize:   8,
		egl.GreenSize: 8,
		egl.BlueSize:  8,
		egl.AlphaSize: 8,
		egl.DepthSize: 24,
	})
	if err != nil {
		log.Println("Platform init failed:", err)
		return nil
	}
	if (int)(C.init_ext()) != 0 {
		log.Println("GL ES 2.0 extensions load failed:", err)
		return nil
	}

	state = newPlatformState()
	state.display = display
	state.ctx = NewContext()
	NkInitDefault(state.ctx, nil)
	deviceCreate()
	return state.ctx
}

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

	var (
		last_prog = make([]int32, 1)
		last_tex  = make([]int32, 1)

		last_ebo = make([]int32, 1)
		last_vbo = make([]int32, 1)
		last_vao = make([]int32, 1)
	)

	// save previous GL state
	gl.GetIntegerv(gl.CURRENT_PROGRAM, last_prog)
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, last_tex)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, last_vbo)
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, last_ebo)
	gl.GetIntegerv(C.GL_VERTEX_ARRAY_BINDING_OES, last_vao)

	// setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
	gl.Enable(gl.SCISSOR_TEST)
	gl.ActiveTexture(gl.TEXTURE0)

	gl.Uniform1i(dev.uniform_tex, 0)
	gl.UniformMatrix4fv(dev.uniform_proj, 1, gl.FALSE, &ortho[0][0])
	C.callBindVertexArray((C.GLuint)(dev.vao[0]))

	// convert from command queue into draw list and draw to screen
	{
		var vertices [1024 * 1024]byte
		var elements [1024 * 1024]byte
		pVertices := unsafe.Pointer(&vertices)
		pElements := unsafe.Pointer(&elements)

		//  setup buffers to load vertices and elements
		vbuf := NewBuffer()
		ebuf := NewBuffer()
		NkBufferInitFixed(vbuf, pVertices, Size(maxVertexBuffer))
		NkBufferInitFixed(ebuf, pElements, Size(maxElementBuffer))

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
		NkConvert(state.ctx, dev.cmds, vbuf, ebuf, config)
		// vbuf.Free()
		// ebuf.Free()
		// config.Free()

		gl.BindBuffer(gl.ARRAY_BUFFER, dev.vbo[0])
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, vbuf.Allocated(), pVertices)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, dev.ebo[0])
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, ebuf.Allocated(), pElements)

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
	gl.BindTexture(gl.TEXTURE_2D, uint32(last_tex[0]))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(last_vbo[0]))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(last_ebo[0]))
	C.callBindVertexArray((C.GLuint)(last_vao[0]))
	gl.Disable(gl.SCISSOR_TEST)

	state.display.SwapBuffers()
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

	// Assign each attribute a name
	var vertex_index uint32
	gl.BindAttribLocation(dev.prog, vertex_index, "Position\x00")
	dev.attrib_pos = vertex_index
	vertex_index++
	gl.BindAttribLocation(dev.prog, vertex_index, "TexCoord\x00")
	dev.attrib_uv = vertex_index
	vertex_index++
	gl.BindAttribLocation(dev.prog, vertex_index, "Color\x00")
	dev.attrib_col = vertex_index
	vertex_index++

	gl.LinkProgram(dev.prog)
	gl.GetProgramiv(dev.prog, gl.LINK_STATUS, &status)
	if status != gl.TRUE {
		panic("gl program failed to link")
	}
	gl.UseProgram(dev.prog)
	dev.uniform_tex = gl.GetUniformLocation(dev.prog, "Texture\x00")
	dev.uniform_proj = gl.GetUniformLocation(dev.prog, "ProjMtx\x00")

	{
		// buffer setup
		vs := int32(platformVertexSize)
		vp := unsafe.Offsetof(emptyVertex.position)
		vt := unsafe.Offsetof(emptyVertex.uv)
		vc := unsafe.Offsetof(emptyVertex.col)

		dev.vao = make([]uint32, 1)
		C.callGenVertexArrays(1, (*C.GLuint)(unsafe.Pointer(&dev.vao[0])))
		C.callBindVertexArray((C.GLuint)(dev.vao[0]))

		dev.vbo = make([]uint32, 1)
		gl.GenBuffers(1, dev.vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, dev.vbo[0])
		gl.BufferData(gl.ARRAY_BUFFER, 1024*1024, nil, gl.DYNAMIC_DRAW)

		dev.ebo = make([]uint32, 1)
		gl.GenBuffers(1, dev.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, dev.ebo[0])
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 1024*1024, nil, gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(dev.attrib_pos)
		gl.EnableVertexAttribArray(dev.attrib_uv)
		gl.EnableVertexAttribArray(dev.attrib_col)

		gl.VertexAttribPointer(dev.attrib_pos, 2, gl.FLOAT, gl.FALSE, vs, unsafe.Pointer(vp))
		gl.VertexAttribPointer(dev.attrib_uv, 2, gl.FLOAT, gl.FALSE, vs, unsafe.Pointer(vt))
		gl.VertexAttribPointer(dev.attrib_col, 4, gl.UNSIGNED_BYTE, gl.TRUE, vs, unsafe.Pointer(vc))
	}

	C.callBindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func deviceUploadAtlas(image unsafe.Pointer, width, height int32) {
	dev := state.ogl
	dev.font_tex = make([]uint32, 1)
	gl.GenTextures(1, dev.font_tex)
	gl.BindTexture(gl.TEXTURE_2D, dev.font_tex[0])
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
	gl.DeleteTextures(1, dev.font_tex)
	gl.DeleteBuffers(1, dev.vbo)
	gl.DeleteBuffers(1, dev.ebo)
	NkBufferFree(dev.cmds)
}

func assignShader(shaderHandle uint32, shaderSource string) {
	gl.ShaderSource(shaderHandle, 1, []string{shaderSource + "\x00"}, nil)
	gl.CompileShader(shaderHandle)
}

var vertexShader = `
uniform mat4 ProjMtx;
attribute vec2 Position;
attribute vec2 TexCoord;
attribute vec4 Color;
varying vec2 Frag_UV;
varying vec4 Frag_Color;

void main() {
   Frag_UV = TexCoord;
   Frag_Color = Color;
   gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}`

var fragmentShader = `
precision mediump float;
uniform sampler2D Texture;
varying vec2 Frag_UV;
varying vec4 Frag_Color;

void main() {
   gl_FragColor = Frag_Color * texture2D(Texture, Frag_UV.st);
}`
