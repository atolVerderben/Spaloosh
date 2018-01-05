package isledef

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
)

//Paused is displayed when the game is over
type Paused struct {
	gameStateMsg  GameStateMsg
	timer         int
	offsetX       int
	offsetY       int
	lunkerMode    bool
	lunkerCommand int
	title         *tentsuyu.TextElement
	menu          *tentsuyu.Menu
	background    *backgroundImageParts
	startButton   *tentsuyu.BasicObject
	currMenu      string
}

//CreatePaused creates either a winning screen or losing screen depending on the bool value
func CreatePaused() *Paused {

	t := &Paused{
		title: tentsuyu.NewTextElementStationary(340, 5, 600, 200, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Game Paused"}, color.White, 24),
	}

	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 100, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Resume"}, color.White, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgUnPause }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 100, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Main Menu"}, color.White, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMainMenu }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 100, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Quit Game"}, color.White, 16)},
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

//Update Paused screen
func (t *Paused) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.menu.Update()
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if tentsuyu.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgUnPause
	}
	return nil
}

//Draw Paused scene
func (t *Paused) Draw(game *Game) error {
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
	t.menu.Draw(game.screen)
	t.title.Draw(game.screen)

	return nil
}

//Msg returns the gamestate msg
func (t *Paused) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *Paused) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
