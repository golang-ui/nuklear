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
		bgColor:     nk.NkRgba(28, 48, 62, 255),
		groupStates: make([]int32, 16),
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
	if nk.NkBegin(ctx, "Group Layout Demo", bounds, nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle) > 0 {

		nk.NkLayoutRowStatic(ctx, 300, 200, 2)
		//create group with border and title
		if nk.NkGroupBegin(ctx, "Group 1", nk.WindowBorder|nk.WindowTitle) > 0 {
			//create row with dynamic calculated height and 1 column
			nk.NkLayoutRowDynamic(ctx, 0, 1)
			//add 16 widgets
			for i := 0; i < 16; i++ {
				nk.NkSelectableLabel(ctx, getGroupStatus(state.groupStates[i]), nk.TextCentered, &state.groupStates[i])
			}
			//Do not forget to end group
			nk.NkGroupEnd(ctx)
		}

		//lets create a floating group. Therefore we need to use layout space api
		nk.NkLayoutSpaceBegin(ctx, nk.Static, 150, 10000)
		//Define the rectangle, in which our group will be placed
		nk.NkLayoutSpacePush(ctx, nk.NkRect(150, 10, 300, 100))
		//try to add the group
		if nk.NkGroupBegin(ctx, "Group 2", nk.WindowBorder|nk.WindowTitle) > 0 {
			//create row with dynamic calculated height and 1 column
			nk.NkLayoutRowDynamic(ctx, 0, 1)
			//add 16 widgets
			for i := 0; i < 16; i++ {
				//add selector widget
				nk.NkSelectableLabel(ctx, getGroupStatus(state.groupStates[i]), nk.TextCentered, &state.groupStates[i])
			}
			//Do not forget to end group
			nk.NkGroupEnd(ctx)
		}
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

func getGroupStatus(b int32) string {
	if b != 0 {
		return "selected"
	} else {
		return "not selected"
	}
}

type State struct {
	bgColor     nk.Color
	groupStates []int32
}
