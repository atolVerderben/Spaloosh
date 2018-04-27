package isledef

import "time"

//GameMode represents the difficulty/mode of the game
type GameMode int

//Currently available game modes
const (
	GameModeNormalTimed GameMode = iota
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

//NewGameData creates a new GameData to keep track of the game
func NewGameData(gameMode GameMode) *GameData {
	g := &GameData{
		currentScale: 1,
		startTime:    time.Now(),
		server:       "127.0.0.1",
		port:         "5555",
	}
	switch gameMode {
	case GameModeNormalTimed:
		g.shotMax = 24
	case GameModeHardcoreTimed:
		g.shotMax = 15

	}
	return g
}

//SetGameMode sets the game mode chosen by the character. In this case used for Timed mode to set number of bombs
func (g *GameData) SetGameMode(gameMode GameMode) {
	switch gameMode {
	case GameModeNormalTimed:
		g.shotMax = 24
	case GameModeHardcoreTimed:
		g.shotMax = 15
	}
	g.gameMode = gameMode
}

//Update the game time
func (g *GameData) Update() {
	g.time++
}

//TimeInSecond returns the current time in seconds
func (g *GameData) TimeInSecond() int {

	return int(time.Now().Sub(g.startTime).Seconds())

	//return g.time / 60
}

//TimeInMilliseconds returns the current time in seconds
func (g *GameData) TimeInMilliseconds() int64 {

	return time.Now().Sub(g.startTime).Nanoseconds() / int64(time.Millisecond)

	//return g.time / 60
}
