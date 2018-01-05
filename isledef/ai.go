package isledef

import (
	"math/rand"

	"github.com/atolVerderben/spaloosh/network"
)

//AI represents the computer controlled opponent
type AI struct {
	Name                                      string
	Level                                     int
	ships                                     []Ship
	prevHit                                   bool
	prevSunk                                  bool
	prevRow, prevCol                          int
	triedLeft, triedRight, triedUp, triedDown bool
	shipIsHoriz, shipIsVert                   bool
	searchIncrememnt                          int
	hitTally                                  int
	umiSunk, krakenSunk, bakeSunk, mishiSunk  bool
	hitLocations                              [][]int
	currHitTally                              int
}

func NewAI() *AI {
	ai := &AI{
		hitLocations: [][]int{},
	}
	return ai
}

//SetBoard is called at the beginning of the game so the AI sets its pieces on the game board
func (ai *AI) SetBoard(g *Grid) {
	set := false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 2, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypePatrol, vertical)
		}
	}
	set = false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 3, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypeCruiser, vertical)
		}
	}
	set = false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 4, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypeBattleShip, vertical)
		}
	}
	if g.Rows == 8 {
		return
	}
	set = false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 3, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypeSubmarine, vertical)
		}
	}
	//g.PrintGrid()
	//g.PlaceShip(1, 1, ShipTypeCruiser, false)
	//g.PlaceShip(5, 4, ShipTypeBattleShip, false)
	//g.PlaceShip(3, 7, ShipTypePatrol, true)
}

//SetBoard is called at the beginning of the game so the AI sets its pieces on the game board
func (ai *AI) NetworkSetMyBoard(g *Grid) []*network.PlaceShipType {

	shipPlacements := []*network.PlaceShipType{}
	set := false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 2, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypePatrol, vertical)
			shipPlacements = append(shipPlacements, &network.PlaceShipType{
				Row:        row,
				Col:        col,
				Vertical:   vertical,
				ShipLength: 2,
			})
		}
	}
	set = false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 3, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypeCruiser, vertical)
			shipPlacements = append(shipPlacements, &network.PlaceShipType{
				Row:        row,
				Col:        col,
				Vertical:   vertical,
				ShipLength: 3,
			})
		}
	}
	set = false
	for set == false {
		vertical := false
		row, col := rand.Intn(g.Rows), rand.Intn(g.Columns)
		set, vertical = ChoosePosition(col, row, 4, g)
		if set == true {
			g.PlaceShip(row, col, ShipTypeBattleShip, vertical)
			shipPlacements = append(shipPlacements, &network.PlaceShipType{
				Row:        row,
				Col:        col,
				Vertical:   vertical,
				ShipLength: 4,
			})
		}
	}
	//g.PrintGrid()
	//g.PlaceShip(1, 1, ShipTypeCruiser, false)
	//g.PlaceShip(5, 4, ShipTypeBattleShip, false)
	//g.PlaceShip(3, 7, ShipTypePatrol, true)
	return shipPlacements
}

//SetBoard is called at the beginning of the game so the AI sets its pieces on the game board
func (ai *AI) NetworkSetOpponentBoard(g *Grid, shipPlacements []*network.PlaceShipType) {

	for _, sp := range shipPlacements {
		g.PlaceShip(sp.Row, sp.Col, ShipType(sp.ShipType), sp.Vertical)
		/*switch sp.ShipType {
		case 2:
			g.PlaceShip(sp.Row, sp.Col, ShipTypePatrol, sp.Vertical)
		case 3:
			g.PlaceShip(sp.Row, sp.Col, ShipTypeCruiser, sp.Vertical)
		case 4:
			g.PlaceShip(sp.Row, sp.Col, ShipTypeBattleShip, sp.Vertical)
		}*/
	}

}

//ResetBoard removes the current ships and resets them for the given grid
func (ai *AI) ResetBoard(g *Grid) {
	g.Ships = nil
	g.area = [][]int{}
	for i := 0; i < g.Rows; i++ {
		g.area = append(g.area, make([]int, g.Columns))
	}
	g.animateArea = [][]int{}
	for i := 0; i < g.Rows; i++ {
		g.animateArea = append(g.animateArea, make([]int, g.Columns))
	}
	ai.SetBoard(g)
}

