## nk-example-layout

This example shows, how to use the different layout options provided by nuklear API.

Please note, that the following explanations are explainationstaken from the nuklear doc, available [here](https://github.com/vurtun/nuklear/tree/master/doc)

## Layouting
Layouting in general describes placing widget inside a window with position and size. While in this particular implementation there are five different APIs for layouting each with different trade offs between control and ease of use.

All layouting methods in this library are based around the concept of a row. A row has a height and a number of columns and each layouting method specifies how each widget is placed inside the row. After a row has been allocated by calling a layouting functions it can be filled with widgets.

To actually define a layout you just call the appropriate layouting function and each subsequent widget call will place the widget as specified. Important here is that if you define more widgets then columns defined inside the layout functions it will allocate the next row without you having to make another layouting call. 

Biggest limitation with using all these APIs outside the `nk.NkLayoutSpace` API is that you have to define the row height for each. However the row height often depends on the height of the font.

To fix that internally nuklear uses a minimum row height that is set to the height plus padding of currently active font and overwrites the row height value if zero.

If you manually want to change the minimum row height then use `nk.NkLayoutSetMinRowHeight`, and use `nk.NkLayoutResetMinRowHeight` to reset it back to be derived from font height.

Also if you change the font in nuklear it will automatically change the minimum row height for you and. This means if you change the font but still want a minimum row height smaller than the font you have to repush your value.

For actually more advanced UI I would even recommend using the `nk.NkLayoutSpace` layouting method in combination with a cassowary constraint solver (there are some versions on github with permissive license model) to take over all control over widget layouting yourself. However for quick and dirty layouting using all the other layouting functions should be fine. 

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
<mkawaganga@gmail.com>
