package isledef

import (
	"fmt"
	"log"
	"math"

	"github.com/atolVerderben/spaloosh/network"
	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//Represent the different tiles on the grid
const (
	TileEmpty int = iota
	TileOccupied
	TileEmptyShot
	TileOccupiedShot
)

//Grid represents the game playing field
type Grid struct {
	*tentsuyu.BasicObject
	imgParts               *tentsuyu.BasicImageParts
	Rows                   int
	Columns                int
	area                   [][]int
	animateArea            [][]int
	tileSize               int
	tileSizeF              float64
	Ships                  []*Ship
	missImg                *tentsuyu.BasicImageParts
	hitImg                 *tentsuyu.BasicImageParts
	Cleared                bool
	shipsDefeated          int
	shipsRemaining         int
	allVisible             bool
	playable               bool
	animCount1, animCount2 int
	placeable              bool
	prevHit                bool
}

//CreateGrid returns a Grid which represents the game field
func CreateGrid(x, y float64, rows, cols int) *Grid {
	g := &Grid{
		Rows:       rows,
		Columns:    cols,
		animCount1: 5,
		animCount2: 10,
		BasicObject: &tentsuyu.BasicObject{
			X:           x,
			Y:           y,
			Width:       int(32*(ZoomLevel)) * cols,
			Height:      int(32*(ZoomLevel)) * rows,
			NotCentered: true,
		},
		hitImg: &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameSpaloosh1].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameSpaloosh1].Frame["y"],
			Width:  32,
			Height: 32,
		},
		missImg: &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameKaboom1].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameKaboom1].Frame["y"],
			Width:  32,
			Height: 32,
		},
		imgParts: &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameGridSmall].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameGridSmall].Frame["y"],
			Width:  257,
			Height: 257,
			//	DestHeight: 128,
			//	DestWidth:  128,
		},
		area:           [][]int{},
		animateArea:    [][]int{},
		tileSize:       int(32 * ZoomLevel),
		tileSizeF:      32 * ZoomLevel,
		shipsRemaining: 3,
		allVisible:     false,
		placeable:      false,
		playable:       true,
	}
	if g.Rows > 8 {
		g.imgParts.Sx = SpalooshSheet.Frames[frameGridBig].Frame["x"]
		g.imgParts.Sy = SpalooshSheet.Frames[frameGridBig].Frame["y"]
		g.imgParts.Width = SpalooshSheet.Frames[frameGridBig].Frame["w"]
		g.imgParts.Height = SpalooshSheet.Frames[frameGridBig].Frame["h"]
	}
	for i := 0; i < g.Rows; i++ {
		g.area = append(g.area, make([]int, cols))
		g.animateArea = append(g.animateArea, make([]int, cols))
	}

	g.imgParts.Width = g.Columns*g.tileSize + 1
	g.imgParts.Height = g.Rows*g.tileSize + 1

	//fmt.Printf("Width: %v Height: %v\n", g.Width, g.Height)
	return g
}

