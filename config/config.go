package config

import (
	"encoding/json"
	"fmt"
	"os"
)
type Type struct{
	ENTER int `json:"ENTER"`
	START int `json:"START"`
	END int `json:"END"`
	LEAVE int `json:"LEAVE"`
	PLAY int `json:"PLAY"`
	REGEN int `json:"REGEN"`
	STEAL int `json:"STEAL"`
	IN int `json:"IN"`
}
type Configuration struct {
	HOST string `json:"HOST"`
	PORT int `json:"PORT"`
	WORLDNUM int `json:"WORLDNUM"`
	USER_PER_ROOM int `json:"USER_PER_ROOM"`
	MESSAGE_BUFFER_SIZE int `json:"MESSAGE_BUFFER_SIZE"`
	PLAY_TIME_SEC int `json:"PLAY_TIME_SEC"`
	SEND_TIME_SEC int `json:"SEND_TIME_SEC"`
	REGEN_TIME_SEC int `json:"REGEN_TIME_SEC"`
	TYPE Type `json:"TYPE"`
}

var Config Configuration

func Get() *Configuration {
	file, _ := os.Open("config/config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Config)
	fmt.Println(Config)
	if err != nil {
		fmt.Println("error:", err)
	}
	return &Config
}