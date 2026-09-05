package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"math/rand/v2"
	"os"

	"meowyplayer/ui/internal/mutil"
	"meowyplayer/ui/internal/mwidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PlaylistEditor struct {
	widget.BaseWidget
	cover         *canvas.Image
	uploadButton  *widget.Button
	refreshButton *widget.Button
	pickButton    *widget.Button
	titleEntry    *widget.Entry
}

// Generates a random GitHub-style 5x5 symmetrical identicon
func generateIdenticon() image.Image {
	const size = 250
	const cells = 5
	const cellSize = size / cells

	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	bgColor := color.NRGBA{240, 240, 240, 255}
	fgColor := color.NRGBA{uint8(rand.N(256)), uint8(rand.N(256)), uint8(rand.N(256)), 255}

	// Generate a 3x5 boolean array for the left half + center vertical line
	pattern := [3][5]bool{}
	for x := 0; x < 3; x++ {
		for y := 0; y < 5; y++ {
			pattern[x][y] = rand.N(2) == 0
		}
	}

	for x := 0; x < cells; x++ {
		for y := 0; y < cells; y++ {
			// Mirror the x coordinate for the right side (columns 3 and 4)
			px := x
			if px > 2 {
				px = 4 - px
			}

			c := bgColor
			if pattern[px][y] {
				c = fgColor
			}

			// Fill the cell block
			for dx := 0; dx < cellSize; dx++ {
				for dy := 0; dy < cellSize; dy++ {
					img.SetNRGBA(x*cellSize+dx, y*cellSize+dy, c)
				}
			}
		}
	}
	return img
}

func newPlaylistEditor() *PlaylistEditor {
	var v PlaylistEditor

	// Cover: use the generated identicon as the default
	v.cover = canvas.NewImageFromImage(generateIdenticon())
	v.cover.SetMinSize(mwidget.PlaylistCardSize)

	// File picker.
	upload := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			slog.Error("failed to read file", "error", err)
		} else if reader != nil {
			v.setImage(reader.URI().Path())
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
	upload.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
	upload.SetConfirmText(lang.L("Upload"))
	upload.SetDismissText(lang.L("Cancel"))
	v.uploadButton = widget.NewButtonWithIcon("", nil, upload.Show)

	// Refresh button to generate a new random identicon.
	v.refreshButton = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		v.cover.Resource = nil
		v.cover.Image = generateIdenticon()
		v.cover.Refresh()
	})

	// Color picker.
	picker := dialog.NewColorPicker("", "", v.setColor, fyne.CurrentApp().Driver().AllWindows()[0])
	picker.Advanced = true
	v.pickButton = widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), picker.Show)

	// Title entry.
	v.titleEntry = widget.NewEntry()
	v.titleEntry.PlaceHolder = lang.L("Enter the title here")

	v.ExtendBaseWidget(&v)
	return &v
}

func newPlaylistEditorWithState(title string, cover fyne.Resource) *PlaylistEditor {
	editor := newPlaylistEditor()
	editor.titleEntry.SetText(title)
	editor.cover.Resource = cover
	return editor
}

func (v *PlaylistEditor) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil,
		container.NewBorder(nil, nil, container.NewHBox(v.refreshButton, v.pickButton), nil, v.titleEntry),
		nil,
		nil,
		container.NewStack(v.uploadButton, v.cover),
	))
}

func (v *PlaylistEditor) state() (string, []byte) {
	// Image has higher priority.
	if v.cover.Image == nil {
		return v.titleEntry.Text, v.cover.Resource.Content()
	}

	buf := bytes.Buffer{}
	if err := png.Encode(&buf, v.cover.Image); err != nil {
		slog.Error("failed to encode image", "error", err)
		return v.titleEntry.Text, theme.BrokenImageIcon().Content()
	}

	return v.titleEntry.Text, buf.Bytes()
}

func (v *PlaylistEditor) setColor(coverColor color.Color) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, coverColor)
	v.cover.Resource = nil
	v.cover.Image = img
	v.cover.Refresh()
}

func (v *PlaylistEditor) setImage(path string) {
	file, err := os.Open(path)
	if err != nil {
		slog.Error("failed to set cover image", "error", err)
		return
	}
	defer file.Close()

	// Resize to reduce UI rendering time.
	img, err := mutil.ScaleImageFromReader(file, mwidget.PlaylistCardSize)
	if err != nil {
		slog.Error("failed to scale the image", "error", err)
		return
	}

	v.cover.Resource = nil
	v.cover.Image = img
	v.cover.Refresh()
}