func (g *Grid) drawShips(screen *ebiten.Image) error {
	for _, s := range g.Ships {
		s.Draw(screen)
		for _, section := range s.sections {
			if section[2] >= 1 {
				if section[2] >= 1 && section[2] <= g.animCount2 {
					section[2]++
				}
				op := &ebiten.DrawImageOptions{}
				if section[2] < g.animCount1 {
					op.ImageParts = &tentsuyu.BasicImageParts{
						Sx:     SpalooshSheet.Frames[frameKaboom1].Frame["x"],
						Sy:     SpalooshSheet.Frames[frameKaboom1].Frame["y"],
						Width:  32,
						Height: 32,
					}
				} else if section[2] >= g.animCount1 && section[2] < g.animCount2 {
					op.ImageParts = &tentsuyu.BasicImageParts{
						Sx:     SpalooshSheet.Frames[frameKaboom2].Frame["x"],
						Sy:     SpalooshSheet.Frames[frameKaboom2].Frame["y"],
						Width:  32,
						Height: 32,
					}
				} else {
					op.ImageParts = &tentsuyu.BasicImageParts{
						Sx:     SpalooshSheet.Frames[frameKaboom3].Frame["x"],
						Sy:     SpalooshSheet.Frames[frameKaboom3].Frame["y"],
						Width:  32,
						Height: 32,
					}
				}
				op.GeoM.Scale(ZoomLevel, ZoomLevel)
				op.GeoM.Translate(g.ReturnXYCoords(section[0], section[1]))
				if s.visible {
					op.ColorM.Scale(1, 1, 1, 0.5)
				}
				if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//Draw the Grid
func (g *Grid) Draw(screen *ebiten.Image) error {

	op := &ebiten.DrawImageOptions{}

	op.ImageParts = g.imgParts
	//op.GeoM.Scale(ZoomLevel, ZoomLevel)
	op.GeoM.Translate(g.X, g.Y)
	//op.ColorM.Scale(0, 0, 0, 1)
	//tentsuyu.ApplyCameraTransform(op, true)

	if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	g.drawShips(screen)
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Columns; j++ {
			switch g.area[i][j] {
			case TileEmptyShot:
				if g.animateArea[i][j] >= 1 && g.animateArea[i][j] <= g.animCount2 {
					g.animateArea[i][j]++
				}
				op := &ebiten.DrawImageOptions{}
				if g.animateArea[i][j] < g.animCount1 {
					op.ImageParts = &tentsuyu.BasicImageParts{
						Sx:     SpalooshSheet.Frames[frameSpaloosh1].Frame["x"],
						Sy:     SpalooshSheet.Frames[frameSpaloosh1].Frame["y"],
						Width:  32,
						Height: 32,
					}
				} else if g.animateArea[i][j] >= g.animCount1 && g.animateArea[i][j] < g.animCount2 {
					op.ImageParts = &tentsuyu.BasicImageParts{
						Sx:     SpalooshSheet.Frames[frameSpaloosh2].Frame["x"],
						Sy:     SpalooshSheet.Frames[frameSpaloosh2].Frame["y"],
						Width:  32,
						Height: 32,
					}
				} else {
					op.ImageParts = &tentsuyu.BasicImageParts{
						Sx:     SpalooshSheet.Frames[frameSpaloosh3].Frame["x"],
						Sy:     SpalooshSheet.Frames[frameSpaloosh3].Frame["y"],
						Width:  32,
						Height: 32,
					}
				}

				op.GeoM.Scale(ZoomLevel, ZoomLevel)
				op.GeoM.Translate(g.ReturnXYCoords(i, j))

				if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
					return err
				}
				//I might come back to this at a later date
			/*case TileOccupiedShot:
			op := &ebiten.DrawImageOptions{}
			op.ImageParts = g.hitImg
			op.GeoM.Scale(ZoomLevel, ZoomLevel)
			op.GeoM.Translate(g.ReturnXYCoords(i, j))
			op.ColorM.Scale(1, 1, 1, 1)
			if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("explosion"), op); err != nil {
				return err
			}*/
			default:
			}
		}
	}

	return nil
}

//Update the grid
func (g *Grid) Update(game *tentsuyu.Game) bool {
	shotHit := false
	if game.Input.LeftClick().JustReleased() && g.playable {
		x, y := game.Input.GetMouseCoords()
		//fmt.Printf("Mouse X: %v Mouse Y: %v\n", x-g.X, y-g.Y)
		if g.Contains(x, y) {
			tile, row, col := g.DetermineTile(x-g.X, y-g.Y)
			switch tile {
			case TileEmpty:
				if GamePlayer.TakeShot() {
					shotHit = true
					switch gs := game.GetGameState().(type) {
					case *GameMain:
						gs.bulletCounter.Fire()
					}

					g.area[row][col] = TileEmptyShot
					g.animateArea[row][col] = 1
					g.prevHit = false
					if err := Game.AudioPlayer.PlaySE(seMiss); err != nil {
						log.Printf("Error: %v\n", err)
						//panic(err)
					}
				}

			case TileOccupied:
				for i := 0; i < len(g.Ships); i++ {
					if g.Ships[i].Contains(x, y) && g.Ships[i].dead == false {
						if GamePlayer.TakeShot() {
							shotHit = true
							switch gs := game.GetGameState().(type) {
							case *GameMain:
								gs.bulletCounter.Fire()
							}
							g.Ships[i].Hit(row, col)
							g.prevHit = true
							if err := Game.AudioPlayer.PlaySE(seHit); err != nil {
								log.Printf("Error: %v\n", err)
								//panic(err)
							}
							if g.Ships[i].dead {
								g.shipsDefeated++
								g.shipsRemaining--
							}
							g.area[row][col] = TileOccupiedShot
							g.animateArea[row][col] = 1
						}

					}
				}

			default:
			}
			//fmt.Printf("(%v)\n", g.DetermineTile(x-g.X, y-g.Y))

		}
	}
	if g.Cleared == false && g.shipsDefeated >= (len(g.Ships)) {
		g.Cleared = true
	}
	return shotHit
}

//SpecialAttack is unique to each game character
func (g *Grid) SpecialAttack(game *tentsuyu.Game) bool {
	shotHit := false
	if game.Input.LeftClick().JustReleased() && g.playable {
		x, y := game.Input.GetMouseCoords()
		//fmt.Printf("Mouse X: %v Mouse Y: %v\n", x-g.X, y-g.Y)
		if g.Contains(x, y) {
			tile, row, col := g.DetermineTile(x-g.X, y-g.Y)
			switch tile {
			case TileEmpty:
				if GamePlayer.TakeShot() {
					shotHit = true
					switch gs := game.GetGameState().(type) {
					case *GameMain:
						gs.bulletCounter.Fire()
					}

					g.area[row][col] = TileEmptyShot
					g.animateArea[row][col] = 1

					if err := Game.AudioPlayer.PlaySE(seMiss); err != nil {
						log.Printf("Error: %v\n", err)
						//panic(err)
					}
				}

			case TileOccupied:
				for i := 0; i < len(g.Ships); i++ {
					if g.Ships[i].Contains(x, y) && g.Ships[i].dead == false {
						if GamePlayer.TakeShot() {
							shotHit = true
							switch gs := game.GetGameState().(type) {
							case *GameMain:
								gs.bulletCounter.Fire()
							}
							g.Ships[i].Hit(row, col)
							if err := Game.AudioPlayer.PlaySE(seHit); err != nil {
								log.Printf("Error: %v\n", err)
								//panic(err)
							}
							if g.Ships[i].dead {
								g.shipsDefeated++
								g.shipsRemaining--
							}
							g.area[row][col] = TileOccupiedShot
							g.animateArea[row][col] = 1
						}

					}
				}

			default:
			}
			//fmt.Printf("(%v)\n", g.DetermineTile(x-g.X, y-g.Y))

		}
	}
	if g.Cleared == false && g.shipsDefeated >= (len(g.Ships)) {
		g.Cleared = true
	}
	return shotHit
}

//Update the grid
func (g *Grid) MPUpdate(game *tentsuyu.Game) (bool, int, int) {
	shotHit := false
	returnRow, returnCol := 0, 0
	if game.Input.LeftClick().JustReleased() && g.playable {
		x, y := game.Input.GetMouseCoords()
		//fmt.Printf("Mouse X: %v Mouse Y: %v\n", x-g.X, y-g.Y)
		if g.Contains(x, y) {
			tile, row, col := g.DetermineTile(x-g.X, y-g.Y)
			switch tile {
			case TileEmpty:
				if GamePlayer.TakeShot() {
					shotHit = true
					switch gs := game.GetGameState().(type) {
					case *GameMain:
						gs.bulletCounter.Fire()
					}

					g.area[row][col] = TileEmptyShot
					g.prevHit = false
					g.animateArea[row][col] = 1
					returnRow = row
					returnCol = col

					if err := Game.AudioPlayer.PlaySE(seMiss); err != nil {
						log.Printf("Error: %v\n", err)
						//panic(err)
					}
				}

			case TileOccupied:
				for i := 0; i < len(g.Ships); i++ {
					if g.Ships[i].Contains(x, y) && g.Ships[i].dead == false {
						if GamePlayer.TakeShot() {
							shotHit = true
							switch gs := game.GetGameState().(type) {
							case *GameMain:
								gs.bulletCounter.Fire()
							}
							g.Ships[i].Hit(row, col)
							if err := Game.AudioPlayer.PlaySE(seHit); err != nil {
								log.Printf("Error: %v\n", err)
								//panic(err)
							}
							if g.Ships[i].dead {
								g.shipsDefeated++
								g.shipsRemaining--
							}
							g.area[row][col] = TileOccupiedShot
							g.prevHit = true
							g.animateArea[row][col] = 1
							returnRow = row
							returnCol = col
						}

					}
				}

			default:
			}
			//fmt.Printf("(%v)\n", g.DetermineTile(x-g.X, y-g.Y))

		}
	}
	if g.Cleared == false && g.shipsDefeated >= (len(g.Ships)) {
		g.Cleared = true
	}
	return shotHit, returnRow, returnCol
}

//Update the grid
func (g *Grid) MPSpecial(game *tentsuyu.Game) (bool, int, int) {
	shotHit := false
	returnRow, returnCol := 0, 0
	if game.Input.LeftClick().JustReleased() && g.playable {
		x, y := game.Input.GetMouseCoords()
		//fmt.Printf("Mouse X: %v Mouse Y: %v\n", x-g.X, y-g.Y)
		if g.Contains(x, y) {
			tile, row, col := g.DetermineTile(x-g.X, y-g.Y)
			switch tile {
			case TileEmpty:
				if GamePlayer.TakeShot() {
					shotHit = true
					switch gs := game.GetGameState().(type) {
					case *GameMain:
						gs.bulletCounter.Fire()
					}

					g.area[row][col] = TileEmptyShot
					g.animateArea[row][col] = 1
					returnRow = row
					returnCol = col

					if err := Game.AudioPlayer.PlaySE(seMiss); err != nil {
						log.Printf("Error: %v\n", err)
						//panic(err)
					}
				}

			case TileOccupied:
				for i := 0; i < len(g.Ships); i++ {
					if g.Ships[i].Contains(x, y) && g.Ships[i].dead == false {
						if GamePlayer.TakeShot() {
							shotHit = true
							switch gs := game.GetGameState().(type) {
							case *GameMain:
								gs.bulletCounter.Fire()
							}
							g.Ships[i].Hit(row, col)
							if err := Game.AudioPlayer.PlaySE(seHit); err != nil {
								log.Printf("Error: %v\n", err)
								//panic(err)
							}
							if g.Ships[i].dead {
								g.shipsDefeated++
								g.shipsRemaining--
							}
							g.area[row][col] = TileOccupiedShot
							g.animateArea[row][col] = 1
							returnRow = row
							returnCol = col
						}

					}
				}

			default:
			}
			//fmt.Printf("(%v)\n", g.DetermineTile(x-g.X, y-g.Y))

		}
	}
	if g.Cleared == false && g.shipsDefeated >= (len(g.Ships)) {
		g.Cleared = true
	}
	return shotHit, returnRow, returnCol
}

//AIShot takes a row and col and returns 3 values: ShotMiss, ShotHit, ShipSunk
func (g *Grid) AIShot(row, col int) (bool, bool, bool) {
	tile := g.area[row][col]
	shotHit := false
	shotEmpty := false
	shipSunk := false
	//tile, col, row := g.DetermineTile(x-g.X, y-g.Y)
	switch tile {
	case TileEmpty:

		shotEmpty = true
		g.prevHit = false

		g.area[row][col] = TileEmptyShot
		g.animateArea[row][col] = 1

		if err := Game.AudioPlayer.PlaySE(seMiss); err != nil {
			log.Printf("Error: %v\n", err)
			//panic(err)
		}

	case TileOccupied:
		for i := 0; i < len(g.Ships); i++ {
			x, y := g.ReturnXYCoords(row, col)
			if g.Ships[i].Contains(x+2, y+2) && g.Ships[i].dead == false {

				shotHit = true
				g.prevHit = true

				g.Ships[i].Hit(row, col)
				if err := Game.AudioPlayer.PlaySE(seHit); err != nil {
					log.Printf("Error: %v\n", err)
					//panic(err)
				}
				if g.Ships[i].dead {
					g.shipsDefeated++
					g.shipsRemaining--
					shipSunk = true
				}
				g.area[row][col] = TileOccupiedShot
				g.animateArea[row][col] = 1
			}
		}

	default:
	}
	return shotEmpty, shotHit, shipSunk
}

//MakeAllVisible is used at the end of a round to see any ships you may have missed
func (g *Grid) MakeAllVisible() {
	for _, s := range g.Ships {
		s.visible = true
	}
	g.allVisible = true
}

//DetermineTile returns the value at the grid
//This is used to use the mouse coords
func (g Grid) DetermineTile(x, y float64) (int, int, int) {
	tx, ty := int(math.Floor(x/(float64(g.tileSize)))), int(math.Floor(y/(float64(g.tileSize))))
	//fmt.Printf("Rows: %v Cols: %v\n", tx, ty)
	return g.area[ty][tx], ty, tx
}

//ReturnXYCoords returns the float values for the given row and col
func (g Grid) ReturnXYCoords(row, col int) (float64, float64) {
	return float64(g.tileSize*col) + g.X, float64(g.tileSize*row) + g.Y
}

//PlaceShip puts the ship at the right place in the grid
func (g *Grid) PlaceShip(row, col int, shipType ShipType, vertical bool) {
	sx, sy := g.ReturnXYCoords(row, col)
	ship := CreateShip(sx, sy, shipType, vertical)

	if vertical {
		for i := 0; i < len(ship.sections); i++ {
			g.area[row+i][col] = TileOccupied
			ship.sections[i] = []int{row + i, col, 0}
		}
	} else {
		for i := 0; i < len(ship.sections); i++ {
			g.area[row][col+i] = TileOccupied
			ship.sections[i] = []int{row, col + i, 0}
		}
	}
	g.Ships = append(g.Ships, ship)
	g.shipsRemaining = len(g.Ships)

}

func (g *Grid) UpdatePlacement(game *tentsuyu.Game) {

	if game.Input.LeftClick().JustReleased() && !GamePlayer.holdingShip {
		x, y := game.Input.GetMouseCoords()
		//fmt.Printf("Mouse X: %v Mouse Y: %v\n", x-g.X, y-g.Y)
		if g.Contains(x, y) {
			tile, _, _ := g.DetermineTile(x-g.X, y-g.Y)
			switch tile {
			case TileOccupied:
				for i := 0; i < len(g.Ships); i++ {
					if g.Ships[i].Contains(x, y) {
						GamePlayer.heldShip = g.Ships[i]
						GamePlayer.holdingShip = true
						g.RemoveShip(GamePlayer.heldShip)
					}
				}

			default:
			}
		}
	} else if game.Input.LeftClick().JustReleased() && GamePlayer.holdingShip {
		x, y := game.Input.GetMouseCoords()
		//fmt.Printf("Mouse X: %v Mouse Y: %v\n", x-g.X, y-g.Y)
		if g.Contains(x, y) {
			tile, row, col := g.DetermineTile(x-g.X, y-g.Y)
			switch tile {
			case TileEmpty:
				if PlayerChoosePosition(col, row, len(GamePlayer.heldShip.sections), GamePlayer.heldShip.vertical, g) {
					g.PlaceShip(row, col, GamePlayer.heldShip.shipType, GamePlayer.heldShip.vertical)
					GamePlayer.heldShip = nil
					GamePlayer.holdingShip = false
					g.MakeAllVisible()
				}

			default:
			}
		}
	}

}

func (g *Grid) RemoveShip(ship *Ship) {

	for _, sec := range ship.sections {
		g.area[sec[0]][sec[1]] = TileEmpty
	}

	delete := -1
	for index, s := range g.Ships {
		if s == ship {
			delete = index
			break
		}
	}
	if delete >= 0 {
		g.Ships = append(g.Ships[:delete], g.Ships[delete+1:]...)
	}

}

func (g *Grid) ExportNetworkPlacements() []*network.PlaceShipType {
	shipPlacements := []*network.PlaceShipType{}

	for _, s := range g.Ships {
		_, row, col := g.DetermineTile(s.X-g.X, s.Y-g.Y)
		shipPlacements = append(shipPlacements, &network.PlaceShipType{
			Row:        row,
			Col:        col,
			Vertical:   s.vertical,
			ShipLength: len(s.sections),
			ShipType:   string(s.shipType),
		})
	}

	return shipPlacements
}

//PrintGrid prints the text representation of the game grid
func (g Grid) PrintGrid() {
	for i := 0; i < g.Rows; i++ {
		fmt.Println(g.area[i])
	}
}
