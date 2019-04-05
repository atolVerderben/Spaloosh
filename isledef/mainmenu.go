package isledef

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

type MainMenu struct {
	gameStateMsg tentsuyu.GameStateMsg
	timer        int
	offsetX      int
	offsetY      int
	title        *tentsuyu.TextElement
	menu         *tentsuyu.Menu
	background   *backgroundImageParts
	startButton  *tentsuyu.BasicObject
	currMenu     string
	desc         *tentsuyu.TextElement
}

func CreateMainMenu(g *tentsuyu.Game) *MainMenu {
	if GamePlayer.conn != nil {
		GamePlayer.conn.Close()
		GamePlayer.conn = nil
	}

	if GameServer != nil {
		GameServer.ShutDown()
	}
	GameServer = nil

	GamePlayer.Reset()

	g.DefaultCamera.SetZoom(2.0)
	t := &MainMenu{
		title: tentsuyu.NewTextElement(100, 5, 100, 20, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"SPALOOSH!"}, color.White, 16),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(25, 275, 1300, 400, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Oh! How rude to sneak up on me while I wash my hair!",
				"I am a Nure-Onna, but you can call me Nure.",
				"How about we play a little game to make up for scaring me?",
				"I'll even let you choose.",
				"",
				"Choice one is a timed match where I will hide 3 monsters,",
				"You have 30 seconds and limited bombs to defeat them all.",
				"",
				"Or we can battle each other in a \"fair\" fight with 4 monsters.",
				"If you win you are free to leave, if I win...",
				"",
				"I get your blood!"}, color.Black, 16),
	}

	testMenu := tentsuyu.NewMenu(ScreenWidth, ScreenHeight)
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 155, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Play Timed Game:"}, color.Black, 16),
		tentsuyu.NewTextElement(0, 0, 70, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Normal"}, color.Black, 16),
		tentsuyu.NewTextElement(0, 0, 50, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Hard"}, color.Black, 16)},
		[]func(){nil,
			func() {
				t.gameStateMsg = GameStateMsgReqMain
				SetGameMode(g.GameData, GameModeNormalTimed)
			},
			func() {
				t.gameStateMsg = GameStateMsgReqMain
				SetGameMode(g.GameData, GameModeHardcoreTimed)
			}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 300, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Play Monster Battle"}, color.Black, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqBattleCharacterSelect }})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 300, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Spaloosh With Friends  (Online)"}, color.Black, 16)},
		[]func(){func() { t.gameStateMsg = GameStateMsgReqMPMainMenu }})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, g.UIController.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MainMenu"
			BuildStatsMenu()
			g.UIController.UIController.ActivateMenu("StatMenu")
			g.UIController.UIController.DeActivateMenu(prevMenu)
		}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, g.UIController.ReturnFont("font1"), []string{"Quit"}, color.Black, 24)}, []func(){func() { os.Exit(0) }})
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

func (t *MainMenu) Update(game *tentsuyu.Game) error {
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
		t.gameStateMsg = GameStateMsgReqTitle
	}
	if game.Input.Button("Enter").Down() {
		t.gameStateMsg = GameStateMsgReqMain
	}
	return nil
}

func (t *MainMenu) Draw(game *tentsuyu.Game) error {
	/*op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 40)
	if err := game.Screen.DrawImage(tentsuyu.ImageManager.ReturnImage("map"), op); err != nil {
		return err
	}
	op.GeoM.Translate(0, -40)
	if err := game.Screen.DrawImage(tentsuyu.ImageManager.ReturnImage("topbar"), op); err != nil {
		return err
	}*/
	/*t.background.Draw(game.Screen, false)
	t.menu.Draw(game.Screen)
	t.title.Draw(game.Screen)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.3, 0.3)
	op.GeoM.Translate(400, 400)

	tentsuyu.ApplyCameraTransform(op, true)
	if err := game.Screen.DrawImage(tentsuyu.ImageManager.ReturnImage("shenanijam"), op); err != nil {
		return err
	}*/
	DrawBackground(game) //background.Draw(game.Screen, true)
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameNureOnna].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameNureOnna].Frame["y"],
		Width:  SpalooshSheet.Frames[frameNureOnna].Frame["w"],
		Height: SpalooshSheet.Frames[frameNureOnna].Frame["h"],
	}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(615, 290)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.Screen.DrawImage(game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(600), float64(205))
	op.GeoM.Translate(20, 270)

	game.Screen.DrawImage(game.ImageManager.ReturnImage("textBubble"), op)

	t.menu.Draw(game.Screen)
	//t.title.Draw(game.Screen)
	t.desc.Draw(game.Screen)

	return nil
}

func (t *MainMenu) Msg() tentsuyu.GameStateMsg {
	return t.gameStateMsg
}

func (t *MainMenu) SetMsg(msg tentsuyu.GameStateMsg) {
	t.gameStateMsg = msg
}
