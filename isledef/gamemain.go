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
	gameStateMsg      tentsuyu.GameStateMsg
	lose, win         bool
	pauseTick         int
	bulletCounter     *BulletCounter
	prevTime          int
	endGameTransition bool
	remainingDisplay  *EnemyDisplay
}

//NewGameMain creates a GameMain state
func NewGameMain(g *tentsuyu.Game) *GameMain {
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
	//g.GameData.SetGameMode(GameModeNormalTimed)
	TimeRanOut = false
	AIBroke = false
	GamePlayer.shots = g.GameData.Settings["ShotMax"].ValueInt
	gm.bulletCounter = NewBulletCounter(96, 152, g.GameData)
	gm.ai.SetBoard(gm.grid)
	gm.hudText = []*tentsuyu.MenuElement{
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(288, 100, 200, 30, g.UIController.ReturnFont(FntSmallPixel), []string{"Time Remaining: "}, color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&gm.timer, 452, 100, 80, 50, g.UIController.ReturnFont(FntSmallPixel), color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(592, 100, 200, 30, g.UIController.ReturnFont(FntSmallPixel), []string{"Monsters Remaining: "}, color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&gm.grid.shipsRemaining, 794, 100, 80, 50, g.UIController.ReturnFont(FntSmallPixel), color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},

		/*&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(520, 60, 250, 50, g.UIController.ReturnFont("font1"), []string{"Shots Left: "}, color.Black, 8),
			Action:     func() {},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&GamePlayer.shots, 780, 60, 80, 50, g.UIController.ReturnFont("font1"), color.Black, 8),
			Action:     func() {},
			Selectable: false,
		},*/
	}
	return gm
}

//Update the main game
func (gm *GameMain) Update(g *tentsuyu.Game) error {
	if g.Input.Button("Escape").JustPressed() {
		gm.gameStateMsg = GameStateMsgPause
		gm.pauseTick = 1
		return nil
	}
	if gm.pauseTick > 0 {
		gm.pauseTick++
		if gm.pauseTick > 60 {
			gm.pauseTick = 0
			gm.prevTime = g.GameData.TimeInSecond()

		}
		return nil
	}
	for _, text := range gm.hudText {
		text.Update(g.Input, 0, 0)
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

	GamePlayer.Update()
	gm.grid.Update(g)
	//gm.tick++
	//log.Printf("Current Time: %v Prev Time: %v\n", g.GameData.TimeInSecond(), gm.prevTime)
	if g.GameData.TimeInSecond()-gm.prevTime == 1 {
		gm.timer--

		gm.prevTime = g.GameData.TimeInSecond()

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
	if GamePlayer.shots == 0 && gm.grid.Cleared == false {
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
func (gm *GameMain) Draw(g *tentsuyu.Game) error {
	DrawBackground(g)
	gm.bulletCounter.Draw(g.Screen)
	gm.grid.Draw(g.Screen)
	gm.remainingDisplay.Draw(g.Screen)
	for _, text := range gm.hudText {
		text.Draw(g.Screen)
	}
	return nil
}

//Msg returns the current state's message
func (gm *GameMain) Msg() tentsuyu.GameStateMsg {
	return gm.gameStateMsg
}

//SetMsg sets the game state message
func (gm *GameMain) SetMsg(msg tentsuyu.GameStateMsg) {
	gm.gameStateMsg = msg
}
