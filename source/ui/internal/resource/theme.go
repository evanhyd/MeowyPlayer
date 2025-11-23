package resource

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type VanillaTheme struct {
	builtin fyne.Theme
}

func NewVanillaTheme() fyne.Theme {
	return &VanillaTheme{builtin: theme.DefaultTheme()}
}

func (t *VanillaTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return t.builtin.Color(name, theme.VariantDark)
}

func (t *VanillaTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.builtin.Font(style)
}

func (t *VanillaTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.builtin.Icon(name)
}

func (t *VanillaTheme) Size(size fyne.ThemeSizeName) float32 {
	return t.builtin.Size(size)
}
