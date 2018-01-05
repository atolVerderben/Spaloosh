package isledef

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
)

//GameMain is the main game...
type GameMain struct {
	grid              *Grid
	ships             []*Ship
	ZoomLevel         float64
	ai                *AI
	hudText           []*tentsuyu.MenuElement
	timer             int
	tick              int
	gameStateMsg      GameStateMsg
	lose, win         bool
	pauseTick         int
	bulletCounter     *BulletCounter
	prevTime          int
	endGameTransition bool
	remainingDisplay  *EnemyDisplay
}

//NewGameMain creates a GameMain state
func NewGameMain(g *Game) *GameMain {
	x := 288.0
	y := 120.0
	//x := 206.0 * 2
	//y := 52.0 * 2
	ZoomLevel = 1.0
	gm := &GameMain{
		grid:             CreateGrid(x, y, 8, 8),
		ships:            []*Ship{},
		ai:               &AI{},
		timer:            30,
		pauseTick:        1,
		win:              false,
		lose:             false,
		remainingDisplay: NewEnemyDisplay(592, 120, true), //&EnemyDisplay{x: 528, y: 120},
	}
	//g.gameData.SetGameMode(GameModeNormalTimed)
	TimeRanOut = false
	AIBroke = false
	g.player.shots = g.gameData.shotMax
	gm.bulletCounter = NewBulletCounter(96, 152, g.gameData)
	gm.ai.SetBoard(gm.grid)
	gm.hudText = []*tentsuyu.MenuElement{
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(288, 100, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Time Remaining: "}, color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&gm.timer, 452, 100, 80, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(592, 100, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Monsters Remaining: "}, color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&gm.grid.shipsRemaining, 794, 100, 80, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},

		/*&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(520, 60, 250, 50, tentsuyu.Components.ReturnFont("font1"), []string{"Shots Left: "}, color.Black, 8),
			Action:     func() {},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&g.player.shots, 780, 60, 80, 50, tentsuyu.Components.ReturnFont("font1"), color.Black, 8),
			Action:     func() {},
			Selectable: false,
		},*/
	}
	return gm
}

//Update the main game
func (gm *GameMain) Update(g *Game) error {
	if tentsuyu.Input.Button("Escape").JustPressed() {
		gm.gameStateMsg = GameStateMsgPause
		gm.pauseTick = 1
		return nil
	}
	if gm.pauseTick > 0 {
		gm.pauseTick++
		if gm.pauseTick > 60 {
			gm.pauseTick = 0
			gm.prevTime = g.gameData.TimeInSecond()

		}
		return nil
	}
	for _, text := range gm.hudText {
		text.Update()
	}
	gm.remainingDisplay.Update(gm.grid.Ships)
	if gm.endGameTransition == true {
		if gm.grid.allVisible == false {
			gm.grid.MakeAllVisible()
		}
		if gm.tick < 120 {
			gm.tick++
		} else {
			if gm.win == true {
				gm.gameStateMsg = GameStateGameWin
			} else if gm.lose == true {
				gm.gameStateMsg = GameStateGameOver
			}
		}
		return nil
	}

	g.player.Update()
	gm.grid.Update(g)
	//gm.tick++
	//log.Printf("Current Time: %v Prev Time: %v\n", g.gameData.TimeInSecond(), gm.prevTime)
	if g.gameData.TimeInSecond()-gm.prevTime == 1 {
		gm.timer--

		gm.prevTime = g.gameData.TimeInSecond()

	}
	/*if gm.tick/60 >= 1 {
		gm.tick = 0
		gm.timer--
	}*/
	if gm.grid.shipsRemaining == 0 && gm.win == false {
		gm.tick = 0
		//gm.timer = 1
		gm.win = true
		gm.endGameTransition = true
	}
	if g.player.shots == 0 && gm.grid.Cleared == false {
		//gm.timer = 0
		gm.lose = true
		gm.endGameTransition = true
	}
	if gm.timer <= 0 {
		gm.timer = 0
		if gm.grid.Cleared {
			gm.win = true

		} else {
			gm.lose = true
			TimeRanOut = true

		}
		gm.endGameTransition = true

	}
	return nil
}

//Draw the main game
func (gm *GameMain) Draw(g *Game) error {
	g.DrawBackground()
	gm.bulletCounter.Draw(g.screen)
	gm.grid.Draw(g.screen)
	gm.remainingDisplay.Draw(g.screen)
	for _, text := range gm.hudText {
		text.Draw(g.screen)
	}
	return nil
}

//Msg returns the current state's message
func (gm *GameMain) Msg() GameStateMsg {
	return gm.gameStateMsg
}

//SetMsg sets the game state message
func (gm *GameMain) SetMsg(msg GameStateMsg) {
	gm.gameStateMsg = msg
}
