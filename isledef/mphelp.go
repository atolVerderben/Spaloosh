package isledef

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

type MPHelp struct {
	gameStateMsg GameStateMsg
	title        *tentsuyu.TextElement
	desc         *tentsuyu.TextElement
}

func CreateMPHelp(g *Game) *MPHelp {
	if g.player.conn != nil {
		g.player.conn.Close()
		g.player.conn = nil
	}

	if GameServer != nil {
		GameServer.ShutDown()
	}
	GameServer = nil

	tentsuyu.Components.Camera.SetZoom(2.0)
	t := &MPHelp{
		title: tentsuyu.NewTextElement(300, 20, 300, 30, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Multiplayer Help"}, color.White, 24),
		//[]string{"Test of a Nure-Onna", "(A.K.A. Spaloosh!)"}, color.White, 8),
		desc: tentsuyu.NewTextElement(85, 100, 1300, 400, tentsuyu.Components.ReturnFont(FntSmallPixel),
			[]string{"Hosting:",
				"If players are having difficulty joining your hosted sessions",
				"It may be necessary to apply port forwarding rules to your router.",
				"In your router configuration, typically reached by opening your",
				"browser and entering an address similar to 192.168.1.1 or 192.168.0.1",
				"(this varies by router), look for port forwarding options and",
				"add a new rule for Spaloosh. Point incoming internet traffic",
				"from your selected port (default 5555) to your local computer's",
				"local IP address and apply the rule",
				"",
				"Joining:",
				"You may either join a direct IP game, or a server game.",
				"After selecting one both will ask for an IP to be entered.",
				"Direct IP the person you're playing against will have to provide their",
				"IP address and port they used when selecting to Host a Game.",
				"Server games will save the last used IP in a config file.",
			}, color.Black, 16),
	}

	return t
}

func (t *MPHelp) Update(game *Game) error {
	if t.gameStateMsg == GameStateMsgReqMain {
		return nil
	}
	if tentsuyu.Input.Button("Escape").JustPressed() {
		t.gameStateMsg = GameStateMsgReqMPMainMenu
	}
	return nil
}

func (t *MPHelp) Draw(game *Game) error {
	game.DrawBackground()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(700), float64(325))
	op.GeoM.Translate(80, 100)

	game.screen.DrawImage(tentsuyu.ImageManager.ReturnImage("textBubble"), op)

	t.title.Draw(game.screen)
	t.desc.Draw(game.screen)

	return nil
}

func (t *MPHelp) Msg() GameStateMsg {
	return t.gameStateMsg
}

func (t *MPHelp) SetMsg(msg GameStateMsg) {
	t.gameStateMsg = msg
}
