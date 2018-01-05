package isledef

import (
	"fmt"
	"image/color"

	"time"

	"github.com/atolVerderben/spaloosh/network"
	"github.com/atolVerderben/tentsuyu"
)

//MPBattle is the main game...
type MPBattle struct {
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
	ready, enemyReady                     bool
	opponentShotRow, opponentShotCol      int
	playerCharacter, enemyCharacter       *Character
}

//NewMPBattle creates a MPBattle state
func NewMPBattle(g *Game) *MPBattle {
	x := 128.0
	y := 88.0
	//x := 206.0 * 2
	//y := 52.0 * 2
	ZoomLevel = 1.0
	gm := &MPBattle{
		playerGrid:      CreateGrid(x, y, 10, 9),
		enemyGrid:       CreateGrid(426, y, 10, 9),
		playerDisplay:   NewEnemyDisplay(5, y, false),
		aiDisplay:       NewEnemyDisplay(724, y, false),
		ships:           []*Ship{},
		ai:              NewAI(),
		timer:           30,
		pauseTick:       1,
		win:             false,
		lose:            false,
		playerTurn:      false,
		playStarted:     false,
		ready:           false,
		enemyReady:      false,
		opponentShotCol: -1,
		opponentShotRow: -1,
	}
	if g.gameData.gameMode == GameModeOnlineHost {
		gm.playerTurn = true
	} else {
		gm.playerTurn = false
	}
	TimeRanOut = false
	AIBroke = false
	gm.playerCharacter = NewCharacter(g.gameData.playerCharacter)
	gm.playerCharacter.SetPosition(600, 14)
	if g.gameData.playerCharacter == nure {
		gm.playerCharacter.SetPosition(588, -8)
	}

	if g.gameData.opponentCharacter != "" {
		gm.enemyCharacter = NewCharacter(g.gameData.opponentCharacter) //458,96
		gm.enemyCharacter.SetPosition(128, 14)
		if g.gameData.opponentCharacter == nure {
			gm.enemyCharacter.SetPosition(96, -8)
		}
	}
	/*go func() {
		timer := time.NewTicker(50 * time.Millisecond)
		for now := range timer.C {
			// entity updates
			// this is called every 100 millisecondes
			gm.SendHello(g.player, now)
		}
	}()*/

	gm.ai.SetBoard(gm.enemyGrid)
	gm.enemyGrid.MakeAllVisible()
	gm.enemyGrid.playable = false
	startText := "Your Turn"
	if g.gameData.gameMode == GameModeOnlineJoin || g.gameData.gameMode == GameModeOnlineRoom {
		startText = "Opponent's Turn"
	}
	gm.hudText = []*tentsuyu.MenuElement{
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(340, 432, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{startText}, color.Black, 18),
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
				if g.player.holdingShip || gm.ready {
					return
				}
				gm.ready = true
				(gm.hudText[2].UIElement).(*tentsuyu.TextElement).SetText([]string{"Waiting for Other Player..."})
				gm.SendShipPlacements(g.player, gm.enemyGrid.ExportNetworkPlacements())
				(gm.hudText[2].UIElement).(*tentsuyu.TextElement).SetFontSize(12)
				(gm.hudText[2].UIElement).(*tentsuyu.TextElement).SetColor(color.White)
				gm.hudText[2].Selectable = false
				(gm.hudText[5].UIElement).(*tentsuyu.TextElement).SetText([]string{""})
				(gm.hudText[6].UIElement).(*tentsuyu.TextElement).SetText([]string{""})
			},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(128, 407, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Player Attack Grid"}, color.Black, 12),
			Action:     func() {},
			Selectable: false,
		},
		&tentsuyu.MenuElement{
			UIElement:  tentsuyu.NewTextElementStationary(426, 407, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Opponent Attack Grid"}, color.Black, 12),
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
func (gm *MPBattle) Update(g *Game) error {
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
			gm.SendHello(g.player, time.Now())
			//shipPlacements := gm.ai.NetworkSetMyBoard(gm.playerGrid)
			//gm.SendShipPlacements(g.player, shipPlacements)

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
				gm.gameStateMsg = GameStateMsgReqMPGameOverWin
			} else if gm.lose == true {
				gm.gameStateMsg = GameStateMsgReqMPGameOverLose
			}
		}
		return nil
	}
	if gm.ready && gm.enemyReady {
		gm.playStarted = true
	}
	if !gm.ready {
		gm.enemyGrid.UpdatePlacement(g)
	}
	g.player.Update()
	if gm.playStarted {
		if gm.playerTurn {
			(gm.hudText[0].UIElement).(*tentsuyu.TextElement).SetText([]string{"Your Turn"})
			valid, row, col := gm.playerGrid.MPUpdate(g)
			if valid {
				if gm.playStarted == false {
					gm.playStarted = true
				}
				gm.SendMessage(g.player, row, col)
				gm.playerTurn = false
				gm.enemyTurn = true
				gm.prevTime = int(g.gameData.TimeInMilliseconds())
				gm.tick = 0

			}
		}
	}
	if gm.enemyTurn {
		(gm.hudText[0].UIElement).(*tentsuyu.TextElement).SetText([]string{"Opponent's Turn"})
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
func (gm *MPBattle) Draw(g *Game) error {
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
func (gm *MPBattle) Msg() GameStateMsg {
	return gm.gameStateMsg
}

//SetMsg sets the GameStateMsg
func (gm *MPBattle) SetMsg(msg GameStateMsg) {
	gm.gameStateMsg = msg
}

func (gm *MPBattle) SendHello(p *Player, now time.Time) {
	comm := &network.Command{
		CommType: network.CommandHello,
		Name:     p.gameData.playerCharacter,
	}
	err := p.commEncoder.Encode(comm)

	if err != nil {
		fmt.Println(err.Error())
		p.conn.Close()
		p.conn = nil
	}

}

func (gm *MPBattle) SendMessage(p *Player, row, col int) {
	comm := &network.Command{
		CommType: network.CommandAttack,
		Row:      row,
		Col:      col,
	}
	err := p.commEncoder.Encode(comm)

	if err != nil {
		fmt.Println(err.Error())
		p.conn.Close()
		p.conn = nil
	}

}

func (gm *MPBattle) SendShipPlacements(p *Player, shipPlacements []*network.PlaceShipType) {
	comm := &network.Command{
		CommType:       network.CommandSetTheBoard,
		ShipPlacements: shipPlacements,
		Name:           p.gameData.playerCharacter,
	}
	err := p.commEncoder.Encode(comm)

	if err != nil {
		fmt.Println(err.Error())
		p.conn.Close()
		p.conn = nil
	}

}
