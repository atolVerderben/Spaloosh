package isledef

import (
	"image/color"

	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

type Icon struct {
	X, Y     float64
	imgParts *tentsuyu.BasicImageParts
}

type EnemyDisplay struct {
	x, y                                                         float64
	krakenDefeated, umibozuDefeated, bakeDefeated, mishiDefeated bool
	smallGrid                                                    bool
	krakenText, umibozuText, bakeText, mishiText                 *tentsuyu.TextElement
}

func NewEnemyDisplay(x, y float64, smallGrid bool) *EnemyDisplay {
	ed := &EnemyDisplay{
		x:           x,
		y:           y,
		smallGrid:   smallGrid,
		umibozuText: tentsuyu.NewTextElementStationary(x, y+73, 100, 50, Game.UIController.ReturnFont(FntSmallPixel), []string{"Umibozu"}, color.Black, 12),
		krakenText:  tentsuyu.NewTextElementStationary(x, y+193, 100, 50, Game.UIController.ReturnFont(FntSmallPixel), []string{"Kraken"}, color.Black, 12),
		bakeText:    tentsuyu.NewTextElementStationary(x, y+255, 100, 50, Game.UIController.ReturnFont(FntSmallPixel), []string{"Bake-Kujira"}, color.Black, 12),
		mishiText:   tentsuyu.NewTextElementStationary(x, y+131, 100, 50, Game.UIController.ReturnFont(FntSmallPixel), []string{"Mishipeshu"}, color.Black, 12),
	}
	if smallGrid {
		ed.krakenText.SetPosition(x, y+131)
		ed.bakeText.SetPosition(x, y+193)
	}
	return ed
}

func (ed *EnemyDisplay) Update(ships []*Ship) {
	for _, s := range ships {
		if s.dead {
			switch s.shipType {
			case ShipTypeBattleShip:
				ed.bakeDefeated = true
			case ShipTypeCruiser:
				ed.krakenDefeated = true
			case ShipTypePatrol:
				ed.umibozuDefeated = true
			case ShipTypeSubmarine:
				ed.mishiDefeated = true
			}
		}
	}
}

func (ed *EnemyDisplay) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameUmibozu].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameUmibozu].Frame["y"],
		Width:  SpalooshSheet.Frames[frameUmibozu].Frame["w"],
		Height: SpalooshSheet.Frames[frameUmibozu].Frame["h"],
	}
	op.GeoM.Translate(-float64(32/2), -float64(64/2))
	op.GeoM.Rotate(1.5708)
	op.GeoM.Translate(float64(32/2), float64(64/2)) //Switch the width and height because it's now vertical
	op.GeoM.Translate(ed.x, ed.y)
	if ed.umibozuDefeated {
		op.ColorM.Scale(0, 0, 0, 1)
	}
	if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	yOffset := 96.0
	if !ed.smallGrid {
		op = &ebiten.DrawImageOptions{}
		op.ImageParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameMishipeshu].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameMishipeshu].Frame["y"],
			Width:  SpalooshSheet.Frames[frameMishipeshu].Frame["w"],
			Height: SpalooshSheet.Frames[frameMishipeshu].Frame["h"],
		}
		op.GeoM.Translate(ed.x, ed.y+yOffset)
		if ed.mishiDefeated {
			op.ColorM.Scale(0, 0, 0, 1)
		}
		if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
			return err
		}
		yOffset += 64
	}
	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameKraken].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameKraken].Frame["y"],
		Width:  SpalooshSheet.Frames[frameKraken].Frame["w"],
		Height: SpalooshSheet.Frames[frameKraken].Frame["h"],
	}
	op.GeoM.Translate(ed.x, ed.y+yOffset)
	if ed.krakenDefeated {
		op.ColorM.Scale(0, 0, 0, 1)
	}
	if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	yOffset += 64
	op = &ebiten.DrawImageOptions{}
	op.ImageParts = &tentsuyu.BasicImageParts{
		Sx:     SpalooshSheet.Frames[frameBake].Frame["x"],
		Sy:     SpalooshSheet.Frames[frameBake].Frame["y"],
		Width:  SpalooshSheet.Frames[frameBake].Frame["w"],
		Height: SpalooshSheet.Frames[frameBake].Frame["h"],
	}
	op.GeoM.Translate(ed.x, ed.y+yOffset)
	if ed.bakeDefeated {
		op.ColorM.Scale(0, 0, 0, 1)
	}
	if err := screen.DrawImage(Game.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}

	ed.krakenText.Draw(screen)
	ed.umibozuText.Draw(screen)
	ed.bakeText.Draw(screen)
	if !ed.smallGrid {
		ed.mishiText.Draw(screen)
	}
	return nil
}
