package isledef

import (
	"time"

	"github.com/atolVerderben/tentsuyu"
)

//GameMode represents the difficulty/mode of the game
type GameMode int

//Currently available game modes
const (
	GameModeNormalTimed int = iota
	GameModeHardcoreTimed
	GameModeNormalBattle
	GameModeOnlineHost
	GameModeOnlineJoin
	GameModeOnlineRoom
)

//GameData contains data for the game
type GameData struct {
	time                               int
	shotMax                            int
	timed                              bool
	currentScale                       int
	startTime                          time.Time
	gameMode                           GameMode
	server, port                       string
	playerCharacter, opponentCharacter string
	youStart                           bool
	joinedRoom                         string
}

func InitGameData(g *tentsuyu.GameData, gameMode int) {
	g.Settings["Server"] = &tentsuyu.GameValuePair{
		Name:      "Server",
		ValueType: tentsuyu.GameValueText,
		ValueText: "127.0.0.1",
	}
	g.Settings["Port"] = &tentsuyu.GameValuePair{
		Name:      "Port",
		ValueType: tentsuyu.GameValueText,
		ValueText: "5555",
	}
	g.Settings["Scale"] = &tentsuyu.GameValuePair{
		Name:      "Scale",
		ValueType: tentsuyu.GameValueInt,
		ValueInt:  1,
	}
	g.Settings["PlayerCharacer"] = &tentsuyu.GameValuePair{
		Name:      "Player Character",
		ValueType: tentsuyu.GameValueText,
		ValueText: "",
	}

	g.Settings["OpponentCharacter"] = &tentsuyu.GameValuePair{
		Name:      "Opponent Character",
		ValueType: tentsuyu.GameValueText,
		ValueText: "",
	}
	g.Settings["JoinedRoom"] = &tentsuyu.GameValuePair{
		Name:      "Joined Room",
		ValueType: tentsuyu.GameValueText,
		ValueText: "",
	}

	g.Settings["GameMode"] = &tentsuyu.GameValuePair{
		Name:      "Game Mode",
		ValueType: tentsuyu.GameValueInt,
		ValueInt:  gameMode,
	}

	switch gameMode {
	case GameModeNormalTimed:
		g.Settings["ShotMax"] = &tentsuyu.GameValuePair{
			Name:      "ShotMax",
			ValueType: tentsuyu.GameValueInt,
			ValueInt:  24,
		}
	case GameModeHardcoreTimed:
		g.Settings["ShotMax"] = &tentsuyu.GameValuePair{
			Name:      "ShotMax",
			ValueType: tentsuyu.GameValueInt,
			ValueInt:  15,
		}

	}
}

//SetGameMode sets the game mode chosen by the character. In this case used for Timed mode to set number of bombs
func SetGameMode(g *tentsuyu.GameData, gameMode int) {
	switch gameMode {
	case GameModeNormalTimed:
		g.Settings["ShotMax"].ValueInt = 24
	case GameModeHardcoreTimed:
		g.Settings["ShotMax"].ValueInt = 15
	}
	g.Settings["GameMode"].ValueInt = gameMode
}
