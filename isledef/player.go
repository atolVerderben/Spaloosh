package isledef

import (
	"encoding/gob"
	"net"

	"github.com/atolVerderben/tentsuyu"
)

//Player represents the player
type Player struct {
	shots       int
	gameData    *tentsuyu.GameData
	holdingShip bool
	heldShip    *Ship
	conn        *net.TCPConn
	commEncoder *gob.Encoder
	wsDecoder   *gob.Decoder
}

//CreatePlayer returns a pointer to a Player
func CreatePlayer(g *tentsuyu.Game) *Player {
	p := &Player{
		shots:    g.GameData.Settings["ShotMax"].ValueInt,
		gameData: g.GameData,
	}
	return p
}

//Update the player
func (p *Player) Update() {
	if p.shots < 0 {
		p.shots = 0
	}
	if p.heldShip != nil {
		p.heldShip.Update()
	}
}

//TakeShot takes a shot and returns true if the player can take a shot. Otherwise return false
func (p *Player) TakeShot() bool {
	if p.gameData.Settings["GameMode"].ValueInt == GameModeNormalBattle ||
		p.gameData.Settings["GameMode"].ValueInt == GameModeOnlineHost ||
		p.gameData.Settings["GameMode"].ValueInt == GameModeOnlineJoin {
		return true
	}
	if p.shots > 0 {
		p.shots--
		return true
	}
	return false
}

//PlayerChoosePosition determines if the ship can fit in the spot given with the given orientation
func PlayerChoosePosition(col, row, length int, vertical bool, g *Grid) bool {
	//fmt.Printf("I'm starting at (%v,%v)\n", col, row)
	if row > g.Rows-1 || col > g.Columns-1 {
		return false
	}
	if g.area[row][col] != TileEmpty {
		return false
	}

	collide := false

	if !vertical {
		collide = checkHoriz(row, col, length, g)
	} else {
		collide = checkVert(row, col, length, g)
	}

	if collide == true {
		return false
	}

	return true
}

func (p *Player) Reset() {
	p.heldShip = nil
	p.holdingShip = false
	p.gameData.Settings["OpponentCharacter"].ValueText = ""
	//p.gameData.youStart = false
	p.gameData.Settings["JoinedRoom"].ValueText = ""
}
