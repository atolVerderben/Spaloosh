package isledef

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

var (
	ghostgirl = "ghostgirl"
	vampire   = "vampire"
	hunter    = "vampirehunter"
	nure      = "nure"
)

type MPMainMenu struct {
	gameStateMsg              GameStateMsg
	timer                     int
	offsetX                   int
	offsetY                   int
	title                     *tentsuyu.TextElement
	menu                      *tentsuyu.Menu
	background                *backgroundImageParts
	startButton               *tentsuyu.BasicObject
	currMenu                  string
	desc                      *tentsuyu.TextElement
	serverBoxInfo             *tentsuyu.TextElement
	portBoxInfo               *tentsuyu.TextElement
	serverBox                 *tentsuyu.TextBox
	portBox                   *tentsuyu.TextBox
	charOptions               []*tentsuyu.MenuElement
	selected                  string
	ghost, vamp, hunter, nure *tentsuyu.BasicObject
}

func CreateMPMainMenu(g *Game) *MPMainMenu {
	if g.player.conn != nil {
		g.player.conn.Close()
		g.player.conn = nil
	}

	if GameServer != nil {
		GameServer.ShutDown()
	}
	GameServer = nil
	g.player.Reset()
	tentsuyu.Components.Camera.SetZoom(2.0)
	t := &MPMainMenu{
		title: tentsuyu.NewTextElement(175, 5, 500, 20, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Choose your character and spaloosh with friends!"}, color.White, 16),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(25, 275, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Challenge your friends to a game of Spaloosh!",
				""}, color.Black, 16),
	}

	t.charOptions = []*tentsuyu.MenuElement{
		&tentsuyu.MenuElement{
			UIElement: tentsuyu.NewTextElementStationary(20, 150, 200, 40, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Tiffany", "the Ghost"}, color.Black, 18),
			Action: func() {
				t.selected = ghostgirl
			},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement: tentsuyu.NewTextElementStationary(645, 150, 200, 40, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Petar", "the Vampire"}, color.Black, 18),
			Action: func() {
				t.selected = vampire
			},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement: tentsuyu.NewTextElementStationary(20, 260, 200, 40, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Archinal", "the Hunter"}, color.Black, 18),
			Action: func() {
				t.selected = hunter
			},
			Selectable: true,
		},
		&tentsuyu.MenuElement{
			UIElement: tentsuyu.NewTextElementStationary(655, 260, 200, 40, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Nure", "Nure-Onna"}, color.Black, 18),
			Action: func() {
				t.selected = nure
			},
			Selectable: true,
		},
	}
	testMenu := tentsuyu.NewMenu()
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 300, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Host Multiplayer Game"}, color.Black, 16),
	},
		[]func(){
			func() {
				t.gameStateMsg = GameStateMsgReqSetIP //GameStateMsgReqMPStage
				g.gameData.SetGameMode(GameModeOnlineHost)
				//g.gameData.server = t.serverBox.Text.ReturnText()
				//g.gameData.port = t.portBox.Text.ReturnText()
				g.gameData.playerCharacter = t.selected
				if t.selected == "" {
					g.gameData.playerCharacter = nure
				}
			},
		})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 300, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Join Direct IP Game"}, color.Black, 16)},
		[]func(){func() {
			t.gameStateMsg = GameStateMsgReqSetIP //GameStateMsgReqMPStage
			g.gameData.SetGameMode(GameModeOnlineJoin)
			//g.gameData.server = t.serverBox.Text.ReturnText()
			//g.gameData.port = t.portBox.Text.ReturnText()
			g.gameData.playerCharacter = t.selected
			if t.selected == "" {
				g.gameData.playerCharacter = nure
			}
		}})
	testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 300, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Join Server Game"}, color.Black, 16)},
		[]func(){func() {
			t.gameStateMsg = GameStateMsgReqSetIP //GameStateMsgReqHostingRooms
			g.gameData.SetGameMode(GameModeOnlineRoom)
			g.gameData.playerCharacter = t.selected
			if t.selected == "" {
				g.gameData.playerCharacter = nure
			}
		}})
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 300, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Setup Help"}, color.Black, 16),
	},
		[]func(){
			func() {
				t.gameStateMsg = GameStateMsgReqMPHelp
			},
		})
	/*t.serverBoxInfo = tentsuyu.NewTextElement(275, 250, 300, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Server IP Address:"}, color.Black, 16)
	t.portBoxInfo = tentsuyu.NewTextElement(275, 300, 300, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Port Number:"}, color.Black, 16)
	t.serverBox = tentsuyu.NewTextBox(475, 250, 200, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{g.gameData.server}, color.Black, 16)
	t.portBox = tentsuyu.NewTextBox(475, 300, 100, 25, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{g.gameData.port}, color.Black, 16)
	*/
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MPMainMenu"
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

	t.ghost = &tentsuyu.BasicObject{
		X:           20,
		Y:           45,
		Width:       96,
		Height:      96,
		NotCentered: true,
	}

	t.vamp = &tentsuyu.BasicObject{
		X:           655,
		Y:           45,
		Width:       96,
		Height:      96,
		NotCentered: true,
	}

	t.hunter = &tentsuyu.BasicObject{
		X:           20,
		Y:           300,
		Width:       96,
		Height:      96,
		NotCentered: true,
	}

	t.nure = &tentsuyu.BasicObject{
		X:           615,
		Y:           290,
		Width:       SpalooshSheet.Frames[frameNureOnna].Frame["w"] * 3,
		Height:      SpalooshSheet.Frames[frameNureOnna].Frame["h"] * 3,
		NotCentered: true,
	}
	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPMainMenu) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++

	t.menu.Update()
	//t.serverBox.Update()
	//t.portBox.Update()
	for _, o := range t.charOptions {
		o.UIElement.(*tentsuyu.TextElement).UnHighlighted()
		o.Update()
	}
	if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.ghost.Contains(tx, ty) {
			t.selected = ghostgirl
		}
		if t.vamp.Contains(tx, ty) {
			t.selected = vampire
		}
		if t.hunter.Contains(tx, ty) {
			t.selected = hunter
		}
		if t.nure.Contains(tx, ty) {
			t.selected = nure
		}
	}
	if tentsuyu.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMainMenu
	}
	if tentsuyu.Input.Button("Enter").Down() {
		t.gameStateMsg = GameStateMsgReqMain
	}
	return nil
}

