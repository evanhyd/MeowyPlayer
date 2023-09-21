package cwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"meowyplayer.com/source/player"
)

type progressController struct {
	widget.BaseWidget
	progressSlider *widget.Slider
	durationLabel  *widget.Label
}

func newProgressController() *progressController {
	controller := &progressController{widget.BaseWidget{}, widget.NewSlider(0.0, 1.0), widget.NewLabel("00:00")}
	controller.progressSlider.Step = 0.001
	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *progressController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, nil, c.durationLabel, c.progressSlider))
}

func (c *progressController) SetMusic(music *player.Music) {
	c.progressSlider.SetValue(0.0)
}

// func (c *progressController) Bind(decoder **mp3.Decoder, mp3Player oto.MP3Player) {
// 	c.progressSlider.OnChanged = func(percent float64) {
// 		//change the seeker bar
// 		loc := int64(float64((*decoder).Length()) * percent)
// 		loc -= loc % 4
// 		_, err := (*decoder).Seek(loc, io.SeekStart)
// 		utility.MustNil(err)
// 	}

// 	go func() {
// 		for (*decoder).Length() != -1 {
// 			loc, _ := (*decoder).Seek(0, io.SeekCurrent)
// 			percent := float64(loc) / float64((*decoder).Length())
// 		}
// 	}()
// }
