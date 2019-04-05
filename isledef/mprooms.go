package isledef

import (
	"encoding/gob"
	"fmt"
	"image/color"
	"log"
	"net"
	"strconv"

	"github.com/atolVerderben/spaloosh/network"
	"github.com/atolVerderben/tentsuyu"
)

type MPRooms struct {
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
	connected    bool
	rooms        []*tentsuyu.MenuElement
	prevTime     int
}

func CreateMPRooms(g *tentsuyu.Game) *MPRooms {
	g.DefaultCamera.SetZoom(2.0)
	t := &MPRooms{
		title: tentsuyu.NewTextElement(300, 5, 400, 20, g.UIController.ReturnFont(FntSmallPixel),
			[]string{"Choose a room to play"}, color.White, 16),
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

	SERVER := g.GameData.Settings["Server"].ValueText
	PORT := g.GameData.Settings["Port"].ValueText
	conn, _ := net.Dial("tcp", SERVER+":"+PORT)

	if conn != nil {
		GamePlayer.conn = conn.(*net.TCPConn)
	}

	if GamePlayer.conn != nil {
		//GamePlayer.bufferReader = bufio.NewReader(GamePlayer.conn)
		//GamePlayer.commandDecoder = json.NewDecoder(GamePlayer.conn)
		GamePlayer.wsDecoder = gob.NewDecoder(GamePlayer.conn)
		GamePlayer.commEncoder = gob.NewEncoder(GamePlayer.conn)
		go HandleConnection(GamePlayer, g)
		comm := &network.Command{
			CommType: network.CommandListRooms,
		}
		err := GamePlayer.commEncoder.Encode(comm)

		if err != nil {
			fmt.Println(err.Error())
			GamePlayer.conn.Close()
			GamePlayer.conn = nil
		}
		t.prevTime = g.GameData.TimeInSecond()
	}
	g.GameData.Settings["GameMode"].ValueInt = GameModeOnlineRoom
	//tentsuyu.SetCustomCursor(30, 30, 30, 482, tentsuyu.ImageManager.ReturnImage("uiSheet"))
	return t
}

func init() {
	//rand.Seed(time.Now().UnixNano())

}

func (t *MPRooms) Update(game *tentsuyu.Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	t.timer++
	for _, r := range t.rooms {
		r.Update(game.Input, 0, 0)
	}
	/*if tentsuyu.Input.LeftClick().JustReleased() {
		tx, ty := tentsuyu.Input.GetMouseCoords()

		if t.startButton.Contains(tx, ty) {
			t.gameStateMsg = GameStateMsgReqMain
		}
	}*/

	if game.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMPMainMenu
	}

	if t.connected == true {
		//t.gameStateMsg = GameStateMsgReqMPMain
	}
	if game.GameData.TimeInSecond()-t.prevTime >= 10 {
		if GamePlayer.conn != nil {
			comm := &network.Command{
				CommType: network.CommandListRooms,
			}
			err := GamePlayer.commEncoder.Encode(comm)

			if err != nil {
				fmt.Println(err.Error())
				GamePlayer.conn.Close()
				GamePlayer.conn = nil
			}
		}
		t.prevTime = game.GameData.TimeInSecond()
	}

	return nil
}

//AddRooms returns all the available rooms from the server and number of people in each one
func (t *MPRooms) AddRooms(rooms []int, p *Player) {
	t.rooms = nil
	log.Println("Load rooms")
	startY := 75.0
	startX := 300.0
	for i := range rooms {
		num := i + 1
		name := "Room " + strconv.Itoa(num)
		me := &tentsuyu.MenuElement{

			UIElement: tentsuyu.NewTextElementStationary(startX, startY, 200, 30, Game.UIController.ReturnFont(FntSmallPixel), []string{name + "   " + strconv.Itoa(rooms[i])}, color.Black, 18),
			Action: func() {
				p.gameData.Settings["JoinedRoom"].ValueText = name
				t.gameStateMsg = GameStateMsgReqMPStage
			},
			Selectable: true,
		}
		t.rooms = append(t.rooms, me)
		startY += 35
	}

	me := &tentsuyu.MenuElement{

		UIElement: tentsuyu.NewTextElementStationary(startX, startY, 200, 30, Game.UIController.ReturnFont(FntSmallPixel), []string{"Cancel"}, color.Black, 18),
		Action: func() {
			t.gameStateMsg = GameStateMsgReqMPMainMenu
		},
		Selectable: true,
	}
	t.rooms = append(t.rooms, me)

}

func (t *MPRooms) Draw(game *tentsuyu.Game) error {
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
	for _, r := range t.rooms {
		r.Draw(game.Screen)
	}
	t.title.Draw(game.Screen)
	//t.desc.Draw(game.screen)

	return nil
}

func (t *MPRooms) Msg() tentsuyu.GameStateMsg {
	return t.gameStateMsg
}

func (t *MPRooms) SetMsg(msg tentsuyu.GameStateMsg) {
	t.gameStateMsg = msg
}
