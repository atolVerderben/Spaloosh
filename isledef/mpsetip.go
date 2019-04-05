package isledef

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/atolVerderben/tentsuyu"
)

type MPSetIP struct {
	gameStateMsg  tentsuyu.GameStateMsg
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

func CreateMPSetIP(g *tentsuyu.Game) *MPSetIP {

	g.DefaultCamera.SetZoom(2.0)

	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineRoom {
		con := readConfigFile("assets/config.json")
		g.GameData.Settings["Server"].ValueText = con.Server
		g.GameData.Settings["Port"].ValueText = con.Port
	}

	t := &MPSetIP{
		title: tentsuyu.NewTextElement(175, 5, 500, 20, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Enter IP Information"}, color.White, 16),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(25, 275, 1300, 400, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Challenge your friends to a game of Spaloosh!",
				""}, color.Black, 16),
	}
	buttonText := "Join Multiplayer Game"
	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineHost {
		buttonText = "Host Multiplayer Game"
	}
	testMenu := tentsuyu.NewMenu(ScreenWidth, ScreenHeight)
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 300, 25, g.UIController.ReturnFont(FntSmallPixel), []string{buttonText}, color.Black, 16)},
		[]func(){func() {

			g.GameData.Settings["Server"].ValueText = t.serverBox.Text.ReturnText()
			g.GameData.Settings["Port"].ValueText = t.portBox.Text.ReturnText()

			if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineRoom {
				conf := &config{
					Server: g.GameData.Settings["Server"].ValueText,
					Port:   g.GameData.Settings["Port"].ValueText,
				}

				c, _ := json.Marshal(conf)

				if err := ioutil.WriteFile("assets/config.json", c, 0644); err != nil {
					log.Printf("Error writing file: %v\n", err.Error())
				}
				t.gameStateMsg = GameStateMsgReqHostingRooms
			} else if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
				t.gameStateMsg = GameStateMsgReqMPStage
			} else if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineHost {
				t.gameStateMsg = GameStateMsgReqMPStage
				g.GameData.Settings["Port"].ValueText = t.portBox.Text.ReturnText()
			}

		}})
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 300, 25, g.UIController.ReturnFont(FntSmallPixel), []string{"Return"}, color.Black, 16),
	},
		[]func(){
			func() {
				t.gameStateMsg = GameStateMsgReqMPMainMenu
			},
		})
	t.serverBoxInfo = tentsuyu.NewTextElement(200, 250, 300, 50, g.UIController.ReturnFont(FntSmallPixel), []string{"Server IP Address:"}, color.Black, 16)
	t.portBoxInfo = tentsuyu.NewTextElement(275, 300, 300, 50, g.UIController.ReturnFont(FntSmallPixel), []string{"Port Number:"}, color.Black, 16)
	t.serverBox = tentsuyu.NewTextBox(400, 250, 400, 25, g.UIController.ReturnFont(FntSmallPixel),
		[]string{g.GameData.Settings["Server"].ValueText}, color.Black, 16)
	t.portBox = tentsuyu.NewTextBox(475, 300, 100, 25, g.UIController.ReturnFont(FntSmallPixel),
		[]string{g.GameData.Settings["Port"].ValueText}, color.Black, 16)
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, g.UIController.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MPSetIP"
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

	t.currMenu = "A"

	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPSetIP) Update(game *tentsuyu.Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.menu.Update(game.Input, 0, 0)
	if game.GameData.Settings["GameMode"].ValueInt != GameModeOnlineHost {
		t.serverBox.Update(game.Input)
	}
	t.portBox.Update(game.Input)

	if game.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMPMainMenu
	}
	return nil
}

func (t *MPSetIP) Draw(game *tentsuyu.Game) error {
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

	if err := game.Screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(600), float64(205))
	op.GeoM.Translate(20, 270)

	game.Screen.DrawImage(tentsuyu.ImageManager.ReturnImage("textBubble"), op)*/

	t.menu.Draw(game.Screen)
	if game.GameData.Settings["GameMode"].ValueInt != GameModeOnlineHost {
		t.serverBoxInfo.Draw(game.Screen)
		t.serverBox.Draw(game.Screen)
	}
	t.portBoxInfo.Draw(game.Screen)

	t.portBox.Draw(game.Screen)
	t.title.Draw(game.Screen)
	//t.desc.Draw(game.Screen)

	return nil
}

func (t *MPSetIP) Msg() tentsuyu.GameStateMsg {
	return t.gameStateMsg
}

func (t *MPSetIP) SetMsg(msg tentsuyu.GameStateMsg) {
	t.gameStateMsg = msg
}
