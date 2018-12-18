## nk-group-layout

This example shows, how to use group layout provided by nuklear API.

Please note, that the following explanations are taken from the nuklear doc, available [here](https://github.com/vurtun/nuklear/tree/master/doc)

## Groups
Groups are basically windows inside windows. They allow to subdivide space in a window to layout widgets as a group. Almost all more complex widget layouting requirements can be solved using groups and basic layouting fuctionality. Groups just like windows are identified by an unique name and internally keep track of scrollbar offsets by default. However additional versions are provided to directly manage the scrollbar. 

To create a group you have to call one of the three `nk.NkGroupBegin` functions to start group declarations and `nk.NkGroupEnd` at the end. Furthermore it is required to check the return value of `nk.NkGroupBegin` and only process widgets inside the window if the value is not 0. Nesting groups is possible and even encouraged since many layouting schemes can only be achieved by nesting. Groups, unlike windows, need `nk.NkGroupEnd` to be only called if the corosponding `nk.NkGroupBegin` call does not return 0.
Note that group names should be unique.

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
