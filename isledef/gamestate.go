package isledef

import "github.com/atolVerderben/tentsuyu"

//These are the available GameStateMsg
const (
	GameStateMsgNone                     tentsuyu.GameStateMsg = "None"
	GameStateMsgReqTitle                                       = "Request Title"
	GameStateMsgReqMain                                        = "Request Main"
	GameStateGameOver                                          = "Game Over Lose"
	GameStateGameWin                                           = "Game Over Win"
	GameStateMsgReqMainMenu                                    = "Request MainMenu"
	GameStateMsgReqBattle                                      = "Request battle"
	GameStateMsgPause                                          = "Request Pause"
	GameStateMsgUnPause                                        = "UnPause"
	GameStateMsgReqMPStage                                     = "Request MP Stage"
	GameStateMsgReqMPMain                                      = "Request MP Main"
	GameStateMsgReqMPMainMenu                                  = "Request MP MainMenu"
	GameStateMsgReqLostConnection                              = "Request Lost Connection"
	GameStateMsgReqBattleCharacterSelect                       = "Battle Character Select"
	GameStateMsgReqMPHelp                                      = "MP Help"
	GameStateMsgReqMPGameOverWin                               = "Game Over Win MP"
	GameStateMsgReqMPGameOverLose                              = "Game Over Lose MP"
	GameStateMsgReqHostingRooms                                = "MP Hosting Rooms"
	GameStateMsgReqSetIP                                       = "Set IP"
)
