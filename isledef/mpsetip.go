package isledef

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/atolVerderben/tentsuyu"
)

type MPSetIP struct {
	gameStateMsg  GameStateMsg
	timer         int
	offsetX       int
	offsetY       int
	title         *tentsuyu.TextElement
	menu          *tentsuyu.Menu
	background    *backgroundImageParts
	startButton   *tentsuyu.BasicObject
	currMenu      string
	desc          *tentsuyu.TextElement
	serverBoxInfo *tentsuyu.TextElement
	portBoxInfo   *tentsuyu.TextElement
	serverBox     *tentsuyu.TextBox
	portBox       *tentsuyu.TextBox
	selected      string
}

func CreateMPSetIP(g *Game) *MPSetIP {

	tentsuyu.Components.Camera.SetZoom(2.0)

	if g.gameData.gameMode == GameModeOnlineRoom {
		con := readConfigFile("assets/config.json")
		g.gameData.server = con.Server
		g.gameData.port = con.Port
	}

	t := &MPSetIP{
		title: tentsuyu.NewTextElement(175, 5, 500, 20, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Enter IP Information"}, color.White, 16),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(25, 275, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Challenge your friends to a game of Spaloosh!",
				""}, color.Black, 16),
	}
	buttonText := "Join Multiplayer Game"
	if g.gameData.gameMode == GameModeOnlineHost {
		buttonText = "Host Multiplayer Game"
	}
	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 300, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{buttonText}, color.Black, 16)},
		[]func(){func() {

			g.gameData.server = t.serverBox.Text.ReturnText()
			g.gameData.port = t.portBox.Text.ReturnText()

			if g.gameData.gameMode == GameModeOnlineRoom {
				conf := &config{
					Server: g.gameData.server,
					Port:   g.gameData.port,
				}

				c, _ := json.Marshal(conf)

				if err := ioutil.WriteFile("assets/config.json", c, 0644); err != nil {
					log.Printf("Error writing file: %v\n", err.Error())
				}
				t.gameStateMsg = GameStateMsgReqHostingRooms
			} else if g.gameData.gameMode == GameModeOnlineJoin {
				t.gameStateMsg = GameStateMsgReqMPStage
			} else if g.gameData.gameMode == GameModeOnlineHost {
				t.gameStateMsg = GameStateMsgReqMPStage
				g.gameData.port = t.portBox.Text.ReturnText()
			}

		}})
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 300, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Return"}, color.Black, 16),
	},
		[]func(){
			func() {
				t.gameStateMsg = GameStateMsgReqMPMainMenu
			},
		})
	t.serverBoxInfo = tentsuyu.NewTextElement(200, 250, 300, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Server IP Address:"}, color.Black, 16)
	t.portBoxInfo = tentsuyu.NewTextElement(275, 300, 300, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Port Number:"}, color.Black, 16)
	t.serverBox = tentsuyu.NewTextBox(400, 250, 400, 25, tentsuyu.Components.ReturnFont(FntSmallPixel),
		[]string{g.gameData.server}, color.Black, 16)
	t.portBox = tentsuyu.NewTextBox(475, 300, 100, 25, tentsuyu.Components.ReturnFont(FntSmallPixel),
		[]string{g.gameData.port}, color.Black, 16)
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MPSetIP"
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

	t.currMenu = "A"

	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPSetIP) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.menu.Update()
	if game.gameData.gameMode != GameModeOnlineHost {
		t.serverBox.Update()
	}
	t.portBox.Update()

	if tentsuyu.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMPMainMenu
	}
	return nil
}

func (t *MPSetIP) Draw(game *Game) error {
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
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameNureOnna].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameNureOnna].Frame["y"],
		Width:  SpalooshSheet.Frames[frameNureOnna].Frame["w"],
		Height: SpalooshSheet.Frames[frameNureOnna].Frame["h"],
	}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(615, 290)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(600), float64(205))
	op.GeoM.Translate(20, 270)

	game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("textBubble"), op)*/

	t.menu.Draw(game.screen)
	if game.gameData.gameMode != GameModeOnlineHost {
		t.serverBoxInfo.Draw(game.screen)
		t.serverBox.Draw(game.screen)
	}
	t.portBoxInfo.Draw(game.screen)

	t.portBox.Draw(game.screen)
	t.title.Draw(game.screen)
	//t.desc.Draw(game.screen)

	return nil
}

func (t *MPSetIP) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *MPSetIP) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
