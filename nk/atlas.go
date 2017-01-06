package nk

func (atlas *FontAtlas) DefaultFont() *Font {
	return (*Font)(atlas.default_font)
}
