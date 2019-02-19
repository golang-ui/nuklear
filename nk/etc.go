package nk

/*
#include "nuklear.h"
*/
import "C"
import "unsafe"

func (ctx *Context) Input() *Input {
	return (*Input)(&ctx.input)
}

func (ctx *Context) Style() *Style {
	return (*Style)(&ctx.style)
}

func (ctx *Context) Memory() *Buffer {
	return (*Buffer)(&ctx.memory)
}

func (ctx *Context) Clip() *Clipboard {
	return (*Clipboard)(&ctx.clip)
}

func (ctx *Context) LastWidgetState() Flags {
	return (Flags)(ctx.last_widget_state)
}

func (ctx *Context) DeltaTimeSeconds() float32 {
	return (float32)(ctx.delta_time_seconds)
}

func (ctx *Context) ButtonBehavior() ButtonBehavior {
	return (ButtonBehavior)(ctx.button_behavior)
}

func (ctx *Context) Stacks() *ConfigurationStacks {
	return (*ConfigurationStacks)(&ctx.stacks)
}

func (input *Input) Mouse() *Mouse {
	return (*Mouse)(&input.mouse)
}

func (input *Input) Keyboard() *Keyboard {
	return (*Keyboard)(&input.keyboard)
}

func (keyboard *Keyboard) Text() string {
	return C.GoStringN(&keyboard.text[0], keyboard.text_len)
}

func (mouse *Mouse) Grab() bool {
	return mouse.grab == True
}

func (mouse *Mouse) Grabbed() bool {
	return mouse.grabbed == True
}

func (mouse *Mouse) Ungrab() bool {
	return mouse.ungrab == True
}

func (mouse *Mouse) ScrollDelta() Vec2 {
	return (Vec2)(mouse.scroll_delta)
}

func (mouse *Mouse) Pos() (int32, int32) {
	return (int32)(mouse.pos.x), (int32)(mouse.pos.y)
}

func (mouse *Mouse) SetPos(x, y int32) {
	mouse.pos.x = (C.float)(x)
	mouse.pos.y = (C.float)(y)
}

func (mouse *Mouse) Prev() (int32, int32) {
	return (int32)(mouse.prev.x), (int32)(mouse.prev.y)
}

func (mouse *Mouse) Delta() (int32, int32) {
	return (int32)(mouse.delta.x), (int32)(mouse.delta.y)
}

var VertexLayoutEnd = DrawVertexLayoutElement{
	Attribute: VertexAttributeCount,
	Format:    FormatCount,
	Offset:    0,
}

func NkDrawForeach(ctx *Context, b *Buffer, fn func(cmd *DrawCommand)) {
	cmd := Nk_DrawBegin(ctx, b)
	for {
		if cmd == nil {
			break
		}
		fn(cmd)
		cmd = Nk_DrawNext(cmd, b, ctx)
	}
}

func NkFontAtlasAddFromBytes(atlas *FontAtlas, data []byte, height float32, config *FontConfig) *Font {
	dataPtr := unsafe.Pointer((*sliceHeader)(unsafe.Pointer(&data)).Data)
	return NkFontAtlasAddFromMemory(atlas, dataPtr, Size(len(data)), height, config)
}

func (fc *FontConfig) SetPixelSnap(b bool) {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	fc.pixel_snap = (C.uchar)(i)
}

func (fc *FontConfig) SetOversample(v, h int) {
	fc.oversample_v = (C.uchar)(v)
	fc.oversample_h = (C.uchar)(h)
}

func (fc *FontConfig) SetRange(r *Rune) {
	fc._range = (*C.nk_rune)(r)
}

func (fc *FontConfig) SetRangeGoRune(r []rune) {
	fc._range = (*C.nk_rune)(unsafe.Pointer(&r[0]))
}

func (h Handle) ID() int {
	return int(*(*int64)(unsafe.Pointer(&h)))
}

func (h Handle) Ptr() uintptr {
	return *(*uintptr)(unsafe.Pointer(&h))
}

func (atlas *FontAtlas) DefaultFont() *Font {
	return (*Font)(atlas.default_font)
}

func (c Color) R() Byte {
	return Byte(c.r)
}

func (c Color) G() Byte {
	return Byte(c.g)
}

func (c Color) B() Byte {
	return Byte(c.b)
}

func (c Color) A() Byte {
	return Byte(c.a)
}

