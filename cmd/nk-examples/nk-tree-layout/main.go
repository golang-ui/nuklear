package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/golang-ui/nuklear/nk"
	"github.com/xlab/closer"
)

const (
	winWidth  = 500
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
	bounds := nk.NkRect(10, 10, 400, 500)
	//Begin Window
	if nk.NkBegin(ctx, "Group Layout Demo", bounds, nk.WindowBorder|nk.WindowMovable|nk.WindowScalable|nk.WindowMinimizable|nk.WindowTitle) > 0 {

		//Define tree. In this example we need to save the collapsed state ourselves in state.rootTreeState
		if nk.NkTreeStatePush(ctx, nk.TreeNode, "Tree (unmanaged / resizable)", &state.rootTreeState) > 0 {
			nk.NkButtonLabel(ctx, "some button")
			nk.NkButtonLabel(ctx, "some button")
			nk.NkButtonLabel(ctx, "some button")
			//Do not forget to end tree
			nk.NkTreePop(ctx)
		}

		//Define tree. In this example we need to save the collapsed state ourselves in state.rootTreeState2
		if nk.NkTreeStatePush(ctx, nk.TreeNode, "Tree (unmanaged / not resizable)", &state.rootTreeState2) > 0 {
			//This will only be fired when tree is not collapsed
			fmt.Print(".")
			//add static layout row to set width of widgets
			nk.NkLayoutRowStatic(ctx, 15, 100, 1)
			nk.NkButtonLabel(ctx, "some button")
			nk.NkButtonLabel(ctx, "some button")
			nk.NkButtonLabel(ctx, "some button")
			//Do not forget to end tree
			nk.NkTreePop(ctx)
		}

		//Define tree with managed collapse state. We need to provide a unique string, which will be used by nuklear to calculate id
		treeUniqueName := "tree_1"
		if nk.NkTreePushHashed(ctx, nk.TreeNode, "Tree (managed / resizable)", nk.Minimized, treeUniqueName, int32(len(treeUniqueName)), 0) > 0 {
			nk.NkLabel(ctx, "another label", nk.TextLeft)
			nk.NkLabel(ctx, "another label", nk.TextLeft)
			nk.NkLabel(ctx, "another label", nk.TextLeft)
			nk.NkTreePop(ctx)
		}

		//Define tree with managed collapse state. We need to provide a unique string, which will be used by nuklear to calculate id
		treeUniqueName = "tree_2"
		if nk.NkTreePushHashed(ctx, nk.TreeNode, "Tree (managed / not resizable)", nk.Minimized, treeUniqueName, int32(len(treeUniqueName)), 0) > 0 {
			//add static layout row to set width of widgets
			nk.NkLayoutRowStatic(ctx, 20, 100, 1)
			nk.NkLabel(ctx, "another label", nk.TextLeft)
			nk.NkLabel(ctx, "another label", nk.TextLeft)
			nk.NkLabel(ctx, "another label", nk.TextLeft)
			nk.NkTreePop(ctx)
		}

		//Define nested tree with managed collapse state. We need to provide a unique string, which will be used by nuklear to calculate id
		treeUniqueName = "tree_3"
		if nk.NkTreePushHashed(ctx, nk.TreeNode, "Nested Tree", nk.Minimized, treeUniqueName, int32(len(treeUniqueName)), 0) > 0 {
			//add some widgets
			nk.NkLabel(ctx, "some label", nk.TextLeft)
			nk.NkButtonLabel(ctx, "some button")

			//add a nested tree
			treeUniqueName := "tree_4"
			if nk.NkTreePushHashed(ctx, nk.TreeNode, "Tree", nk.Minimized, treeUniqueName, int32(len(treeUniqueName)), 0) > 0 {
				nk.NkLabel(ctx, "another label", nk.TextLeft)
				nk.NkLabel(ctx, "another label", nk.TextLeft)
				nk.NkLabel(ctx, "another label", nk.TextLeft)
				nk.NkTreePop(ctx)
			}
			nk.NkTreePop(ctx)
		}

		//This is what happens, when given string is not unique
		treeUniqueName = "tree_1"
		if nk.NkTreePushHashed(ctx, nk.TreeNode, "ID of this three is not unique", nk.Minimized, treeUniqueName, int32(len(treeUniqueName)), 0) > 0 {
			//add static layout row to set width of widgets
			nk.NkLayoutRowStatic(ctx, 20, 100, 1)
			nk.NkButtonLabel(ctx, "some button")
			nk.NkButtonLabel(ctx, "some button")
			nk.NkButtonLabel(ctx, "some button")
			nk.NkTreePop(ctx)
		}
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
	bgColor        nk.Color
	rootTreeState  nk.CollapseStates
	rootTreeState2 nk.CollapseStates
}
