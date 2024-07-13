package view

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"playground/view/internal/cwidget"
	"playground/view/internal/resource"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AlbumEditor struct {
	widget.BaseWidget
	cover        *canvas.Image
	uploadButton *widget.Button
	pickButton   *widget.Button
	titleEntry   *widget.Entry
}

func newAlbumEditor() *AlbumEditor {
	var v AlbumEditor

	//cover
	v.cover = canvas.NewImageFromResource(theme.DocumentCreateIcon())
	v.cover.SetMinSize(resource.KAlbumCoverSize)
	v.cover.FillMode = canvas.ImageFillContain

	//file picker
	upload := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			fyne.LogError("failed to read file", err)
		} else if reader != nil {
			v.setImage(reader.URI().Path())
		}
	}, getWindow())
	upload.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", "jpeg", ".bmp"}))
	upload.SetConfirmText(resource.KUploadText)
	upload.SetDismissText(resource.KCancelText)
	v.uploadButton = cwidget.NewButtonIcon(nil, upload.Show)

	//color picker
	picker := dialog.NewColorPicker("", "", v.setColor, getWindow())
	picker.Advanced = true
	v.pickButton = cwidget.NewButtonIcon(theme.ColorPaletteIcon(), picker.Show)

	//title entry
	v.titleEntry = widget.NewEntry()
	v.titleEntry.PlaceHolder = resource.KEnterTitleHint

	v.ExtendBaseWidget(&v)
	return &v
}

func newAlbumEditorWithState(title string, cover fyne.Resource) *AlbumEditor {
	editor := newAlbumEditor()
	editor.titleEntry.SetText(title)
	editor.cover.Resource = cover
	return editor
}

func (v *AlbumEditor) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil,
		container.NewBorder(nil, nil, v.pickButton, nil, v.titleEntry),
		nil,
		nil,
		container.NewStack(v.uploadButton, v.cover),
	))
}

func (v *AlbumEditor) state() (string, fyne.Resource) {
	return v.titleEntry.Text, v.cover.Resource
}

func (v *AlbumEditor) setColor(coverColor color.Color) {
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, coverColor)
	data := bytes.Buffer{}
	if err := png.Encode(&data, img); err != nil {
		fyne.LogError("failed to set cover color", err)
	}
	v.cover.Resource = fyne.NewStaticResource("", data.Bytes())
	v.cover.Refresh()
}

func (v *AlbumEditor) setImage(path string) {
	var err error
	if v.cover.Resource, err = fyne.LoadResourceFromPath(path); err != nil {
		fyne.LogError("failed to set cover image", err)
	}
	v.cover.Refresh()
}