func (t *MPMainMenu) Draw(game *Game) error {
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

	//Nure
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameNureOnna].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameNureOnna].Frame["y"],
		Width:  SpalooshSheet.Frames[frameNureOnna].Frame["w"],
		Height: SpalooshSheet.Frames[frameNureOnna].Frame["h"],
	}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(t.nure.X, t.nure.Y)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	//Ghost Girl
	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameGhostGirl].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameGhostGirl].Frame["y"],
		Width:  SpalooshSheet.Frames[frameGhostGirl].Frame["w"],
		Height: SpalooshSheet.Frames[frameGhostGirl].Frame["h"],
	}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(t.ghost.X, t.ghost.Y)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	//Vampire
	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameVamp].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameVamp].Frame["y"],
		Width:  SpalooshSheet.Frames[frameVamp].Frame["w"],
		Height: SpalooshSheet.Frames[frameVamp].Frame["h"],
	}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(t.vamp.X, t.vamp.Y)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	//Vampire Hunter
	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameVampHunter].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameVampHunter].Frame["y"],
		Width:  SpalooshSheet.Frames[frameVampHunter].Frame["w"],
		Height: SpalooshSheet.Frames[frameVampHunter].Frame["h"],
	}
	op.GeoM.Scale(3, 3)
	op.GeoM.Translate(t.hunter.X, t.hunter.Y)
	if err := game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

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

	switch t.selected {
	case "ghostgirl":
		t.charOptions[0].UIElement.(*tentsuyu.TextElement).Highlighted()
	case "vampire":
		t.charOptions[1].UIElement.(*tentsuyu.TextElement).Highlighted()
	case "vampirehunter":
		t.charOptions[2].UIElement.(*tentsuyu.TextElement).Highlighted()
	case "nure":
		t.charOptions[3].UIElement.(*tentsuyu.TextElement).Highlighted()
	}
	for _, o := range t.charOptions {
		o.Draw(game.screen)
	}
	t.menu.Draw(game.screen)
	//t.serverBoxInfo.Draw(game.screen)
	//t.portBoxInfo.Draw(game.screen)
	//t.serverBox.Draw(game.screen)
	//t.portBox.Draw(game.screen)
	t.title.Draw(game.screen)
	//t.desc.Draw(game.screen)

	return nil
}

func (t *MPMainMenu) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *MPMainMenu) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
