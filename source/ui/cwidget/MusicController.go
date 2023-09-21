package cwidget

import (
	"bytes"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/source/path"
	"meowyplayer.com/source/player"
	"meowyplayer.com/source/utility"
)

type MusicController struct {
	widget.BaseWidget
	*titleDisplay
	progressController *progressController
	*buttonController
	*volumeController

	*oto.Context
	playChan chan player.Play
}

func NewMusicController() *MusicController {
	//initialize oto
	context, ready, err := oto.NewContext(player.SAMPLING_RATE, player.NUM_OF_CHANNELS, player.AUDIO_BIT_DEPTH)
	utility.MustNil(err)
	<-ready

	controller := &MusicController{
		widget.BaseWidget{},
		newTitleDisplay(),
		newProgressController(),
		newButtonController(),
		newVolumeController(),
		context,
		make(chan player.Play),
	}
	defer func() { go controller.start() }()

	controller.ExtendBaseWidget(controller)
	return controller
}

func (c *MusicController) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		c.titleDisplay,
		container.NewGridWithRows(1, layout.NewSpacer(), c.buttonController, layout.NewSpacer(), c.volumeController),
		nil,
		nil,
		c.progressController,
	))
}

func (c *MusicController) Notify(play *player.Play) {
	c.playChan <- *play
}

func (c *MusicController) start() {

	//wait for the user to load the music
	currentPlay := <-c.playChan

	for {
		//load music
		mp3Data, err := os.ReadFile(path.Music(currentPlay.Music()))
		utility.MustNil(err)
		mp3Decoder, err := mp3.NewDecoder(bytes.NewReader(mp3Data))
		utility.MustNil(err)
		mp3Player := c.NewPlayer(mp3Decoder)

		//update GUI
		c.SetMusicTitle(currentPlay.Music())
		c.BindVolume(mp3Player)
		c.BindButton(mp3Player)

		//play music
		log.Printf("playing %v\n", currentPlay.Music().GoodTitle())
		interrupt := false
	MusicLoop:
		for mp3Player.Play(); mp3Player.UnplayedBufferSize() > 0; {
			select {
			case currentPlay = <-c.playChan:
				interrupt = true
				break MusicLoop
			default:
			}
		}

		//next music
		if !interrupt {
			currentPlay.NextMusic()
		}

		mp3Player.Close()
	}
}
