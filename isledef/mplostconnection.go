package isledef

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
)

//LostConnection is displayed when the game is over
type LostConnection struct {
	gameStateMsg  tentsuyu.GameStateMsg
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

//CreateLostConnection creates either a winning screen or losing screen depending on the bool value
func CreateLostConnection() *LostConnection {

	t := &LostConnection{
		title: tentsuyu.NewTextElementStationary(250, 5, 600, 200, Game.UIController.ReturnFont(FntSmallPixel), []string{"Lost Connection to Opponent"}, color.White, 24),
	}

	testMenu := tentsuyu.NewMenu(ScreenWidth, ScreenHeight)
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 100, 30, Game.UIController.ReturnFont(FntSmallPixel), []string{"Main Menu"}, color.White, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMainMenu }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 100, 25, Game.UIController.ReturnFont(FntSmallPixel), []string{"Quit Game"}, color.White, 16)},
		[]func(){func() {
			os.Exit(0)
		}})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, Game.UIController.ReturnFont(FntSmallPixel), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
	testMenu.SetBackground(tentsuyu.ImageManager.ReturnImage("topbar-light"), &tentsuyu.BasicImageParts{
		Sx:     0,
		Sy:     0,
		Width:  100,
		Height: 40,
	})
	t.menu = testMenu
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("bgDark"), count: 20, width: 1920, height: 1080}
	*/
	t.background = &backgroundImageParts{image: Game.ImageManager.ReturnImage("blue"), count: 20, width: 1920, height: 1080}
	t.menu = testMenu

	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

//Update LostConnection screen
func (t *LostConnection) Update(game *tentsuyu.Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.menu.Update(game.Input, 0, 0)
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if game.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMainMenu
	}
	return nil
}

//Draw LostConnection scene
func (t *LostConnection) Draw(game *tentsuyu.Game) error {
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
	t.menu.Draw(game.Screen)
	t.title.Draw(game.Screen)

	return nil
}

//Msg returns the gamestate msg
func (t *LostConnection) Msg() tentsuyu.GameStateMsg {
	return t.gameStateMsg
}

func (t *LostConnection) SetMsg(msg tentsuyu.GameStateMsg) {
	t.gameStateMsg = msg
}
