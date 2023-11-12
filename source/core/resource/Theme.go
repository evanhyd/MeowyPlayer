package resource

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type vanillaTheme struct {
	darkTheme fyne.Theme
}

func VanillaTheme() fyne.Theme {
	return &vanillaTheme{theme.DarkTheme()}
}

func (t *vanillaTheme) Color(colorName fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	primary := fyne.CurrentApp().Settings().PrimaryColor()
	switch colorName {
	case theme.ColorNamePrimary:
		return theme.PrimaryColorNamed(primary)
	case theme.ColorNameFocus:
		return t.focusColorNamed(primary)
	case theme.ColorNameSelection:
		return t.selectionColorNamed(primary)
	default:
		return t.darkPaletColorNamed(colorName)
	}
}

func (t *vanillaTheme) Font(textStyle fyne.TextStyle) fyne.Resource {
	if textStyle.Bold && textStyle.Italic {
		return BoldItalicFont
	}

	if textStyle.Bold {
		return BoldFont
	}

	if textStyle.Italic {
		return ItalicFont
	}

	// font pack bloats up the binary
	// if textStyle.Monospace {
	// 	return t.monospace
	// }

	return RegularFont
}

func (t *vanillaTheme) Icon(iconName fyne.ThemeIconName) fyne.Resource {
	return t.darkTheme.Icon(iconName)
}

func (t *vanillaTheme) Size(size fyne.ThemeSizeName) float32 {
	switch size {
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameLineSpacing:
		return 4
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameText:
		return 13
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInputBorder:
		return 1
	default:
		return 0
	}
}

func (t *vanillaTheme) darkPaletColorNamed(name fyne.ThemeColorName) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x14, G: 0x14, B: 0x15, A: 0xff}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x39, G: 0x39, B: 0x3a, A: 0xff}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xf3, G: 0xf3, B: 0xf3, A: 0xff}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x0f}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x20, G: 0x20, B: 0x23, A: 0xff}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 0x39, G: 0x39, B: 0x3a, A: 0xff}
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 0x28, G: 0x29, B: 0x2e, A: 0xff}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0x18, G: 0x1d, B: 0x25, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0xb2, G: 0xb2, B: 0xb2, A: 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x66}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0x99}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xff}
	case theme.ColorNameShadow:
		return color.NRGBA{A: 0x66}
	case theme.ColorNameSuccess:
		return theme.SuccessColor()
	case theme.ColorNameWarning:
		return theme.WarningColor()
	case theme.ColorNameError:
		return theme.ErrorColor()
	}

	return color.Transparent
}

func (t *vanillaTheme) focusColorNamed(name string) color.NRGBA {
	switch name {
	case theme.ColorRed:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0x7f}
	case theme.ColorOrange:
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0x7f}
	case theme.ColorYellow:
		return color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x7f}
	case theme.ColorGreen:
		return color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0x7f}
	case theme.ColorPurple:
		return color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0x7f}
	case theme.ColorBrown:
		return color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0x7f}
	case theme.ColorGray:
		return color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0x7f}
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return color.NRGBA{R: 0x00, G: 0x6C, B: 0xff, A: 0x2a}
}
func (t *vanillaTheme) selectionColorNamed(name string) color.NRGBA {
	switch name {
	case theme.ColorRed:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0x3f}
	case theme.ColorOrange:
		return color.NRGBA{R: 0xff, G: 0x98, B: 0x00, A: 0x3f}
	case theme.ColorYellow:
		return color.NRGBA{R: 0xff, G: 0xeb, B: 0x3b, A: 0x3f}
	case theme.ColorGreen:
		return color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0x3f}
	case theme.ColorPurple:
		return color.NRGBA{R: 0x9c, G: 0x27, B: 0xb0, A: 0x3f}
	case theme.ColorBrown:
		return color.NRGBA{R: 0x79, G: 0x55, B: 0x48, A: 0x3f}
	case theme.ColorGray:
		return color.NRGBA{R: 0x9e, G: 0x9e, B: 0x9e, A: 0x3f}
	}

	// We return the value for ColorBlue for every other value.
	// There is no need to have it in the switch above.
	return color.NRGBA{R: 0x00, G: 0x6C, B: 0xff, A: 0x40}
}
