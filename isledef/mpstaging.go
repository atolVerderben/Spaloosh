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
	gameStateMsg tentsuyu.GameStateMsg
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

func CreateMPStage(g *tentsuyu.Game) *MPStage {
	g.DefaultCamera.SetZoom(2.0)
	t := &MPStage{
		title: tentsuyu.NewTextElement(300, 5, 400, 20, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Waiting for a player to join..."}, color.White, 16),
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
				"Or we can battle each other in a \"fair\" fight",
				"If you win you are free to leave, if I win...",
				"",
				"I get your blood!"}, color.Black, 16),
	}

	testMenu := tentsuyu.NewMenu(ScreenHeight, ScreenWidth)
	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
		testMenu.AddElement([]tentsuyu.UIElement{
			tentsuyu.NewTextElement(0, 0, 155, 50, g.UIController.ReturnFont(FntSmallPixel), []string{"Retry"}, color.Black, 16),
		},
			[]func(){
				func() {
					if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin && GamePlayer.conn == nil {
						SERVER := g.GameData.Settings["Server"].ValueText
						PORT := g.GameData.Settings["Port"].ValueText
						conn, _ := net.Dial("tcp", SERVER+":"+PORT)

						if conn != nil && g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
							GamePlayer.conn = conn.(*net.TCPConn)
						}
					}
					if GamePlayer.conn != nil && g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
						//GamePlayer.bufferReader = bufio.NewReader(GamePlayer.conn)
						//GamePlayer.commandDecoder = json.NewDecoder(GamePlayer.conn)
						GamePlayer.wsDecoder = gob.NewDecoder(GamePlayer.conn)
						GamePlayer.commEncoder = gob.NewEncoder(GamePlayer.conn)
						go HandleConnection(GamePlayer, g)
					}
				},
			})
	}
	testMenu.AddElement([]tentsuyu.UIElement{
		tentsuyu.NewTextElement(0, 0, 155, 50, g.UIController.ReturnFont(FntSmallPixel), []string{"Cancel"}, color.Black, 16),
	},
		[]func(){
			func() {
				t.gameStateMsg = GameStateMsgReqMPMainMenu
				//g.gameData.SetGameMode(GameModeOnlineHost)
			},
		})
	/*testMenu.AddElement([]tentsuyu.UIElement{tentsuyu.NewTextElement(0, 0, 200, 25, g.UIController.ReturnFont("font1"), []string{"Continue"}, color.Black, 24)},
		[]func(){func() {
			/*prevMenu = "MPStage"
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
	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineRoom {
		comm := &network.Command{
			CommType: network.CommandJoinRoom,
			Name:     g.GameData.Settings["JoinedRoom"].ValueText,
		}
		err := GamePlayer.commEncoder.Encode(comm)

		if err != nil {
			fmt.Println(err.Error())
			GamePlayer.conn.Close()
			GamePlayer.conn = nil
		}
	}
	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineHost {
		GameServer = network.CreateGameServer(g.GameData.Settings["Port"].ValueText)
		GameServer.Run()
	}
	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineHost {
		SERVER := "127.0.0.1"
		PORT := g.GameData.Settings["Port"].ValueText
		conn, _ := net.Dial("tcp", SERVER+":"+PORT)
		if conn != nil {
			GamePlayer.conn = conn.(*net.TCPConn)
		}
	}
	if GamePlayer.conn != nil && g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineHost {
		//GamePlayer.bufferReader = bufio.NewReader(GamePlayer.conn)
		//GamePlayer.commandDecoder = json.NewDecoder(GamePlayer.conn)
		GamePlayer.wsDecoder = gob.NewDecoder(GamePlayer.conn)
		GamePlayer.commEncoder = gob.NewEncoder(GamePlayer.conn)
		go HandleConnection(GamePlayer, g)
	}
	if g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin && GamePlayer.conn == nil {
		SERVER := g.GameData.Settings["Server"].ValueText
		PORT := g.GameData.Settings["Port"].ValueText
		conn, _ := net.Dial("tcp", SERVER+":"+PORT)

		if conn != nil && g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
			GamePlayer.conn = conn.(*net.TCPConn)
		}
	}
	if GamePlayer.conn != nil && g.GameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
		//GamePlayer.bufferReader = bufio.NewReader(GamePlayer.conn)
		//GamePlayer.commandDecoder = json.NewDecoder(GamePlayer.conn)
		GamePlayer.wsDecoder = gob.NewDecoder(GamePlayer.conn)
		GamePlayer.commEncoder = gob.NewEncoder(GamePlayer.conn)
		go HandleConnection(GamePlayer, g)
	}
	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPStage) Update(game *tentsuyu.Game) error {
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
	if game.Input.Button("Enter").Down() {
		t.gameStateMsg = GameStateMsgReqMain
	}

	if t.connected == true {
		t.gameStateMsg = GameStateMsgReqMPMain
	}

	return nil
}

func (t *MPStage) Draw(game *tentsuyu.Game) error {
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
	DrawBackground(game) //background.Draw(game.screen, true)
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

	t.menu.Draw(game.Screen)
	t.title.Draw(game.Screen)
	//t.desc.Draw(game.screen)

	return nil
}

func (t *MPStage) Msg() tentsuyu.GameStateMsg {
	return t.gameStateMsg
}

func (t *MPStage) SetMsg(msg tentsuyu.GameStateMsg) {
	t.gameStateMsg = msg
}
