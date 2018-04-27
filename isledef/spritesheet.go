package isledef

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

//SpriteSheet holds all frames of a spritesheet from a json
type SpriteSheet struct {
	Frames []*Frame `json:"frames"`
}

//Frame represents a single frame of a spritesheet
type Frame struct {
	Filename         string             `json:"filename"`
	Frame            map[string]int     `json:"frame"`
	Rotated          bool               `json:"rotated"`
	Trimmed          bool               `json:"trimmed"`
	SpriteSourceSize map[string]int     `json:"spriteSourceSize"`
	SourceSize       map[string]int     `json:"sourceSize"`
	Pivot            map[string]float64 `json:"pivot"`
}

//ReadSpriteSheet reads a json file and returns a SpriteSheet struct
func ReadSpriteSheet(filename string) *SpriteSheet {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	m := &SpriteSheet{}
	json.Unmarshal(raw, m)
	return m
}

func ReadSpriteSheetJSON(jsonByte []byte) *SpriteSheet {
	m := &SpriteSheet{}
	json.Unmarshal(jsonByte, m)
	return m
}
