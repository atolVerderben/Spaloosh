package isledef

import (
	"image/color"
	"os"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

type TitleMain struct {
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
	desc          *tentsuyu.TextElement
	start         *tentsuyu.MenuElement
}

func CreateTitleMain() *TitleMain {
	tentsuyu.Components.Camera.SetZoom(2.0)
	t := &TitleMain{
		title: tentsuyu.NewTextElement(400, 25, 100, 20, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"SPALOOSH!"}, color.Black, 16),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(25, 150, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Oh! I was just washing my hair human.",
				"I am a Nure-Onna, we enjoy our privacy.",
				"I will overlook this if you can win my game.",
				"Find my 3 friends in the depths to live.",
				"You have 30 seconds, and 24 shots.",
				"Make them count... or you are mine!"}, color.White, 16),
	}

	t.start = &tentsuyu.MenuElement{
		UIElement:  tentsuyu.NewTextElementStationary(200, 325, 200, 30, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Start Game"}, color.Black, 24),
		Action:     func() { t.gameStateMsg = GameStateMsgReqMainMenu },
		Selectable: true,
	}

	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 150, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Start Game"}, color.Black, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMainMenu }})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MainMenu"
			BuildStatsMenu()
			tentsuyu.Components.UIController.ActivateMenu("StatMenu")
			tentsuyu.Components.UIController.DeActivateMenu(prevMenu)
		}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
	testMenu.SetBackground(tentsuyu.ImageManager.ReturnImage("topbar-light"), &tentsuyu.BasicImageParts{
		Sx:     0,
		Sy:     0,
		Width:  100,
		Height: 40,
	})
	t.menu = testMenu
	t.background = &backgroundImageParts{image: tentsuyu.ImageManager.ReturnImage("bgDark"), count: 20, width: 1920, height: 1080}
	*/
	t.menu = testMenu
	t.startButton = &tentsuyu.BasicObject{
		X:           518,
		Y:           578,
		Width:       294,
		Height:      68,
		NotCentered: true,
	}
	t.currMenu = "A"
	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *TitleMain) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.start.Update()
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/
	if tentsuyu.Input.Button("Escape").JustPressed() {
		os.Exit(0)
	}
	if tentsuyu.Input.Button("Enter").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMain
	}
	return nil
}

func (t *TitleMain) Draw(game *Game) error {
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
	game.DrawBackground() //background.Draw(game.screen, true)
	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 5)
	op.GeoM.Translate(90, 80)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("girl"), op); err != nil {
		return err
	}*/

	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameSpalooshLogo].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameSpalooshLogo].Frame["y"],
		Width:  SpalooshSheet.Frames[frameSpalooshLogo].Frame["w"],
		Height: SpalooshSheet.Frames[frameSpalooshLogo].Frame["h"],
	}
	//op.GeoM.Scale(2, 2)
	//op.GeoM.Translate(250, 20)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameKraken].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameKraken].Frame["y"],
		Width:  SpalooshSheet.Frames[frameKraken].Frame["w"],
		Height: SpalooshSheet.Frames[frameKraken].Frame["h"],
	}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(620, 175)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	t.start.Draw(game.screen)
	//t.title.Draw(game.screen)
	//t.desc.Draw(game.screen)

	return nil
}

func (t *TitleMain) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *TitleMain) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
