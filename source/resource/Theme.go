package resource

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type vanillaTheme struct {
}

func VanillaTheme() vanillaTheme {
	return vanillaTheme{}
}

func (t *vanillaTheme) Color(colorName fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	panic("not implemented")
	// return color.Color("none")
	// if colorName == theme.ColorNamePrimary {
	// 	return primaryColorNamed(primary)
	// } else if colorName == theme.ColorNameFocus {
	// 	return focusColorNamed(primary)
	// } else if colorName == theme.ColorNameSelection {
	// 	return selectionColorNamed(primary)
	// }

	// if variant == theme.VariantLight {
	// 	return lightPaletColorNamed(colorName)
	// }

	// return darkPaletColorNamed(colorName)
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

	// if textStyle.Monospace {
	// 	return t.monospace
	// }

	return RegularFont
}

func (t *vanillaTheme) Icon(iconName fyne.ThemeIconName) fyne.Resource {
	return WindowIcon
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
