package config

import (
	"encoding/json"
	"fmt"
	"os"
)
type Configuration struct {
	HOST string
	PORT int
	WORLDNUM int
	USER_PER_ROOM int
	MESSAGE_BUFFER_SIZE int
	PLAY_TIME_SEC int
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