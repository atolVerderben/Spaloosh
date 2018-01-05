package isledef

import (
	"encoding/json"
	"io/ioutil"
)

//config is the configuration file
type config struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

//readConfigFile reads the configuration json file...
func readConfigFile(filename string) *config {
	m := &config{}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		m.Port = "5555"
		m.Server = "127.0.0.1"
	} else {
		json.Unmarshal(raw, m)
	}

	return m
}