func checkHoriz(row, col, length int, g *Grid) bool {
	for i := 0; i < length; i++ {
		if col+i > g.Columns-1 {
			return true
		}
		if g.area[row][col+i] != TileEmpty {
			return true
		}
	}
	return false
}

func checkVert(row, col, length int, g *Grid) bool {
	for i := 0; i < length; i++ {
		if row+i > g.Rows-1 {
			return true
		}
		if g.area[row+i][col] != TileEmpty {
			return true
		}
	}
	return false
}

//ChoosePosition determines if the ship can fit in the spot given, including turning the ship vertically. Otherwise returns false
func ChoosePosition(col, row, length int, g *Grid) (bool, bool) {
	//fmt.Printf("I'm starting at (%v,%v)\n", col, row)
	if row > g.Rows-1 || col > g.Columns-1 {
		return false, false
	}
	if g.area[row][col] != TileEmpty {
		return false, false
	}
	r := rand.Intn(2)

	vertical := false
	if r == 1 {
		vertical = true
	}
	collide := false

	if !vertical {
		collide = checkHoriz(row, col, length, g)
	} else {
		collide = checkVert(row, col, length, g)
	}

	if collide == true { //Try to go vertical
		if !vertical {
			collide = checkHoriz(row, col, length, g)
		} else {
			collide = checkVert(row, col, length, g)
		}
		vertical = !vertical
	}
	if collide == true {
		return false, false
	}
	return true, vertical
}

