package seeker

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"meowyplayer.com/src/custom_canvas"
	"meowyplayer.com/src/resource"
)

const (
	MAX_MUSIC_QUEUE_SIZE = 128

	//music quality
	SAMPLING_RATE   = 44100
	NUM_OF_CHANNELS = 2
	AUDIO_BIT_DEPTH = 2

	LOOP_ORDER   = 0
	REPEAT_ORDER = 1
	RANDOM_ORDER = 2
	ORDER_LEN    = 3
)

var MusicQueue *MusicSeeker

type MusicSeeker struct {
	curr      int
	musicInfo []custom_canvas.MusicInfo
	mux       sync.Mutex

	//widgets
	Title    *widget.Label
	Play     *SeekerPlay
	Prev     *SeekerPrev
	Next     *SeekerNext
	Progress *SeekerProgress
	Volume   *SeekerVolume
	Order    *SeekerOrder
}

func NewMusicSeeker() (*MusicSeeker, error) {

	s := MusicSeeker{}
	s.musicInfo = make([]custom_canvas.MusicInfo, 0, MAX_MUSIC_QUEUE_SIZE)

	s.Title = widget.NewLabel("")

	//play button
	var err error
	s.Play, err = NewSeekerPlay()
	if err != nil {
		return nil, err
	}
	s.Play.OnTapped = func() {
		if !s.IsEmpty() {
			s.Play.Signal <- struct{}{}
		}
	}

	//prev button
	s.Prev = NewSeekerPrev()
	s.Prev.OnTapped = func() {
		if !s.IsEmpty() {
			s.PrevMusic()
			s.Prev.Signal <- struct{}{}
		}
	}

	//next button
	s.Next = NewSeekerNext()
	s.Next.OnTapped = func() {
		if !s.IsEmpty() {
			s.NextMusic()
			s.Next.Signal <- struct{}{}
		}
	}

	//music progress bar
	s.Progress = NewSeekerProgress()

	//music volume bar
	s.Volume = NewSeekerVolume()

	//music playing order
	s.Order = NewSeekerOrder()

	return &s, nil
}

func (s *MusicSeeker) NewPlaylist() *container.Scroll {

	playlist := container.NewVBox()
	playlistScroll := container.NewVScroll(playlist)
	return playlistScroll
}

func (s *MusicSeeker) SetPlaylist(musicInfo []custom_canvas.MusicInfo, i int) {
	s.mux.Lock()
	s.musicInfo = make([]custom_canvas.MusicInfo, len(musicInfo))
	copy(s.musicInfo, musicInfo)
	s.curr = i

	//must reload seeker
	//assume the music info list is never empty
	s.Next.Signal <- struct{}{}
	s.mux.Unlock()
}

func (s *MusicSeeker) IsEmpty() bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.musicInfo) == 0
}

func (s *MusicSeeker) GetCurrMusic() custom_canvas.MusicInfo {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.musicInfo[s.curr]
}

func (s *MusicSeeker) PrevMusic() {
	s.mux.Lock()
	s.curr--
	if s.curr < 0 {
		s.curr += len(s.musicInfo)
	}
	s.mux.Unlock()
}

func (s *MusicSeeker) NextMusic() {
	s.mux.Lock()

	switch s.Order.Mode {
	case LOOP_ORDER:
		s.curr++
		if s.curr >= len(s.musicInfo) {
			s.curr -= len(s.musicInfo)
		}

	case REPEAT_ORDER:
		//literally nothing

	case RANDOM_ORDER:
		rand.Seed(time.Now().UnixNano())
		s.curr = rand.Int() % len(s.musicInfo)
	}
	s.mux.Unlock()
}

/*

Components

*/

type SeekerPlay struct {
	widget.Button
	Signal chan struct{}

	playIcon  fyne.Resource
	pauseIcon fyne.Resource
}

func NewSeekerPlay() (*SeekerPlay, error) {
	play := &SeekerPlay{}
	play.ExtendBaseWidget(play)

	//load
	var err error
	play.playIcon, err = fyne.LoadResourceFromPath(resource.GetImagePath("seeker_play.png"))
	if err != nil {
		return nil, err
	}

	play.pauseIcon, err = fyne.LoadResourceFromPath(resource.GetImagePath("seeker_pause.png"))
	if err != nil {
		return nil, err
	}

	play.SetText("")
	play.SetIcon(play.playIcon)
	play.Signal = make(chan struct{}, 4)

	return play, nil
}

