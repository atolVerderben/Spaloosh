package isledef

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
)

//GameBattle is the main game...
type GameBattle struct {
	playerGrid                            *Grid
	enemyGrid                             *Grid
	playerDisplay                         *EnemyDisplay
	aiDisplay                             *EnemyDisplay
	ships                                 []*Ship
	ZoomLevel                             float64
	ai                                    *AI
	hudText                               []*tentsuyu.MenuElement
	timer                                 int
	tick                                  int
	gameStateMsg                          GameStateMsg
	lose, win                             bool
	pauseTick                             int
	prevTime                              int
	endGameTransition                     bool
	playerTurn, enemyTurn, startEnemyTurn bool
	playStarted                           bool
	playerCharacter, enemyCharacter       *Character
}

//NewGameBattle creates a GameBattle state
func NewGameBattle(g *Game) *GameBattle {
	x := 128.0
	y := 88.0
	//x := 206.0 * 2
	//y := 52.0 * 2
	ZoomLevel = 1.0
	gm := &GameBattle{
		playerGrid:    CreateGrid(x, y, 10, 9),
		enemyGrid:     CreateGrid(426, y, 10, 9),
		playerDisplay: NewEnemyDisplay(5, y, false),
		aiDisplay:     NewEnemyDisplay(724, y, false),
		ships:         []*Ship{},
		ai:            NewAI(),
		timer:         30,
		pauseTick:     1,
		win:           false,
		lose:          false,
		playerTurn:    true,
		playStarted:   false,
	}
	TimeRanOut = false
	AIBroke = false
	g.gameData.SetGameMode(GameModeNormalBattle)
	gm.ai.SetBoard(gm.playerGrid)
	gm.ai.SetBoard(gm.enemyGrid)
	gm.enemyGrid.MakeAllVisible()
	gm.enemyGrid.playable = false

	gm.playerCharacter = NewCharacter(g.gameData.playerCharacter)
	gm.playerCharacter.SetPosition(600, 14)
	if g.gameData.playerCharacter == nure {
		gm.playerCharacter.SetPosition(96, 0)
	}

	gm.enemyCharacter = NewCharacter(nure) //458,96
	//gs.enemyCharacter.SetPosition(458, 50)
	//if comm.Name == nure {
	gm.enemyCharacter.SetPosition(96, -8)
	//}

	gm.hudText = []*tentsuyu.MenuElement{
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(340, 432, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Player's Turn"}, color.Black, 18),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement: tentsuyu.NewTextElementStationary(426, 407, 300, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Click to Randomize", "(Click a monster to manually move it)"}, color.Black, 12),
			Action: func() {
				if !gm.playStarted {
					gm.ai.ResetBoard(gm.enemyGrid)
					gm.enemyGrid.MakeAllVisible()
					gm.enemyGrid.playable = false
				}
			},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement: tentsuyu.NewTextElementStationary(325, 30, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Click When Ready"}, color.RGBA{R: 124, G: 255, B: 0, A: 255}, 20),
			Action: func() {
				if g.player.holdingShip {
					return
				}
				gm.playStarted = true
				(gm.hudText[5].UIElement).(*tentsuyu.TextElement).SetText([]string{""})
				(gm.hudText[6].UIElement).(*tentsuyu.TextElement).SetText([]string{""})
			},
			Selectable: true,
		},
		/*&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(450, 432, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Turn"}, color.Black, 18),
			Action:     func() {},
			Selectable: false,
		},*/
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(128, 407, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Player Attack Grid"}, color.Black, 12),
			Action:     func() {},
			Selectable: false,
		},

		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(426, 407, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Nure Attack Grid"}, color.Black, 12),
			Action:     func() {},
			Selectable: false,
		},

		&tentsuyu.MenuElement{ //5
			UIElement:  tentsuyu.NewTextElementStationary(144, 220, 300, 40, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Click in this grid to attack", "Sink all ships to win"}, color.Black, 14),
			Action:     func() {},
			Selectable: false,
		},

		&tentsuyu.MenuElement{ // 6
			UIElement:  tentsuyu.NewTextElementStationary(430, 220, 300, 40, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Click to select and move pieces", "Right and Left Arrow rotates"}, color.Black, 14),
			Action:     func() {},
			Selectable: false,
		},
		/*&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&gm.timer, 160, 10, 80, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},*/
		/*&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(65, 40, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Monsters Remaining: "}, color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewUINumberDisplayIntStationary(&gm.playerGrid.shipsRemaining, 170, 40, 80, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), color.Black, 16),
			Action:     func() {},
			Selectable: false,
		},*/

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
	//gm.ships = append(gm.ships, CreateShip(x, y, ShipTypeCruiser, true))
	return gm
}

//Update the main game
func (gm *GameBattle) Update(g *Game) error {
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
	gm.playerDisplay.Update(gm.playerGrid.Ships)
	gm.aiDisplay.Update(gm.enemyGrid.Ships)
	if gm.endGameTransition == true {
		if gm.playerGrid.allVisible == false {
			gm.playerGrid.MakeAllVisible()
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
	if !gm.playStarted {
		gm.enemyGrid.UpdatePlacement(g)
	}
	g.player.Update()
	if gm.playStarted {
		if gm.playerTurn {

			if gm.playerGrid.Update(g) {
				gm.playerTurn = false
				gm.startEnemyTurn = true
				gm.prevTime = int(g.gameData.TimeInMilliseconds())
				gm.tick = 0
				(gm.hudText[0].UIElement).(*tentsuyu.TextElement).SetText([]string{"Nure's Turn"})
			}
		}
	}
	if gm.startEnemyTurn {
		gm.tick++
		if g.gameData.TimeInMilliseconds()-int64(gm.prevTime) >= 1500 {
			gm.startEnemyTurn = false
			gm.enemyTurn = true
			gm.prevTime = g.gameData.TimeInSecond()
		}
	}
	if gm.enemyTurn {
		if g.gameData.TimeInSecond()-gm.prevTime >= 10 {
			AIBroke = true
			gm.win = true
			gm.endGameTransition = true
		}
		if gm.ai.TakeShot(gm.enemyGrid) {
			gm.enemyTurn = false
			gm.playerTurn = true
			(gm.hudText[0].UIElement).(*tentsuyu.TextElement).SetText([]string{"Player's Turn"})
		}
	}

	if gm.playerGrid.shipsRemaining == 0 && gm.win == false && gm.enemyGrid.shipsRemaining > 0 {
		gm.tick = 0
		//gm.timer = 1
		gm.win = true
		gm.endGameTransition = true
	}

	if gm.enemyGrid.shipsRemaining == 0 && gm.lose == false {
		gm.lose = true
		gm.endGameTransition = true
	}
	return nil
}

//Draw the main game
func (gm *GameBattle) Draw(g *Game) error {
	g.DrawBackground()
	if gm.enemyGrid.prevHit {
		gm.playerCharacter.DrawBustSad(g.screen)
	} else {
		gm.playerCharacter.DrawBust(g.screen)
	}
	if gm.enemyCharacter != nil {
		if gm.playerGrid.prevHit {
			gm.enemyCharacter.DrawBustSad(g.screen)
		} else {
			gm.enemyCharacter.DrawBust(g.screen)
		}
	}
	gm.playerGrid.Draw(g.screen)
	gm.enemyGrid.Draw(g.screen)
	if g.player.heldShip != nil {
		g.player.heldShip.Draw(g.screen)
	}
	gm.playerDisplay.Draw(g.screen)
	gm.aiDisplay.Draw(g.screen)
	for i, text := range gm.hudText {
		if i != 1 || !gm.playStarted { // Hide the Reset Board option
			if i != 2 || !gm.playStarted { //Hide click to play text
				if i != 4 || gm.playStarted {
					if i != 0 || gm.playStarted {
						text.Draw(g.screen)
					}
				}
			}
		}
	}
	return nil
}

//Msg returns the current state's message
func (gm *GameBattle) Msg() GameStateMsg {
	return gm.gameStateMsg
}

//SetMsg sets the GameStateMsg
func (gm *GameBattle) SetMsg(msg GameStateMsg) {
	gm.gameStateMsg = msg
}
