package isledef

import (
	"github.com/atolVerderben/tentsuyu"
	"github.com/hajimehoshi/ebiten"
)

//Character represents the player characters
type Character struct {
	*tentsuyu.BasicObject
	imgParts        *tentsuyu.BasicImageParts
	imgPartsBust    *tentsuyu.BasicImageParts
	imgPartsBustSad *tentsuyu.BasicImageParts
	Name            string
	scale           float64
}

func NewCharacter(name string) *Character {
	c := &Character{
		Name:  name,
		scale: 3,
	}

	switch name {
	case ghostgirl:
		c.BasicObject = &tentsuyu.BasicObject{
			X:           20,
			Y:           45,
			Width:       96,
			Height:      96,
			NotCentered: true,
		}
		c.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameGhostGirl].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameGhostGirl].Frame["y"],
			Width:  SpalooshSheet.Frames[frameGhostGirl].Frame["w"],
			Height: SpalooshSheet.Frames[frameGhostGirl].Frame["h"],
		}
		c.imgPartsBust = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameGhostGirl].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameGhostGirl].Frame["y"],
			Width:  SpalooshSheet.Frames[frameGhostGirl].Frame["w"],
			Height: SpalooshSheet.Frames[frameGhostGirl].Frame["h"] - SpalooshSheet.Frames[frameGhostGirl].Frame["h"]/3 + 3,
		}
		c.imgPartsBustSad = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameGhostGirlSad].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameGhostGirlSad].Frame["y"],
			Width:  SpalooshSheet.Frames[frameGhostGirlSad].Frame["w"],
			Height: SpalooshSheet.Frames[frameGhostGirlSad].Frame["h"] - SpalooshSheet.Frames[frameGhostGirl].Frame["h"]/3 + 3,
		}
	case vampire:
		c.BasicObject = &tentsuyu.BasicObject{
			X:           20,
			Y:           45,
			Width:       96,
			Height:      96,
			NotCentered: true,
		}
		c.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameVamp].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameVamp].Frame["y"],
			Width:  SpalooshSheet.Frames[frameVamp].Frame["w"],
			Height: SpalooshSheet.Frames[frameVamp].Frame["h"],
		}
		c.imgPartsBust = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameVamp].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameVamp].Frame["y"],
			Width:  SpalooshSheet.Frames[frameVamp].Frame["w"],
			Height: SpalooshSheet.Frames[frameVamp].Frame["h"] - SpalooshSheet.Frames[frameVamp].Frame["h"]/3 + 3,
		}
		c.imgPartsBustSad = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameVampSad].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameVampSad].Frame["y"],
			Width:  SpalooshSheet.Frames[frameVampSad].Frame["w"],
			Height: SpalooshSheet.Frames[frameVampSad].Frame["h"] - SpalooshSheet.Frames[frameVamp].Frame["h"]/3 + 3,
		}
	case hunter:
		c.BasicObject = &tentsuyu.BasicObject{
			X:           20,
			Y:           45,
			Width:       96,
			Height:      96,
			NotCentered: true,
		}
		c.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameVampHunter].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameVampHunter].Frame["y"],
			Width:  SpalooshSheet.Frames[frameVampHunter].Frame["w"],
			Height: SpalooshSheet.Frames[frameVampHunter].Frame["h"],
		}
		c.imgPartsBust = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameVampHunter].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameVampHunter].Frame["y"],
			Width:  SpalooshSheet.Frames[frameVampHunter].Frame["w"],
			Height: SpalooshSheet.Frames[frameVampHunter].Frame["h"] - SpalooshSheet.Frames[frameVampHunter].Frame["h"]/3 + 3,
		}
		c.imgPartsBustSad = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameVampHunterSad].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameVampHunterSad].Frame["y"],
			Width:  SpalooshSheet.Frames[frameVampHunterSad].Frame["w"],
			Height: SpalooshSheet.Frames[frameVampHunterSad].Frame["h"] - SpalooshSheet.Frames[frameVamp].Frame["h"]/3 + 3,
		}
	case nure:
		c.BasicObject = &tentsuyu.BasicObject{
			X:           615,
			Y:           290,
			Width:       SpalooshSheet.Frames[frameNureOnna].Frame["w"] * 3,
			Height:      SpalooshSheet.Frames[frameNureOnna].Frame["h"] * 3,
			NotCentered: true,
		}
		c.imgParts = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameNureOnna].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameNureOnna].Frame["y"],
			Width:  SpalooshSheet.Frames[frameNureOnna].Frame["w"],
			Height: SpalooshSheet.Frames[frameNureOnna].Frame["h"],
		}
		c.imgPartsBust = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameNureOnna].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameNureOnna].Frame["y"],
			Width:  SpalooshSheet.Frames[frameNureOnna].Frame["w"],
			Height: SpalooshSheet.Frames[frameNureOnna].Frame["h"] / 2,
		}
		c.imgPartsBustSad = &tentsuyu.BasicImageParts{
			Sx:     SpalooshSheet.Frames[frameNureLost].Frame["x"],
			Sy:     SpalooshSheet.Frames[frameNureLost].Frame["y"],
			Width:  SpalooshSheet.Frames[frameNureLost].Frame["w"],
			Height: SpalooshSheet.Frames[frameNureLost].Frame["h"] / 2,
		}
	}

	return c
}

//Draw the character
func (c *Character) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = c.imgParts
	op.GeoM.Scale(c.scale, c.scale)
	op.GeoM.Translate(c.X, c.Y)
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	return nil
}

//DrawBust the character
func (c *Character) DrawBust(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = c.imgPartsBust
	op.GeoM.Scale(c.scale, c.scale)
	//op.GeoM.Translate(c.X-float64(c.Width/2), c.Y-float64(c.Height/2))
	op.GeoM.Translate(c.X, c.Y)
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	return nil
}

//DrawBustSad the character
func (c *Character) DrawBustSad(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.ImageParts = c.imgPartsBustSad
	op.GeoM.Scale(c.scale, c.scale)
	//op.GeoM.Translate(c.X-float64(c.Width/2), c.Y-float64(c.Height/2))
	op.GeoM.Translate(c.X, c.Y)
	if err := screen.DrawImage(tentsuyu.ImageManager.ReturnImage("spaloosh-sheet"), op); err != nil {
		return err
	}
	return nil
}

func (c *Character) SetScale(scale float64) {
	c.scale = scale
}
