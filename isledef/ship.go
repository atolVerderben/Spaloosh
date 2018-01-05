package isledef

import (
	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//Ship represents the ships on the board (in other words the game pieces)
type Ship struct {
	imgParts *tentsuyu.BasicImageParts
	*tentsuyu.BasicObject
	hitImg      *tentsuyu.BasicImageParts
	shipType    ShipType
	vertical    bool
	sections    [][]int
	dead        bool
	visible     bool
	sectionSize float64
	beingMoved  bool
}

//ShipType represents the different ship types
type ShipType string

//All the available ship types
const (
	ShipTypeBattleShip ShipType = "Bake-Kujira"
	ShipTypeCruiser             = "Kraken"
	ShipTypePatrol              = "Umibozu"
	ShipTypeSubmarine           = "Mishipeshu"
)

//CreateShip at point (x,y) or type shipType
func CreateShip(x, y float64, shipType ShipType, vertical bool) *Ship {
	sectionSize := 32.0
	s := &Ship{
		BasicObject: &tentsuyu.BasicObject{
			X:           x,
			Y:           y,
			NotCentered: true,
		},
		sectionSize: sectionSize,
		hitImg: &tentsuyu.BasicImageParts{
			Sx:     64,
			Sy:     0,
			Width:  int(32),
			Height: int(32),
		},
		shipType: shipType,
		vertical: vertical,
		//visible:  true,
	}
	switch shipType {
	case ShipTypeCruiser:
		s.sections = make([][]int, 3)
		s.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameKraken].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameKraken].Frame["y"],
			Width:  SpalooshSheet.Frames[frameKraken].Frame["w"],
			Height: SpalooshSheet.Frames[frameKraken].Frame["h"],
		}
	case ShipTypeSubmarine:
		s.sections = make([][]int, 3)
		s.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameMishipeshu].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameMishipeshu].Frame["y"],
			Width:  SpalooshSheet.Frames[frameMishipeshu].Frame["w"],
			Height: SpalooshSheet.Frames[frameMishipeshu].Frame["h"],
		}
	case ShipTypePatrol:
		s.sections = make([][]int, 2)
		s.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameUmibozu].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameUmibozu].Frame["y"],
			Width:  SpalooshSheet.Frames[frameUmibozu].Frame["w"],
			Height: SpalooshSheet.Frames[frameUmibozu].Frame["h"],
		}
	case ShipTypeBattleShip:
		s.sections = make([][]int, 4)
		s.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameBake].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameBake].Frame["y"],
			Width:  SpalooshSheet.Frames[frameBake].Frame["w"],
			Height: SpalooshSheet.Frames[frameBake].Frame["h"],
		}
	default:
		s.sections = make([][]int, 0)
	}
	/*for i := 0; i < len(s.sections); i++ {
		s.sections[i] = []int{0, 0, 0}
	}*/
	s.Width = int(float64(len(s.sections)) * sectionSize * ZoomLevel)
	s.Height = int(sectionSize * ZoomLevel)
	if vertical {
		s.SetAngle(1.5708)
		s.Height = s.Width
		s.Width = int(sectionSize * ZoomLevel)
	}
	//fmt.Printf("%s Width: %v Height: %v\n", s.shipType, s.Width, s.Height)
	return s
}

//Update the ship while it is being moved about the board
func (s *Ship) Update() {
	//if s.beingMoved {
	s.X, s.Y = tentsuyu.Input.GetMouseCoords()
	if tentsuyu.Input.Button("RotateLeft").JustPressed() {
		s.AddAngle(-rightAngle)
		if s.vertical {
			s.vertical = false
			s.Width = int(float64(len(s.sections)) * s.sectionSize * ZoomLevel)
			s.Height = int(s.sectionSize * ZoomLevel)
		} else {
			s.vertical = true
			s.Height = s.Width
			s.Width = int(s.sectionSize * ZoomLevel)
		}
	}
	if tentsuyu.Input.Button("RotateRight").JustPressed() {
		s.AddAngle(rightAngle)
		if s.vertical {
			s.vertical = false
			s.Width = int(float64(len(s.sections)) * s.sectionSize * ZoomLevel)
			s.Height = int(s.sectionSize * ZoomLevel)
		} else {
			s.vertical = true
			s.Height = s.Width
			s.Width = int(s.sectionSize * ZoomLevel)
		}
	}
	//}
}

//Draw the ship
func (s *Ship) Draw(screen *ebiten.Image) error {
	if s.visible == false {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = s.imgParts
	op.GeoM.Scale(ZoomLevel, ZoomLevel)
	if s.vertical {
		op.GeoM.Translate(-float64(s.Height/2), -float64(s.Width/2))
		op.GeoM.Rotate(s.Angle)
		op.GeoM.Translate(float64(s.Width/2), float64(s.Height/2)) //Switch the width and height because it's now vertical
	}
	op.GeoM.Translate(s.GetPosition())
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	return nil
}

//Hit is called when the ship is hit with a correct shot
func (s *Ship) Hit(x, y int) {
	hitCount := 0
	for i := 0; i < len(s.sections); i++ {
		if s.sections[i][2] >= 1 {
			hitCount++
		}
		if s.sections[i][0] == x && s.sections[i][1] == y && s.sections[i][2] == 0 {
			hitCount++
			//fmt.Printf("%s: I'm HIT! \n", s.shipType)
			s.sections[i][2] = 1
		}
	}
	if hitCount == len(s.sections) {
		//fmt.Printf("%s: I'm dead\n", s.shipType)
		s.dead = true
		s.visible = true
	}
}

//ReturnHitCount returns the number of hits the ship can take before sinking
func (s *Ship) ReturnHitCount() int {
	return len(s.sections)
}
