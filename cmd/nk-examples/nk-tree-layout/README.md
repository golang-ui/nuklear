## nk-tree-layout

This example shows, how to use tree layout provided by nuklear API.

Please note, that the following explanations are taken from the nuklear doc, available [here](https://github.com/vurtun/nuklear/tree/master/doc)

## Trees
Trees represent two different concept. First the concept of a collapsable UI section that can be either in a hidden or visibile state. They allow the UI user to selectively minimize the current set of visible UI to comprehend. The second concept are tree widgets for visual UI representation of trees.

Trees thereby can be nested for tree representations and multiple nested collapsible UI sections. All trees are started by calling one of the `nk.NkTreeXXXPushXXX` functions and ended by calling the `nk.NkTreePop` function. Each starting functions takes a title label and optionally an image to be displayed and the initial collapse state of type `nk.CollapseStates`.

Each starting function will either return false(0) if the tree is collapsed or hidden and therefore does not need to be filled with content or true(1) if visible and required to be filled.

The runtime state of the tree is either stored outside the library by the caller or inside which requires a unique ID. If do not want to manage the collapse state yourself, use `nk.NkTreePushHashed` and provide a unique string in param `hash`

The tree header does not require any layouting function and instead calculates a auto height based on the currently used font size

The tree ending functions only need to be called if the tree content is actually visible. So make sure the tree push function is guarded by if and the pop call is only taken if the tree is visible. 
## Install
Make sure, that the following packages are installed
  - github.com/go-gl/gl/v3.2-core/gl
  - github.com/go-gl/glfw/v3.2/glfw
  - github.com/golang-ui/nuklear/nk
  - github.com/xlab/closer

Then run
```
$ go run main.go
```
Now a runnable binary should have been created in your $GOBIN path

### Maintainer
jannst <mkawaganga@gmail.com>
