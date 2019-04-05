package isledef

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//GameOver is displayed when the game is over
type GameOver struct {
	gameStateMsg tentsuyu.GameStateMsg
	timer        int
	offsetX      int
	offsetY      int
	title        *tentsuyu.TextElement
	menu         *tentsuyu.Menu
	background   *backgroundImageParts
	startButton  *tentsuyu.BasicObject
	currMenu     string
	win          bool
	text         *tentsuyu.TextElement
}

//CreateGameOver creates either a winning screen or losing screen depending on the bool value
func CreateGameOver(g *tentsuyu.Game, win bool) *GameOver {

	t := &GameOver{
		title: tentsuyu.NewTextElementStationary(288, 5, 600, 200, g.UIController.ReturnFont(FntSmallPixel), []string{"Game Over!", "You ran out of bombs!"}, color.White, 16),
	}
	t.win = win
	if TimeRanOut {
		t.title = tentsuyu.NewTextElementStationary(288, 5, 600, 200, g.UIController.ReturnFont(FntSmallPixel), []string{"Game Over!", "You ran out of time!"}, color.White, 16)
	}
	if g.GameData.Settings["GameMode"].ValueInt == GameModeNormalBattle {
		t.title = tentsuyu.NewTextElementStationary(288, 5, 600, 200, g.UIController.ReturnFont(FntSmallPixel), []string{"Game Over!", "Nure Sank Your Monsters!!"}, color.White, 16)
	}
	if win {
		t.title = tentsuyu.NewTextElementStationary(288, 5, 725, 200, g.UIController.ReturnFont(FntSmallPixel), []string{"Congratulations!", "All Monsters Defeated!"}, color.White, 16)
		if AIBroke {
			t.title.SetText([]string{"Congratulations!", "You broke the AI!"})
		}
	}
	if win {
		t.text = tentsuyu.NewTextElement(160, 430, 1300, 400, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Oh, that's no fair! Why don't we play again?"}, color.Black, 16)
		if AIBroke {
			t.text.SetText([]string{"Tell atol he tried to hard and broke my brain!",
				"It's rare though I promise! Please play again!"})
		}
	} else {
		t.text = tentsuyu.NewTextElement(160, 430, 1300, 400, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Stand Still! This will only hurt for a minute!"}, color.Black, 16)
	}
	testMenu := tentsuyu.NewMenu(ScreenHeight, ScreenWidth)
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 30, g.UIController.ReturnFont(FntSmallPixel), []string{"Play Again"}, color.White, 16)},
		[]func(){func() {
			if g.GameData.Settings["GameMode"].ValueInt == GameModeHardcoreTimed || g.GameData.Settings["GameMode"].ValueInt == GameModeNormalTimed {
				t.gameStateMsg = GameStateMsgReqMain
			}
			if g.GameData.Settings["GameMode"].ValueInt == GameModeNormalBattle {
				t.gameStateMsg = GameStateMsgReqBattle
			}
		}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 30, g.UIController.ReturnFont(FntSmallPixel), []string{"Main Menu"}, color.White, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMainMenu }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Quit Game"}, color.White, 16)},
		[]func(){func() {
			os.Exit(0)
		}})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
	testMenu.SetBackground(tentsuyu.ImageManager.ReturnImage("topbar-light"), &tentsuyu.BasicImageParts{
		Sx:     0,
		Sy:     0,
		Width:  100,
		Height: 40,
	})
	t.menu = testMenu
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("bgDark"), count: 20, width: 1920, height: 1080}
	*/
	t.background = &backgroundImageParts{image: g.ImageManager.ReturnImage("blue"), count: 20, width: 1920, height: 1080}
	t.menu = testMenu

	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

//Update GameOver screen
func (t *GameOver) Update(game *tentsuyu.Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++
	t.title.Update()
	t.menu.Update(game.Input, 0, 0)
	t.text.Update()
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if game.Input.Button("Escape").Down() {
		//os.Exit(0)
	}
	return nil
}

//Draw GameOver scene
func (t *GameOver) Draw(game *tentsuyu.Game) error {
	game.PausedState.Draw(game)
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
	t.background.Draw(game.Screen, false)

	op := &ebiten.DrawImageOptions{}
	if t.win {
		op.ImageParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameNureLost].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameNureLost].Frame["y"],
			Width:  SpalooshSheet.Frames[frameNureLost].Frame["w"],
			Height: SpalooshSheet.Frames[frameNureLost].Frame["h"],
		}
	} else {
		op.ImageParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameNureWin].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameNureWin].Frame["y"],
			Width:  SpalooshSheet.Frames[frameNureLost].Frame["w"],
			Height: SpalooshSheet.Frames[frameNureLost].Frame["h"],
		}
	}
	op.GeoM.Scale(5, 5)
	op.GeoM.Translate(235, 280)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.Screen.DrawImage(game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(600), float64(100))
	op.GeoM.Translate(100, 425)
	game.Screen.DrawImage(game.ImageManager.ReturnImage("textBubble"), op)

	t.text.Draw(game.Screen)
	t.menu.Draw(game.Screen)
	t.title.Draw(game.Screen)

	return nil
}

//Msg returns the gamestate msg
func (t *GameOver) Msg() tentsuyu.GameStateMsg {
	return t.gameStateMsg
}

func (t *GameOver) SetMsg(msg tentsuyu.GameStateMsg) {
	t.gameStateMsg = msg
}
