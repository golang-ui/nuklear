package main

import (
	"log"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/golang-ui/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
)

const (
	winWidth  = 400
	winHeight = 500

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func init() {
	runtime.LockOSThread()
}

func main() {
	glfw.SetErrorCallback(onError)
	if ok := b(glfw.Init()); !ok {
		closer.Fatalln("glfw: init failed")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	glfw.WindowHint(glfw.OpenglForwardCompat, glfw.True)
	win := glfw.CreateWindow(winWidth, winHeight, "Nuklear Demo\x00", nil, nil)
	if win == nil {
		closer.Fatalln("glfw: window creation failed")
	}
	glfw.MakeContextCurrent(win)

	var width, height int32
	glfw.GetWindowSize(win, &width, &height)
	log.Printf("glfw: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, width, height)

	ctx := nk.NkGLFW3Init((*nk.GLFWwindow)(unsafe.Pointer(win)), nk.GLFW3InstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkGLFW3FontStashBegin(&atlas)
	sansFont := nk.NkFontAtlasAddFromFile(atlas, s("assets/FreeSans.ttf"), 16, nil)
	nk.NkGLFW3FontStashEnd()
	if sansFont != nil {
		nk.NkStyleSetFont(ctx, sansFont.Handle())
	}

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	state := &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}
	fpsTicker := time.NewTicker(time.Second / 30)
	for {
		select {
		case <-exitC:
			nk.NkGLFW3Shutdown()
			glfw.Terminate()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			if b(glfw.WindowShouldClose(win)) {
				close(exitC)
				continue
			}
			glfw.PollEvents()
			gfxMain(win, ctx, state)
		}
	}
}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkGLFW3NewFrame()

	// Layout
	panel := nk.NewPanel()
	bounds := nk.NkRect(50, 50, 230, 250)
	update := nk.NkBegin(ctx, panel, s("Demo"), bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update > 0 {
		nk.NkLayoutRowStatic(ctx, 30, 80, 1)
		if nk.NkButtonLabel(ctx, s("button")) > 0 {
			log.Println("[INFO] button pressed!")
		}
		nk.NkLayoutRowDynamic(ctx, 30, 2)
		if nk.NkOptionLabel(ctx, s("easy"), flag(state.opt == Easy)) > 0 {
			state.opt = Easy
		}
		if nk.NkOptionLabel(ctx, s("hard"), flag(state.opt == Hard)) > 0 {
			state.opt = Hard
		}
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		nk.NkPropertyInt(ctx, s("Compression:"), 0, &state.prop, 100, 10, 1)
		{
			combo := nk.NewPanel()
			nk.NkLayoutRowDynamic(ctx, 20, 1)
			nk.NkLabel(ctx, s("background:"), nk.TextLeft)
			nk.NkLayoutRowDynamic(ctx, 25, 1)
			size := nk.NkVec2(nk.NkWidgetWidth(ctx), 400)
			if nk.NkComboBeginColor(ctx, combo, state.bgColor, size) > 0 {
				nk.NkLayoutRowDynamic(ctx, 120, 1)
				state.bgColor = nk.NkColorPicker(ctx, state.bgColor, nk.ColorFormatRGBA)
				nk.NkLayoutRowDynamic(ctx, 25, 1)
				r, g, b, a := state.bgColor.RGBAi()
				r = nk.NkPropertyi(ctx, s("#R:"), 0, r, 255, 1, 1)
				g = nk.NkPropertyi(ctx, s("#G:"), 0, g, 255, 1, 1)
				b = nk.NkPropertyi(ctx, s("#B:"), 0, b, 255, 1, 1)
				a = nk.NkPropertyi(ctx, s("#A:"), 0, a, 255, 1, 1)
				state.bgColor.SetRGBAi(r, g, b, a)
				nk.NkComboEnd(ctx)
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := []float32{38, 38, 38, 255}
	nk.NkColorFv(bg, state.bgColor)
	var width, height int32
	glfw.GetWindowSize(win, &width, &height)
	gl.Viewport(0, 0, width, height)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkGLFW3Render(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	glfw.SwapBuffers(win)
}

type Option uint8

const (
	Easy Option = 0
	Hard Option = 1
)

type State struct {
	bgColor nk.Color
	prop    int32
	opt     Option
}

func onError(code int32, msg string) {
	log.Printf("[glfw ERR]: error %d: %s", code, msg)
}
