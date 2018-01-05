package network

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

//IDMAP is a map of current in use IDs... this is overkill for this game, but it's here
var IDMAP map[int]bool

//GameServer is our main netorking entity
type GameServer struct {
	remotePlayer  *Player
	Players       []*Player
	ipConn        *net.TCPListener
	port          string
	tick          int
	startTickTime time.Time
	quit          chan bool
}

//CreateGameServer returns a pointer to a GameServer struct
func CreateGameServer(port string) *GameServer {
	s := &GameServer{
		port: port,
	}
	return s
}

//Run the GameServer
func (s *GameServer) Run() {

	server, err := net.Listen("tcp", ":"+s.port)
	if server == nil {
		panic("listen failed: " + err.Error() + "\n")
	} else {
		//defer server.Close()
	}
	s.ipConn = server.(*net.TCPListener)
	IDMAP = make(map[int]bool)
	s.quit = make(chan bool)
	// connection handling
	go func(s *GameServer) {
		for {
			select {
			case <-s.quit:
				log.Println("Server shutdown")
				return
			default:
				conn, err := s.ipConn.Accept()
				if err != nil {

					fmt.Printf("client error: %s\n", err.Error())
					return
				} else {
					s.CreatePlayer(conn.(*net.TCPConn))
				}
			}

		}
	}(s)

	go func(s *GameServer) {
		timer := time.NewTicker(50 * time.Millisecond)
		for now := range timer.C {
			// entity updates
			// this is called every 100 millisecondes
			select {
			case <-s.quit:
				return
			default:
				s.CheckPlayers(now)
			}
		}
	}(s)

}

func (s *GameServer) ShutDown() {
	s.quit <- true
	for _, p := range s.Players {
		p.conn.Close()
	}
	s.ipConn.Close()
	s.ipConn = nil
}

func Read(c *net.TCPConn, buffer []byte) bool {
	bytes, err := c.Read(buffer)
	if err != nil {
		c.Close()
		log.Println(err)
		return false
	}
	//byt, _ := ioutil.ReadAll(c)
	log.Println("Read ", bytes, " bytes")
	return true
}

func Write(p *Player) bool {
	comm := &Command{
		CommType: CommandHealthCheck,
	}

	err := p.wsEncoder.Encode(comm)
	if err != nil {
		return false
	}

	return true

}

func (s *GameServer) CheckPlayers(t time.Time) {
	if s.ipConn == nil {
		return
	}
	//log.Printf("Server: I'm looking for players, have %v\n", len(s.Players))
	removeList := []int{}
	for _, p := range s.Players {
		if !Write(p) {
			comm := &Command{
				CommType: CommandDisconnected,
			}
			for _, player := range s.Players {
				if player != p {
					player.wsEncoder.Encode(comm)
				}
			}
			p.conn.Close()
			p.conn = nil
			removeList = append(removeList, p.ID)

		}
	}

	for _, id := range removeList {
		delete := -1
		for index, e := range s.Players {
			if e.ID == id {
				delete = index
				break
			}
		}
		if delete >= 0 {
			s.Players = append(s.Players[:delete], s.Players[delete+1:]...)
			log.Printf("Server Message: Player %v left game. %v Players left\n", id, len(s.Players))
		}
	}
}

//CreateNewID returns a unique ID for a new player
func CreateNewID() int {
	max := 99999999999
	var i = rand.Intn(max)
	for IDMAP[i] == false {
		IDMAP[i] = true
		i = rand.Intn(max)
	}
	return i
}
