package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/golang-ui/nuklear/nk"
	"github.com/veandco/go-sdl2/sdl"
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
	var err error
	sdl.Init(sdl.INIT_EVERYTHING)

	win, err := sdl.CreateWindow("Nuklear Demo", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		closer.Fatalln(err)
	}
	defer win.Destroy()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	context, err := win.GLCreateContext()
	if err != nil {
		closer.Fatalln(err)
	}

	width, height := win.GetSize()
	log.Printf("SDL2: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, context, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	sansFont := nk.NkFontAtlasAddFromBytes(atlas, MustAsset("assets/FreeSans.ttf"), 16, nil)
	nk.NkFontStashEnd()
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
			nk.NkPlatformShutdown()
			sdl.Quit()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					close(exitC)
				}
			}
			gfxMain(win, ctx, state)
		}
	}
}

func gfxMain(win *sdl.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(50, 50, 230, 250)
	update := nk.NkBegin(ctx, "Demo", bounds,
		nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle)

	if update > 0 {
		nk.NkLayoutRowStatic(ctx, 30, 80, 1)
		{
			if nk.NkButtonLabel(ctx, "button") > 0 {
				log.Println("[INFO] button pressed!")
			}
		}
		nk.NkLayoutRowDynamic(ctx, 30, 2)
		{
			if nk.NkOptionLabel(ctx, "easy", flag(state.opt == Easy)) > 0 {
				state.opt = Easy
			}
			if nk.NkOptionLabel(ctx, "hard", flag(state.opt == Hard)) > 0 {
				state.opt = Hard
			}
		}
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		{
			nk.NkPropertyInt(ctx, "Compression:", 0, &state.prop, 100, 10, 1)
		}
		nk.NkLayoutRowDynamic(ctx, 20, 1)
		{
			nk.NkLabel(ctx, "background:", nk.TextLeft)
		}
		nk.NkLayoutRowDynamic(ctx, 25, 1)
		{
			size := nk.NkVec2(nk.NkWidgetWidth(ctx), 400)
			if nk.NkComboBeginColor(ctx, state.bgColor, size) > 0 {
				nk.NkLayoutRowDynamic(ctx, 120, 1)
				cf := nk.NkColorCf(state.bgColor)
				cf = nk.NkColorPicker(ctx, cf, nk.ColorFormatRGBA)
				state.bgColor = nk.NkRgbCf(cf)
				nk.NkLayoutRowDynamic(ctx, 25, 1)
				r, g, b, a := state.bgColor.RGBAi()
				r = nk.NkPropertyi(ctx, "#R:", 0, r, 255, 1, 1)
				g = nk.NkPropertyi(ctx, "#G:", 0, g, 255, 1, 1)
				b = nk.NkPropertyi(ctx, "#B:", 0, b, 255, 1, 1)
				a = nk.NkPropertyi(ctx, "#A:", 0, a, 255, 1, 1)
				state.bgColor.SetRGBAi(r, g, b, a)
				nk.NkComboEnd(ctx)
			}
		}
	}
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.GLSwap()
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
	log.Printf("[ERR]: error %d: %s", code, msg)
}
