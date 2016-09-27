package nk

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
