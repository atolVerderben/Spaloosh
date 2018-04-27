package isledef

import (
	"time"

	assets "github.com/atolVerderben/spaloosh/isledef/internal"
	"github.com/hajimehoshi/ebiten/audio"
)

var (
	audioContext   *audio.Context
	soundFilenames = []string{
		"kaboom-ah2.wav",
		"spaloosh.wav",
	}
	soundPlayers = map[string]*audio.Player{}
)

type SoundPlayer struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
}

func init() {
	const sampleRate = 44100
	var err error
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}
}

type emptyAudio struct {
}

func (e *emptyAudio) Read(b []byte) (int, error) {
	// TODO: Clear b?
	return len(b), nil
}

func (e *emptyAudio) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (e *emptyAudio) Close() error {
	return nil
}

func loadAudio() error {
	/*for _, n := range soundFilenames {
		/*b, err := assets.Asset("resources/sound/" + n)
		if err != nil {
			return err
		}

		oggF, err := ebitenutil.OpenFile("assets/audio/" + n)
		if err != nil {
			log.Fatal(err)
		}

		//f := audio.BytesReadSeekCloser(b)
		var s audio.ReadSeekCloser
		switch {
		case strings.HasSuffix(n, ".ogg"):
			stream, err := vorbis.Decode(audioContext, oggF)
			if err != nil {
				s = &emptyAudio{}
			} else {
				//if n == "FindYouMarch.ogg" {
				s = audio.NewInfiniteLoop(stream, stream.Size())
				//} else {
				//	s = stream
				//}
			}
		case strings.HasSuffix(n, ".wav"):
			stream, err := wav.Decode(audioContext, oggF)
			if err != nil {
				return err
			}
			s = stream

		default:
			panic("invalid file name")
		}
		p, err := audio.NewPlayer(audioContext, s)
		if err != nil {
			return err
		}

		soundPlayers[n] = p

	}*/
	return nil
}

func finalizeAudio() error {
	for _, p := range soundPlayers {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}

type BGM string

const (
	BGM0 BGM = "FindYouMarch.ogg"
	BGM1 BGM = "high seas.ogg"
)

func SetBGMVolume(volume float64) {
	for _, b := range []BGM{BGM0, BGM1} {
		p := soundPlayers[string(b)]
		if !p.IsPlaying() {
			continue
		}
		p.SetVolume(volume)
		return
	}
}

func PauseBGM() error {
	for _, b := range []BGM{BGM0, BGM1} {
		p := soundPlayers[string(b)]
		if err := p.Pause(); err != nil {
			return err
		}
	}
	return nil
}

func ResumeBGM(bgm BGM) error {
	if err := PauseBGM(); err != nil {
		return err
	}
	p := soundPlayers[string(bgm)]
	p.SetVolume(1)
	return p.Play()
}

func PlayBGM(bgm BGM) error {
	if err := PauseBGM(); err != nil {
		return err
	}
	p := soundPlayers[string(bgm)]
	p.SetVolume(1)
	if err := p.Rewind(); err != nil {
		return err
	}
	return p.Play()
}

type SE string

const (
	seMiss string = "spaloosh.wav"
	seHit         = "kaboom-ah2.wav"
)

func PlaySE(se string) error {
	if !PlaySoundEffects {
		return nil
	}
	if se != seMiss && se != seHit {
		return nil
	}

	sePlayer, _ := audio.NewPlayerFromBytes(audioContext, assets.ReturnSE(se))
	sePlayer.SetVolume(0.05)

	if err := sePlayer.Play(); err != nil {
		return err
	}
	return nil

}
