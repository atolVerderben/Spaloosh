package main

import (
	"github.com/atolVerderben/spaloosh/isledef"
	"github.com/hajimehoshi/ebiten"
)

//Size of the game
const (
	ScreenWidth  = 852 //512
	ScreenHeight = 480
)

func main() {
	game, err := isledef.NewGame(ScreenWidth, ScreenHeight)
	if err != nil {
		panic(err)
	}
	if err := ebiten.Run(game.Loop, ScreenWidth, ScreenHeight, 1, "SPALOOSH!"); err != nil {
		panic(err)
	}
}
