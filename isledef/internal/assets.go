package assets

import (
	"encoding/base64"
	"image"
	"image/png"

	"github.com/golang/freetype/truetype"

	"github.com/hajimehoshi/ebiten"
)

import "bytes"

//DecodeB64 decodes a base64 string to byte slice
func DecodeB64(message string) []byte {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, _ := base64.StdEncoding.Decode(base64Text, []byte(message))
	//log.Printf("%s\n", base64Text[:l])
	return base64Text[:l]
}

//LoadSpalooshSheet loads the sprite sheet for the game
func LoadSpalooshSheet() (*ebiten.Image, error) {

	img, err := png.Decode(bytes.NewReader(DecodeB64(spalooshSheet)))
	if err != nil {
		return nil, err
	}
	img2, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return nil, err
	}
	return img2, err
}

func LoadImage(imgString string) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(DecodeB64(imgString)))
	if err != nil {
		return nil
	}
	img2, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return nil
	}
	return img2
}

func LoadSpriteSheetJSON() []byte {
	return ([]byte(spalooshSheetJSON))
}

func ReturnSpalooshSE() []byte {
	return DecodeB64(spalooshWav)
}

func ReturnKaboomSE() []byte {
	return DecodeB64(kaboomWav)
}

func ReturnSE(seName string) []byte {
	switch seName {
	case "kaboom-ah2.wav":
		return ReturnKaboomSE()
	case "spaloosh.wav":
		return ReturnSpalooshSE()

	}
	return nil
}

func ReturnPixelFont() *truetype.Font {
	tt, _ := truetype.Parse(DecodeB64(smallPixelFont))
	/*if err != nil {
		return err
	}*/
	return tt
}
