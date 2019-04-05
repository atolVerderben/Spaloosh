package isledef

import (
	assets "github.com/atolVerderben/spaloosh/isledef/internal"
	"github.com/atolVerderben/tentsuyu"
)

func loadAudio() *tentsuyu.AudioPlayer {
	audioPlayer, _ := tentsuyu.NewAudioPlayer()
	audioPlayer.AddSoundEffectFromBytes("spaloosh", assets.ReturnSE(seMiss), 1.0)
	audioPlayer.AddSoundEffectFromBytes("kaboom", assets.ReturnSE(seHit), 1.0)
	return audioPlayer
}

const (
	seMiss string = "spaloosh.wav"
	seHit         = "kaboom-ah2.wav"
)
