package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"joy/World"
	"log"
)

var upgrader = websocket.Upgrader{} // use default options

func Connect(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer,c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	user := World.User{Name:"test"}
	World.NewClient(conn, &user)
}