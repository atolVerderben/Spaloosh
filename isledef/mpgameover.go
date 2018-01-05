package isledef

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//MPGameOver is displayed when the game is over
type MPGameOver struct {
	gameStateMsg      GameStateMsg
	timer             int
	offsetX           int
	offsetY           int
	title             *tentsuyu.TextElement
	menu              *tentsuyu.Menu
	background        *backgroundImageParts
	startButton       *tentsuyu.BasicObject
	currMenu          string
	win               bool
	text              *tentsuyu.TextElement
	opponentCharacter *Character
}

//CreateMPGameOver creates either a winning screen or losing screen depending on the bool value
func CreateMPGameOver(g *Game, win bool) *MPGameOver {

	t := &MPGameOver{
		title: tentsuyu.NewTextElementStationary(288, 5, 600, 200, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Game Over!", "You've lost Your Monsters!!"}, color.White, 16),
	}
	t.opponentCharacter = NewCharacter(g.gameData.opponentCharacter)
	t.opponentCharacter.SetPosition(267, 312)
	t.opponentCharacter.SetScale(5)
	if g.gameData.opponentCharacter == nure {
		t.opponentCharacter.imgPartsBust = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameNureWin].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameNureWin].Frame["y"],
			Width:  SpalooshSheet.Frames[frameNureWin].Frame["w"],
			Height: SpalooshSheet.Frames[frameNureWin].Frame["h"] / 2,
		}
		t.opponentCharacter.SetPosition(235, 280)
	}
	t.win = win
	if win {
		t.title = tentsuyu.NewTextElementStationary(288, 5, 725, 200, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Congratulations!", "All Monsters Defeated!"}, color.White, 16)
	}
	if win {
		t.text = tentsuyu.NewTextElement(160, 430, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Oh, that's no fair! Why don't we play again?"}, color.Black, 16)
	} else {
		t.text = tentsuyu.NewTextElement(160, 430, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Haha! Better luck next time! Care for another?"}, color.Black, 16)
	}
	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Play Again"}, color.White, 16)},
		[]func(){func() {
			if g.gameData.gameMode == GameModeOnlineHost {
				t.gameStateMsg = GameStateMsgReqMPMain
			}
			if g.gameData.gameMode == GameModeOnlineJoin {
				t.gameStateMsg = GameStateMsgReqMPMain
			}
		}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Main Menu"}, color.White, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMainMenu }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Multiplayer Menu"}, color.White, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMPMainMenu }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Quit Game"}, color.White, 16)},
		[]func(){func() {
			os.Exit(0)
		}})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
	testMenu.SetBackground(tentsuyu.ImageManager.ReturnImage("topbar-light"), &tentsuyu.BasicImageParts{
		Sx:     0,
		Sy:     0,
		Width:  100,
		Height: 40,
	})
	t.menu = testMenu
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("bgDark"), count: 20, width: 1920, height: 1080}
	*/
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("blue"), count: 20, width: 1920, height: 1080}
	t.menu = testMenu

	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPGameOver) LostConnection() {
	t.menu.Elements[0][0].Selectable = false
	t.title.SetText([]string{"Lost Connection to Opponent"})
	//(gm.hudText[2].UIElement).(*tentsuyu.TextElement).SetFontSize(12)
	//(gm.hudText[2].UIElement).(*tentsuyu.TextElement).SetColor(color.White)
}

//Update MPGameOver screen
func (t *MPGameOver) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++
	t.title.Update()
	t.menu.Update()
	t.text.Update()
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if tentsuyu.Input.Button("Escape").Down() {
		//os.Exit(0)
	}
	return nil
}

//Draw MPGameOver scene
func (t *MPGameOver) Draw(game *Game) error {
	game.pausedState.Draw(game)
	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 40)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("map"), op); err != nil {
		return err
	}
	op.GeoM.Translate(0, -40)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("topbar"), op); err != nil {
		return err
	}*/
	/*t.background.Draw(game.screen, false)
	t.menu.Draw(game.screen)
	t.title.Draw(game.screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.3, 0.3)
	op.GeoM.Translate(400, 400)

	tentsuyu.ApplyCameraTransform(op, true)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("shenanijam"), op); err != nil {
		return err
	}*/
	t.background.Draw(game.screen, false)

	op := &ebiten.DrawImageOptions{}
	if t.win {
		t.opponentCharacter.DrawBustSad(game.screen)
	} else {
		t.opponentCharacter.DrawBust(game.screen)
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(600), float64(100))
	op.GeoM.Translate(100, 425)
	game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("textBubble"), op)

	t.text.Draw(game.screen)
	t.menu.Draw(game.screen)
	t.title.Draw(game.screen)

	return nil
}

//Msg returns the gamestate msg
func (t *MPGameOver) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *MPGameOver) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