//TakeShot has the AI take a shot at the battleship style grid
//Returns true if there was a successful shot taken (hit or miss)
func (ai *AI) TakeShot(grid *Grid) bool {
	if ai.prevSunk == true {
		//fmt.Printf("%v\n", ai.hitLocations)
		totalHit := 0
		for _, s := range grid.Ships {
			if s.dead == true {
				totalHit += s.ReturnHitCount()
				if s.shipType == ShipTypeBattleShip {
					ai.bakeSunk = true
				} else if s.shipType == ShipTypeCruiser {
					ai.krakenSunk = true
				} else if s.shipType == ShipTypePatrol {
					ai.umiSunk = true
				} else if s.shipType == ShipTypeSubmarine {
					ai.mishiSunk = true
				}
			}
		}
		if totalHit < ai.hitTally { //AI has hit one and not sunk it (happens when they are grouped together)
			//This should set the hit that is not part of a sunk ship to active
			for _, coords := range ai.hitLocations {
				x, y := grid.ReturnXYCoords(coords[0], coords[1])
				for i := range grid.Ships {
					if grid.Ships[i].Contains(x+2, y+2) && grid.Ships[i].dead == false {
						//Set these values so the AI skips to the targeted shooting phase
						ai.prevRow = coords[0]
						ai.prevCol = coords[1]
						ai.prevSunk = false
						ai.prevHit = true
						ai.currHitTally = 1
						ai.searchIncrememnt = 1
						ai.resetTryDirections()
						//log.Println("Should Get Here!!")
					}
				}
			}
		}
	}
	//If last shot wasn't a hit or if the last shot sunk a ship, aim randomly on the board
	if ai.prevHit == false || (ai.prevHit == true && ai.prevSunk == true) {
		ai.currHitTally = 0
		ai.resetTryDirections()
		row, col := rand.Intn(grid.Rows), rand.Intn(grid.Columns)
		if row == ai.prevRow && col == ai.prevCol {
			return false
		}
		if grid.area[row][col] == TileEmptyShot || grid.area[row][col] == TileOccupiedShot {
			return false // already shot here
		}

		//Check all around to see if this is viable spot
		leftBad, rightBad, upBad, downBad := false, false, false, false
		minLength := 2
		maxLength := 4
		if ai.bakeSunk {
			maxLength = 3
			if ai.krakenSunk && ai.mishiSunk {
				maxLength = 2
			}
		}
		if ai.umiSunk {
			minLength = 3
			if ai.krakenSunk && ai.mishiSunk {
				minLength = 4
			}
		}
		//log.Printf("Max: %v, Min: %v\n", maxLength, minLength)
		/*
			//Check Up
			if row != 0 {
				if grid.area[row-1][col] == TileEmptyShot || grid.area[row-1][col] == TileOccupiedShot {
					UpBad = true
				}
			} else {
				UpBad = true
			}
			//Check Down
			if row < grid.Rows-1 {
				if grid.area[row+1][col] == TileEmptyShot || grid.area[row+1][col] == TileOccupiedShot {
					DownBad = true
				}
			} else {
				DownBad = true
			}
			//Check Left
			if col != 0 {
				if grid.area[row][col-1] == TileEmptyShot || grid.area[row][col-1] == TileOccupiedShot {
					leftBad = true
				}
			} else {
				leftBad = true
			}
			//Check Right
			if col < grid.Columns-1 {
				if grid.area[row][col+1] == TileEmptyShot || grid.area[row][col+1] == TileOccupiedShot {
					rightBad = true
				}
			} else {
				rightBad = true
			}
		*/
		//Experimental Checks =======================================================================

		//Check Left
		if col != 0 {
			for i := 1; i < maxLength; i++ {
				if col-i >= 0 {
					if grid.area[row][col-i] == TileEmptyShot || grid.area[row][col-i] == TileOccupiedShot {
						if i < minLength-1 {
							leftBad = true
						}
					}
				} else if i < minLength-1 {
					leftBad = true
				}
			}
		} else {
			leftBad = true
		}

		//Check Right
		if col < grid.Columns {
			for i := 1; i < maxLength; i++ {
				if col+i < grid.Columns {
					if grid.area[row][col+i] == TileEmptyShot || grid.area[row][col+i] == TileOccupiedShot {
						if i < minLength-1 {
							rightBad = true
						}
					}
				} else if i < minLength-1 {
					rightBad = true
				}
			}

		} else {
			rightBad = true
		}

		//Check Up
		if row != 0 {
			for i := 1; i < maxLength; i++ {
				if row-i >= 0 {
					if grid.area[row-i][col] == TileEmptyShot || grid.area[row-i][col] == TileOccupiedShot {
						if i < minLength-1 {
							upBad = true
						}
					}
				} else if i < minLength-1 {
					upBad = true
				}
			}

		} else {
			upBad = true
		}

		//Check Down
		if row < grid.Rows-1 {
			col := ai.prevCol

			for i := 1; i < maxLength; i++ {
				if row+i < grid.Rows {
					if grid.area[row+i][col] == TileEmptyShot || grid.area[row+i][col] == TileOccupiedShot {
						if i < minLength-1 {
							downBad = true
						}
					}
				} else if i < minLength-1 {
					downBad = true
				}
			}

		} else {
			downBad = true
		}

		//===========================================================================================

		//square is surrounded by shots, no ship could fit here
		if leftBad && rightBad && upBad && downBad {
			return false
		}

		if grid.area[row][col] != TileEmptyShot && grid.area[row][col] != TileOccupiedShot {
			miss, hit, sunk := grid.AIShot(row, col)
			ai.prevHit = hit
			ai.prevSunk = sunk
			ai.prevCol = col
			ai.prevRow = row
			//log.Printf("Row: %v, Col:%v\n", row, col)
			if hit {
				ai.searchIncrememnt = 1
				ai.hitLocations = append(ai.hitLocations, []int{row, col})
				ai.hitTally++
				ai.currHitTally++
			}
			return miss || hit
		}
	} else { //Otherwise shoot around that shot until ship is sunk

		//Check all around to see if this is viable spot
		leftBad, rightBad, upBad, downBad := false, false, false, false
		minLength := 2
		maxLength := 4
		if ai.bakeSunk {
			maxLength = 3
			if ai.krakenSunk && ai.mishiSunk {
				maxLength = 2
			}
		}
		if ai.umiSunk {
			minLength = 3
			if ai.krakenSunk && ai.mishiSunk {
				minLength = 4
			}
		}

		minLength -= ai.currHitTally
		maxLength -= ai.currHitTally

		if ai.triedLeft != true {
			col := ai.prevCol - ai.searchIncrememnt
			if col >= 0 && col < grid.Columns {
				//Experimental====================================
				//Check Left
				if col != 0 {
					row := ai.prevRow

					for i := 2; i < maxLength; i++ {
						if col-i >= 0 {
							if grid.area[row][col-i] == TileEmptyShot || grid.area[row][col-i] == TileOccupiedShot {
								if i < minLength {
									leftBad = true
								} else {
									break
								}
							}
						}
					}

				}
				leftBad = false
				//================================================
				if !leftBad {
					miss, hit, sunk := grid.AIShot(ai.prevRow, col)
					ai.prevSunk = sunk
					if miss {
						ai.triedLeft = true
						ai.searchIncrememnt = 1
					} else if hit {
						ai.searchIncrememnt++
						ai.hitTally++
						ai.currHitTally++
						ai.hitLocations = append(ai.hitLocations, []int{ai.prevRow, col})
					} else {
						ai.triedLeft = true
						ai.searchIncrememnt = 1
					}
					return miss || hit
				}
			}
			ai.triedLeft = true
			ai.searchIncrememnt = 1

		}

		if ai.triedRight == false {
			col := ai.prevCol + ai.searchIncrememnt
			if col >= 0 && col < grid.Columns {
				//Experimental====================================
				//Check Right
				if col < grid.Columns {
					row := ai.prevRow

					for i := 2; i < maxLength; i++ {
						if col+i < grid.Columns {
							if grid.area[row][col+i] == TileEmptyShot || grid.area[row][col+i] == TileOccupiedShot {
								if i < minLength {
									rightBad = true
								} else {
									break
								}
							}
						}
					}

				}
				//================================================
				rightBad = false
				if !rightBad {
					miss, hit, sunk := grid.AIShot(ai.prevRow, col)
					ai.prevSunk = sunk
					if miss {
						ai.triedRight = true
						ai.searchIncrememnt = 1
					} else if hit {
						ai.searchIncrememnt++
						ai.hitTally++
						ai.currHitTally++
						ai.hitLocations = append(ai.hitLocations, []int{ai.prevRow, col})
					} else {
						ai.triedRight = true
						ai.searchIncrememnt = 1
					}
					return miss || hit
				}
			}
			ai.triedRight = true
			ai.searchIncrememnt = 1
		}

		if ai.triedUp == false {
			row := ai.prevRow - ai.searchIncrememnt
			if row >= 0 && row < grid.Rows {
				//Experimental====================================
				//Check CheckUp
				if row != 0 {
					col := ai.prevCol

					for i := 2; i < maxLength; i++ {
						if row-i >= 0 {
							if grid.area[row-i][col] == TileEmptyShot || grid.area[row-i][col] == TileOccupiedShot {
								if i < minLength {
									upBad = true
								} else {
									break
								}
							}
						}
					}

				}
				//================================================
				upBad = false
				if !upBad {
					miss, hit, sunk := grid.AIShot(row, ai.prevCol)
					ai.prevSunk = sunk
					if miss {
						ai.triedUp = true
						ai.searchIncrememnt = 1
					} else if hit {
						ai.searchIncrememnt++
						ai.hitTally++
						ai.currHitTally++
						ai.hitLocations = append(ai.hitLocations, []int{row, ai.prevCol})
					} else {
						ai.triedUp = true
						ai.searchIncrememnt = 1
					}
					return miss || hit
				}
			}
			ai.triedUp = true
			ai.searchIncrememnt = 1
		}

		if ai.triedDown == false {
			row := ai.prevRow + ai.searchIncrememnt
			if row >= 0 && row < grid.Rows {
				//Experimental====================================
				//Check Down
				if row < grid.Rows-1 {
					col := ai.prevCol

					for i := 2; i < maxLength; i++ {
						if row+i < grid.Rows {
							if grid.area[row+i][col] == TileEmptyShot || grid.area[row+i][col] == TileOccupiedShot {
								if i < minLength {
									downBad = true
								} else {
									break
								}
							}
						}
					}

				}
				//================================================
				downBad = false
				if !downBad {
					miss, hit, sunk := grid.AIShot(row, ai.prevCol)
					ai.prevSunk = sunk
					if miss {
						ai.triedDown = true
						ai.searchIncrememnt = 1
					} else if hit {
						ai.searchIncrememnt++
						ai.hitTally++
						ai.currHitTally++
						ai.hitLocations = append(ai.hitLocations, []int{row, ai.prevCol})
					} else {
						ai.triedDown = true
						ai.searchIncrememnt = 1
					}
					return miss || hit
				}
			}
			ai.triedDown = true
			ai.searchIncrememnt = 1
		}
		ai.resetTryDirections()
	}

	return false
}

func (ai *AI) resetTryDirections() {
	ai.triedDown = false
	ai.triedLeft = false
	ai.triedRight = false
	ai.triedUp = false
}
