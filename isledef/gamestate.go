package isledef

//GameState represents the current state of the game
type GameState interface {
	Update(g *Game) error
	Draw(g *Game) error
	Msg() GameStateMsg
	SetMsg(GameStateMsg)
}

//GameStateMsg represents the messages sent from the GameState that may change states
type GameStateMsg int

//These are the available GameStateMsg
const (
	GameStateMsgNone GameStateMsg = iota
	GameStateMsgReqTitle
	GameStateMsgReqMain
	GameStateGameOver
	GameStateGameWin
	GameStateMsgReqMainMenu
	GameStateMsgReqBattle
	GameStateMsgPause
	GameStateMsgUnPause
	GameStateMsgReqMPStage
	GameStateMsgReqMPMain
	GameStateMsgReqMPMainMenu
	GameStateMsgReqLostConnection
	GameStateMsgReqBattleCharacterSelect
	GameStateMsgReqMPHelp
	GameStateMsgReqMPGameOverWin
	GameStateMsgReqMPGameOverLose
	GameStateMsgReqHostingRooms
	GameStateMsgReqSetIP
)
