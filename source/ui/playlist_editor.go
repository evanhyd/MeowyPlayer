package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"meowyplayer/ui/internal/widgets"
	"os"

	_ "image/jpeg"

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
	cover        *canvas.Image
	uploadButton *widget.Button
	pickButton   *widget.Button
	titleEntry   *widget.Entry
}

func newPlaylistEditor() *PlaylistEditor {
	var v PlaylistEditor

	// Cover.
	v.cover = canvas.NewImageFromResource(theme.DocumentCreateIcon())
	v.cover.SetMinSize(fyne.NewSize(180, 180))

	// File picker.
	upload := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			slog.Error("failed to read file", "error", err)
		} else if reader != nil {
			v.setImage(reader.URI().Path())
		}
	}, fyne.CurrentApp().Driver().AllWindows()[0])
	upload.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg"}))
	upload.SetConfirmText(lang.L("Upload"))
	upload.SetDismissText(lang.L("Cancel"))
	v.uploadButton = widget.NewButtonWithIcon("", nil, upload.Show)

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
		container.NewBorder(nil, nil, v.pickButton, nil, v.titleEntry),
		nil,
		nil,
		container.NewStack(v.uploadButton, v.cover),
	))
}

func (v *PlaylistEditor) state() (string, fyne.Resource) {
	return v.titleEntry.Text, v.cover.Resource
}

func (v *PlaylistEditor) setColor(coverColor color.Color) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, coverColor)
	data := bytes.Buffer{}
	if err := png.Encode(&data, img); err != nil {
		slog.Error("failed to set cover color", "error", err)
		return
	}
	v.cover.Resource = fyne.NewStaticResource("", data.Bytes())
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
	image, err := widgets.ScaleImageFromReader(file, 140, 140)
	if err != nil {
		slog.Error("failed to scale the image", "error", err)
	}

	buffer := bytes.Buffer{}
	png.Encode(&buffer, image)
	v.cover.Resource = &fyne.StaticResource{StaticContent: buffer.Bytes()}
	v.cover.Refresh()
}