func (c Color) RGBA() (Byte, Byte, Byte, Byte) {
	return Byte(c.r), Byte(c.g), Byte(c.b), Byte(c.a)
}

func (c Color) RGBAi() (int32, int32, int32, int32) {
	return int32(c.r), int32(c.g), int32(c.b), int32(c.a)
}

func (c *Color) SetR(r Byte) {
	c.r = (C.nk_byte)(r)
}

func (c *Color) SetG(g Byte) {
	c.g = (C.nk_byte)(g)
}

func (c *Color) SetB(b Byte) {
	c.b = (C.nk_byte)(b)
}

func (c *Color) SetA(a Byte) {
	c.a = (C.nk_byte)(a)
}

func (c *Color) SetRGBA(r, g, b, a Byte) {
	c.r = (C.nk_byte)(r)
	c.g = (C.nk_byte)(g)
	c.b = (C.nk_byte)(b)
	c.a = (C.nk_byte)(a)
}

func (c *Color) SetRGBAi(r, g, b, a int32) {
	c.r = (C.nk_byte)(r)
	c.g = (C.nk_byte)(g)
	c.b = (C.nk_byte)(b)
	c.a = (C.nk_byte)(a)
}

func (cmd *DrawCommand) ElemCount() int {
	return (int)(cmd.elem_count)
}

func (cmd *DrawCommand) Texture() Handle {
	return (Handle)(cmd.texture)
}

func (cmd *DrawCommand) ClipRect() *Rect {
	return (*Rect)(&cmd.clip_rect)
}

func (f Font) Handle() *UserFont {
	return NewUserFontRef(unsafe.Pointer(&f.handle))
}

func (r *Rect) X() float32 {
	return (float32)(r.x)
}

func (r *Rect) Y() float32 {
	return (float32)(r.y)
}

func (r *Rect) W() float32 {
	return (float32)(r.w)
}

func (r *Rect) H() float32 {
	return (float32)(r.h)
}

func (v *Vec2) X() float32 {
	return (float32)(v.x)
}

func (v *Vec2) Y() float32 {
	return (float32)(v.y)
}

func (v *Vec2) SetX(x float32) {
	v.x = (C.float)(x)
}

func (v *Vec2) SetY(y float32) {
	v.y = (C.float)(y)
}

func (v *Vec2) Reset() {
	v.x = 0
	v.y = 0
}

// Allocated is the total amount of memory allocated.
func (b *Buffer) Allocated() int {
	return (int)(b.allocated)
}

// Size is the current size of the buffer.
func (b *Buffer) Size() int {
	return (int)(b.size)
}

// Type is the memory management type of the buffer.
func (b *Buffer) Type() AllocationType {
	return (AllocationType)(b._type)
}

func (l *ListView) Begin() int {
	return (int)(l.begin)
}

func (l *ListView) End() int {
	return (int)(l.end)
}

func (l *ListView) Count() int {
	return (int)(l.count)
}

func (panel *Panel) Bounds() *Rect {
	return (*Rect)(&panel.bounds)
}

func (t *StyleText) Color() *Color {
	return (*Color)(&t.color)
}

func (s *Style) Text() *StyleText {
	return (*StyleText)(&s.text)
}

func (s *Style) Window() *StyleWindow {
	return (*StyleWindow)(&s.window)
}

func (w *StyleWindow) Background() *Color {
	return (*Color)(&w.background)
}

func (w *StyleWindow) Spacing() *Vec2 {
	return (*Vec2)(&w.spacing)
}

func (w *StyleWindow) Padding() *Vec2 {
	return (*Vec2)(&w.padding)
}

func (w *StyleWindow) GroupPadding() *Vec2 {
	return (*Vec2)(&w.group_padding)
}

func SetSpacing(ctx *Context, v Vec2) {
	*ctx.Style().Window().Spacing() = v
}

func SetPadding(ctx *Context, v Vec2) {
	*ctx.Style().Window().Padding() = v
}

func SetGroupPadding(ctx *Context, v Vec2) {
	*ctx.Style().Window().GroupPadding() = v
}

func SetTextColor(ctx *Context, color Color) {
	*ctx.Style().Text().Color() = color
}

func SetBackgroundColor(ctx *Context, color Color) {
	ctx.Style().Window().fixed_background = C.struct_nk_style_item(NkStyleItemColor(color))
}
