package main

import (
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
)

const (
	winWidth  = 800
	winHeight = 550

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		closer.Fatalln(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(winWidth, winHeight, "Nuklear Demo", nil, nil)
	if err != nil {
		closer.Fatalln(err)
	}
	win.MakeContextCurrent()

	width, height := win.GetSize()
	log.Printf("glfw: created window %dx%d", width, height)

	if err := gl.Init(); err != nil {
		closer.Fatalln("opengl: init failed:", err)
	}
	gl.Viewport(0, 0, int32(width), int32(height))

	ctx := nk.NkPlatformInit(win, nk.PlatformInstallCallbacks)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)
	//Use default font so we dont have to manage assets
	sansFont := nk.NkFontAtlasAddDefault(atlas, 14, nil)
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
			glfw.Terminate()
			fpsTicker.Stop()
			close(doneC)
			return
		case <-fpsTicker.C:
			if win.ShouldClose() {
				close(exitC)
				continue
			}
			glfw.PollEvents()
			gfxMain(win, ctx, state)
		}
	}
}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	// Layout
	bounds := nk.NkRect(10, 10, 700, 500)
	//Begin Window
	if nk.NkBegin(ctx, "Wiget Layout Demo", bounds, nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle) > 0 {
		//create dynamic Row with height 30 and 1 column
		nk.NkLayoutRowDynamic(ctx, 30, 1)
		nk.NkLabel(ctx, "Dynamic fixed column layout with generated position and size (resizable)", nk.TextLeft)
		//create dynamic Row with height 20 and 3 columns
		nk.NkLayoutRowDynamic(ctx, 30, 3)
		nk.NkButtonLabel(ctx, "Button 1")
		nk.NkButtonLabel(ctx, "Button 2")
		nk.NkButtonLabel(ctx, "Button 3")

		//create static Row with height 30, width 100 and 1 column
		nk.NkLayoutRowStatic(ctx, 30, 400, 1)
		nk.NkLabel(ctx, "static fixed column layout with generated position and size (not resizable)", nk.TextLeft)
		//create static Row with height 30, width 100 and 3 columns
		nk.NkLayoutRowStatic(ctx, 30, 100, 3)
		nk.NkButtonLabel(ctx, "Button 1")
		nk.NkButtonLabel(ctx, "Button 2")
		nk.NkButtonLabel(ctx, "Button 3")

		//create dynamic Row with height 30 and 1 column
		nk.NkLayoutRowDynamic(ctx, 30, 1)
		nk.NkLabel(ctx, "Dynamic array-based custom column layout with generated position and custom size (resizable)", nk.TextLeft)
		//Define size of widgets in percent
		ratio := []float32{0.1, 0.2, 0.7}
		//Define dynamic layout row with height 30, 3 columns and sizing information
		//WARNING: You have to use constant nk.Dynamic and NOT nk.LayoutDynamic for param LayoutFormat
		nk.NkLayoutRow(ctx, nk.Dynamic, 30, 3, ratio)
		nk.NkButtonLabel(ctx, "10 %")
		nk.NkButtonLabel(ctx, "20 %")
		nk.NkButtonLabel(ctx, "70 %")

		//create dynamic Row with height 30 and 1 column
		nk.NkLayoutRowDynamic(ctx, 30, 1)
		nk.NkLabel(ctx, "Static array-based custom column layout with generated position and custom size (not resizable)", nk.TextLeft)
		//Define total size of widgets
		ratio = []float32{50, 100, 200}
		//Define static layout row with height 30, 3 columns and sizing information
		//WARNING: You have to use constant nk.Static and NOT nk.LayoutStatic for param LayoutFormat
		nk.NkLayoutRow(ctx, nk.Static, 30, 3, ratio)
		nk.NkButtonLabel(ctx, "size 50")
		nk.NkButtonLabel(ctx, "size 100")
		nk.NkButtonLabel(ctx, "size 200")

		//create dynamic Row with height 30 and 1 column
		nk.NkLayoutRowDynamic(ctx, 30, 1)
		nk.NkLabel(ctx, "Dynamic immediate mode custom column layout with generated position and custom size (resizable)", nk.TextLeft)
		//Begin dynamic layout row with height 30 and 4 columns
		//WARNING: You have to use constant nk.Dynamic and NOT nk.LayoutDynamic for param LayoutFormat
		nk.NkLayoutRowBegin(ctx, nk.Dynamic, 30, 4)
		//Define, that all following widgets will have 10% of total width
		nk.NkLayoutRowPush(ctx, 0.1)
		nk.NkButtonLabel(ctx, "10 %")
		nk.NkButtonLabel(ctx, "10 %")
		//Define, that all following widgets will have 30% of total width
		nk.NkLayoutRowPush(ctx, 0.3)
		nk.NkButtonLabel(ctx, "30 %")
		//Define, that all following widgets will have 40% of total width
		nk.NkLayoutRowPush(ctx, 0.4)
		nk.NkButtonLabel(ctx, "40 %")
		nk.NkLayoutRowEnd(ctx)

		//create dynamic Row with height 30 and 1 column
		nk.NkLayoutRowDynamic(ctx, 30, 1)
		nk.NkLabel(ctx, "Dynamic immediate mode custom column layout with generated position and custom size (not resizable)", nk.TextLeft)
		//Begin static layout row with height 30 and 3 columns
		//WARNING: You have to use constant nk.Static and NOT nk.LayoutStatic for param LayoutFormat
		nk.NkLayoutRowBegin(ctx, nk.Static, 30, 3)
		//Define, that all following widgets will have width=50
		nk.NkLayoutRowPush(ctx, 50)
		nk.NkButtonLabel(ctx, "size 50")
		nk.NkButtonLabel(ctx, "size 50")
		//Define, that all following widgets will have width=100
		nk.NkLayoutRowPush(ctx, 100)
		nk.NkButtonLabel(ctx, "size 100")
		nk.NkLayoutRowEnd(ctx)

		//create dynamic Row with height 30 and 1 column
		nk.NkLayoutResetMinRowHeight(ctx)
		nk.NkLayoutRowDynamic(ctx, 0, 1)

		nk.NkLabel(ctx, "Static free space with custom position and custom size (not resizable)", nk.TextLeft)
		//Begin custom layout space with height 100 and 4 widgets
		nk.NkLayoutSpaceBegin(ctx, nk.Static, 100, 4)
		//Define a rectangle at with w=60 and h=60.
		nk.NkLayoutSpacePush(ctx, nk.NkRect(0, 0, 60, 60))
		//This button will be placed in the previously define rectangle
		nk.NkButtonLabel(ctx, "60x60")
		//Define a rectangle at with w=60 and h=30 at pos x=70 and y=10.
		nk.NkLayoutSpacePush(ctx, nk.NkRect(70, 70, 60, 30))
		//This button will be placed in the previously define rectangle
		nk.NkButtonLabel(ctx, "60x30")
		//Define a rectangle at with w=30 and h=60 at pos x=150 and y=45.
		nk.NkLayoutSpacePush(ctx, nk.NkRect(150, 45, 60, 30))
		//This button will be placed in the previously define rectangle
		nk.NkButtonLabel(ctx, "30x60")
		//Do not forget to end layout space (only needed when nk.NkLayoutSpaceBegin was called)
		nk.NkLayoutSpaceEnd(ctx)
	}
	//End Window
	nk.NkEnd(ctx)

	// Render
	bg := make([]float32, 4)
	nk.NkColorFv(bg, state.bgColor)
	width, height := win.GetSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.ClearColor(bg[0], bg[1], bg[2], bg[3])
	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
	win.SwapBuffers()
}

type State struct {
	bgColor nk.Color
}
