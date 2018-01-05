package isledef

import (
	"log"

	"github.com/atolVerderben/spaloosh/network"
)

func HandleConnection(p *Player, g *Game) {
	for p.conn != nil {

		comm := network.Command{}
		err := p.wsDecoder.Decode(&comm)

		if err == nil {
			switch gs := g.gameState.(type) {
			case *MPStage:
				if comm.CommType == network.CommandConnected && comm.Row == 2 {
					gs.connected = true
				}
				if comm.CommType == network.CommandJoinRoom && comm.Row == 2 {
					gs.connected = true
				}
				if comm.CommType == network.CommandFull {
					gs.title.SetText([]string{"Server is Full"})
				}
				if comm.CommType == network.CommandDisconnected {
					log.Println("Server is probably full")
					gs.title.SetText([]string{"Server is Full or lost connection"})
					return
				}
			case *MPBattle:
				if comm.CommType == network.CommandAttack {
					gs.enemyGrid.AIShot(comm.Row, comm.Col)
					gs.enemyTurn = false
					gs.playerTurn = true
				}
				if comm.CommType == network.CommandYouStart { //Used for room connections
					gs.playerTurn = true
					gs.enemyTurn = false

				}
				if comm.CommType == network.CommandSetTheBoard {
					p.gameData.opponentCharacter = comm.Name
					gs.ai.NetworkSetOpponentBoard(gs.playerGrid, comm.ShipPlacements)
					//gs.enemyGrid.MakeAllVisible()
					//gs.enemyGrid.playable = false
					gs.enemyReady = true

				}
				if comm.CommType == network.CommandHello {
					if comm.Name != "" {
						p.gameData.opponentCharacter = comm.Name
						gs.enemyCharacter = NewCharacter(comm.Name) //458,96
						gs.enemyCharacter.SetPosition(128, 14)
						if comm.Name == nure {
							gs.enemyCharacter.SetPosition(96, -8)
						}
					}

				}
			case *MPGameOver:
				if comm.CommType == network.CommandDisconnected {
					log.Println("Other Player has left the game.")
					gs.LostConnection()
					return
				}
			case *MPRooms:
				if comm.CommType == network.CommandListRooms {
					gs.AddRooms(comm.Rooms, p)

				}
				if comm.CommType == network.CommandJoinRoom && comm.Row == 2 {
					gs.connected = true
				}
			}

			if comm.CommType == network.CommandDisconnected {
				log.Println("Other Player has left the game.")
				g.gameState.SetMsg(GameStateMsgReqLostConnection)
			}
		}

		if err != nil {
			log.Println(err.Error())
			//p.conn.Close()
			p.conn = nil
			switch gs := g.gameState.(type) {
			case *MPGameOver:
				if comm.CommType == network.CommandDisconnected {
					log.Println("Other Player has left the game.")
					gs.LostConnection()
					return
				}
			case *MPStage:
				if comm.CommType == network.CommandDisconnected {
					log.Println("Server is probably full")
					gs.title.SetText([]string{"Server is Full or Lost Connection"})
					return
				}
			default:
				g.gameState.SetMsg(GameStateMsgReqLostConnection)
			}
		}

	}
}
