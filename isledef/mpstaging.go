package isledef

import (
	"encoding/gob"
	"fmt"
	"image/color"
	"net"

	"github.com/atolVerderben/spaloosh/network"
	"github.com/atolVerderben/tentsuyu"
)

var GameServer *network.GameServer

type MPStage struct {
	gameStateMsg GameStateMsg
	timer        int
	offsetX      int
	offsetY      int
	lunkerMode   bool
	title        *tentsuyu.TextElement
	menu         *tentsuyu.Menu
	background   *backgroundImageParts
	startButton  *tentsuyu.BasicObject
	currMenu     string
	desc         *tentsuyu.TextElement
	connected    bool
}

func CreateMPStage(g *Game) *MPStage {
	tentsuyu.Components.Camera.SetZoom(2.0)
	t := &MPStage{
		title: tentsuyu.NewTextElement(300, 5, 400, 20, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Waiting for a player to join..."}, color.White, 16),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(25, 275, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Oh! How rude to sneak up on me while I wash my hair!",
				"I am a Nure-Onna, but you can call me Nure.",
				"How about we play a little game to make up for scaring me?",
				"I'll even let you choose.",
				"",
				"Choice one is a timed match where I will hide 3 monsters,",
				"You have 30 seconds and limited bombs to defeat them all.",
				"",
				"Or we can battle each other in a \"fair\" fight",
				"If you win you are free to leave, if I win...",
				"",
				"I get your blood!"}, color.Black, 16),
	}

	testMenu := tentsuyu.NewMenu()
	if g.gameData.gameMode == GameModeOnlineJoin {
		testMenu.AddElement([]tentsuyu.UIElement{
			tentsuyu.NewTextElement(0, 0, 155, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Retry"}, color.Black, 16),
		},
			[]func(){
				func() {
					if g.gameData.gameMode == GameModeOnlineJoin && g.player.conn == nil {
						SERVER := g.gameData.server
						PORT := g.gameData.port
						conn, _ := net.Dial("tcp", SERVER+":"+PORT)

						if conn != nil && g.gameData.gameMode == GameModeOnlineJoin {
							g.player.conn = conn.(*net.TCPConn)
						}
					}
					if g.player.conn != nil && g.gameData.gameMode == GameModeOnlineJoin {
						//g.player.bufferReader = bufio.NewReader(g.player.conn)
						//g.player.commandDecoder = json.NewDecoder(g.player.conn)
						g.player.wsDecoder = gob.NewDecoder(g.player.conn)
						g.player.commEncoder = gob.NewEncoder(g.player.conn)
						go HandleConnection(g.player, g)
					}
				},
			})
	}
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 155, 50, tentsuyu.Components.ReturnFont(FntSmallPixel), []string{"Cancel"}, color.Black, 16),
	},
		[]func(){
			func() {
				t.gameStateMsg = GameStateMsgReqMPMainMenu
				//g.gameData.SetGameMode(GameModeOnlineHost)
			},
		})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, tentsuyu.Components.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MPStage"
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
	if g.gameData.gameMode == GameModeOnlineRoom {
		comm := &network.Command{
			CommType: network.CommandJoinRoom,
			Name:     g.gameData.joinedRoom,
		}
		err := g.player.commEncoder.Encode(comm)

		if err != nil {
			fmt.Println(err.Error())
			g.player.conn.Close()
			g.player.conn = nil
		}
	}
	if g.gameData.gameMode == GameModeOnlineHost {
		GameServer = network.CreateGameServer(g.gameData.port)
		GameServer.Run()
	}
	if g.gameData.gameMode == GameModeOnlineHost {
		SERVER := "127.0.0.1"
		PORT := g.gameData.port
		conn, _ := net.Dial("tcp", SERVER+":"+PORT)
		if conn != nil {
			g.player.conn = conn.(*net.TCPConn)
		}
	}
	if g.player.conn != nil && g.gameData.gameMode == GameModeOnlineHost {
		//g.player.bufferReader = bufio.NewReader(g.player.conn)
		//g.player.commandDecoder = json.NewDecoder(g.player.conn)
		g.player.wsDecoder = gob.NewDecoder(g.player.conn)
		g.player.commEncoder = gob.NewEncoder(g.player.conn)
		go HandleConnection(g.player, g)
	}
	if g.gameData.gameMode == GameModeOnlineJoin && g.player.conn == nil {
		SERVER := g.gameData.server
		PORT := g.gameData.port
		conn, _ := net.Dial("tcp", SERVER+":"+PORT)

		if conn != nil && g.gameData.gameMode == GameModeOnlineJoin {
			g.player.conn = conn.(*net.TCPConn)
		}
	}
	if g.player.conn != nil && g.gameData.gameMode == GameModeOnlineJoin {
		//g.player.bufferReader = bufio.NewReader(g.player.conn)
		//g.player.commandDecoder = json.NewDecoder(g.player.conn)
		g.player.wsDecoder = gob.NewDecoder(g.player.conn)
		g.player.commEncoder = gob.NewEncoder(g.player.conn)
		go HandleConnection(g.player, g)
	}
	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPStage) Update(game *Game) error {
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
	if tentsuyu.Input.Button("Enter").Down() {
		t.gameStateMsg = GameStateMsgReqMain
	}

	if t.connected == true {
		t.gameStateMsg = GameStateMsgReqMPMain
	}

	return nil
}

func (t *MPStage) Draw(game *Game) error {
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
	t.title.Draw(game.screen)
	//t.desc.Draw(game.screen)

	return nil
}

func (t *MPStage) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *MPStage) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
