package network

import (
	"encoding/gob"
	"log"
	"net"
)

type Player struct {
	ID int
	// The connection to the player
	conn        *net.TCPConn
	wsEncoder   *gob.Encoder
	charDecoder *gob.Decoder
}

//CreatePlayer adds the player to the slice of players
func (s *GameServer) CreatePlayer(conn *net.TCPConn) bool {
	if s.remotePlayer != nil { // we only want to allow one other player
		return false
	}
	if len(s.Players) >= 2 { // only two players per game
		gob.NewEncoder(conn).Encode(&Command{CommType: CommandFull})
		conn.Close()
		return false
	}
	p := &Player{
		conn: conn,
		ID:   CreateNewID(),
	}
	p.wsEncoder = gob.NewEncoder(conn)
	p.charDecoder = gob.NewDecoder(conn)
	//s.remotePlayer = p
	s.Players = append(s.Players, p)
	comm := &Command{
		CommType: CommandConnected,
		Row:      len(s.Players),
	}
	log.Printf("Joined Game %v\n", len(s.Players))
	for _, player := range s.Players {

		player.wsEncoder.Encode(comm)

	}
	go HandleMessages(p, s)
	//Handle input from this player
	/*go func() {
		for {
			// Make a buffer to hold incoming data.
			buf := make([]byte, 1024)
			// Read the incoming connection into the buffer.
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}

			fmt.Println(string(buf))
		}
	}()*/
	return true
}

func HandleMessages(p *Player, s *GameServer) {
	select {
	case <-s.quit:
		return
	default:
	}
	for p.conn != nil {
		comm := &Command{}
		p.charDecoder.Decode(comm)

		for _, player := range s.Players {
			if player != p {
				player.wsEncoder.Encode(comm)
			}
		}
	}

	if p.conn == nil {
		comm := &Command{
			CommType: CommandDisconnected,
		}
		for _, player := range s.Players {
			if player != p {
				player.wsEncoder.Encode(comm)
			}
		}

		delete := -1
		for index, e := range s.Players {
			if e == p {
				delete = index
				break
			}
		}
		if delete >= 0 {
			s.Players = append(s.Players[:delete], s.Players[delete+1:]...)
		}

		log.Printf("Server Message: Player %v left game. %v Players left\n", p.ID, len(s.Players))
	}
}
