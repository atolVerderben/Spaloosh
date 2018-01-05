package network

type Command struct {
	Row, Col       int
	Name           string
	Message        string
	Rooms          []int
	CommType       CommandType
	ShipPlacements []*PlaceShipType
}

type CommandType int

const (
	CommandConnected CommandType = iota
	CommandAttack
	CommandDisconnected
	CommandSetTheBoard
	CommandHello
	CommandFull
	CommandHealthCheck
	CommandJoinRoom
	CommandListRooms
	CommandLeaveRoom
	CommandRoomFull
	CommandYouStart
)

type PlaceShipType struct {
	Row, Col   int
	Vertical   bool
	ShipLength int
	ShipType   string
}
