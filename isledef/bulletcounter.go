package isledef

import "github.com/atolVerderben/tentsuyu"
import "github.com/hajimehoshi/ebiten"

type BulletCounter struct {
	*tentsuyu.BasicObject
	bombImgParts     *tentsuyu.BasicImageParts
	bullets          []bool
	remainingBullets int
	bulletPerRow     int
}

func NewBulletCounter(x, y float64, gd *GameData) *BulletCounter {
	bc := &BulletCounter{
		BasicObject: &tentsuyu.BasicObject{
			X:      x,
			Y:      y,
			Width:  100,
			Height: 100,
		},
		bombImgParts: &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameBomb].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameBomb].Frame["y"],
			Width:  32,
			Height: 32,
		},
		bullets:          make([]bool, gd.shotMax),
		bulletPerRow:     4,
		remainingBullets: gd.shotMax,
	}
	for i := range bc.bullets {
		bc.bullets[i] = true
	}
	return bc
}

func (bc *BulletCounter) Fire() {
	if bc.remainingBullets > 0 {
		bc.remainingBullets--
		bc.bullets[bc.remainingBullets] = false
	}
}

func (bc *BulletCounter) Draw(screen *ebiten.Image) error {
	bX, bY := 0.0, 0.0
	justEven := false
	justStart := true
	for i, b := range bc.bullets {
		i = i + 1
		if b {
			if bc.bulletPerRow == 2 {
				if i%2 == 0 {
					//this is an even number

					bX = 34.0

					justEven = true

				} else {

					bX = 0

					if justEven { // Move down to the next row
						bY += 32
					}
					justEven = false
				}
			}
			if bc.bulletPerRow == 4 {
				if i%4 == 0 {
					//this is an even number

					bX += 34.0

					justEven = true

				} else {

					if justEven { // Move down to the next row
						bY += 32
						bX = 0
					} else {
						if justStart {
							bX = 0
							justStart = false
						} else {
							bX += 34
						}
					}
					justEven = false
				}
			}
			op := &ebiten.DrawImageOptions{}
			op.ImageParts = bc.bombImgParts
			//op.ImageParts = g.missImg
			op.GeoM.Scale(ZoomLevel, ZoomLevel)
			op.GeoM.Translate(bc.X+bX, bc.Y+bY)

			if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
				return err
			}
		}
	}
	return nil
}
