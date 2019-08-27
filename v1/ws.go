package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"joy/world"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func Connect(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer,c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	world.NewClient(conn)
}