func (play *SeekerPlay) UpdateIcon(shouldPlay bool) {
	if shouldPlay {
		play.SetIcon(play.pauseIcon)
	} else {
		play.SetIcon(play.playIcon)
	}
}

type SeekerPrev struct {
	widget.Button
	Signal chan struct{}
}

func NewSeekerPrev() *SeekerPrev {
	prev := &SeekerPrev{}
	prev.ExtendBaseWidget(prev)

	prev.SetText("<<")
	prev.Signal = make(chan struct{}, 4)

	return prev
}

type SeekerNext struct {
	widget.Button
	Signal chan struct{}
}

func NewSeekerNext() *SeekerNext {
	next := &SeekerNext{}
	next.ExtendBaseWidget(next)

	next.SetText(">>")
	next.Signal = make(chan struct{}, 4)

	return next
}

type SeekerProgress struct {
	widget.Slider
	Ignore     chan struct{}
	decodedMP3 *mp3.Decoder
}

func NewSeekerProgress() *SeekerProgress {
	s := SeekerProgress{}
	s.ExtendBaseWidget(&s)

	s.Min = 0.0
	s.Max = 1.0
	s.Step = 1.0
	return &s
}

func (s *SeekerProgress) BindMP3(mp3 *mp3.Decoder, player *oto.Player) {

	s.decodedMP3 = mp3

	//set up the progress bar range
	s.Max = float64(mp3.Length())

	//sometime the play list reload sends an extra signal
	//2 avoids buffer locks
	s.Ignore = make(chan struct{}, 2)
	s.OnChanged = func(tick float64) {
		select {
		case <-s.Ignore:

		default:

			wasPlaying := (*player).IsPlaying()

			if wasPlaying {
				(*player).Pause()
			}

			iTick := int64(tick)
			mp3.Seek(iTick-iTick%4, io.SeekStart)

			if wasPlaying {
				(*player).Play()
			}
		}
	}
}

func (s *SeekerProgress) UpdateBar() {

	s.Ignore <- struct{}{}
	currTick, err := s.decodedMP3.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Panic(err)
	}
	s.SetValue(float64(currTick))
	log.Printf("%v / %v (%.4f)\n", currTick, s.decodedMP3.Length(), (float64(currTick) / float64(s.decodedMP3.Length())))
}

type SeekerVolume struct {
	widget.Slider
}

func NewSeekerVolume() *SeekerVolume {

	s := SeekerVolume{}
	s.ExtendBaseWidget(&s)

	s.Min = 0.0
	s.Max = 1.0
	s.Step = 0.01

	return &s
}

func (s *SeekerVolume) BindMP3(mp3 *mp3.Decoder, player *oto.Player) {

	s.SetValue((*player).Volume())
	s.OnChanged = func(volume float64) {
		(*player).SetVolume(volume)
		log.Printf("volume: %v\n", volume)
	}
}

type SeekerOrder struct {
	widget.Button
	Mode  int
	icons []fyne.Resource
}

func NewSeekerOrder() *SeekerOrder {
	s := SeekerOrder{}
	s.ExtendBaseWidget(&s)

	//load icon resource
	s.icons = make([]fyne.Resource, 0, ORDER_LEN)
	for i := 0; i < ORDER_LEN; i++ {
		icon, err := fyne.LoadResourceFromPath(resource.GetImagePath(fmt.Sprintf("seeker_order_icon_%v.png", i)))
		if err != nil {
			log.Panic(err)
		}
		s.icons = append(s.icons, icon)
	}

	s.Mode = RANDOM_ORDER
	s.updateIcon()

	return &s
}

func (s *SeekerOrder) Tapped(_ *fyne.PointEvent) {

	if s.Mode == RANDOM_ORDER {
		s.Mode = LOOP_ORDER
	} else {
		s.Mode++
	}
	s.updateIcon()
}

func (s *SeekerOrder) updateIcon() {
	log.Printf("set icon %v\n", s.Mode)
	s.SetIcon(s.icons[s.Mode])
}
