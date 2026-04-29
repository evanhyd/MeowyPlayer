package ui

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"

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
	cover        *canvas.Image
	uploadButton *widget.Button
	pickButton   *widget.Button
	titleEntry   *widget.Entry
}

func newPlaylistEditor() *PlaylistEditor {
	var v PlaylistEditor

	// Cover.
	v.cover = canvas.NewImageFromResource(theme.UploadIcon())
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
	img, err := mwidget.ScaleImageFromReader(file, mwidget.PlaylistCardSize)
	if err != nil {
		slog.Error("failed to scale the image", "error", err)
		return
	}

	v.cover.Resource = nil
	v.cover.Image = img
	v.cover.Refresh()
}
