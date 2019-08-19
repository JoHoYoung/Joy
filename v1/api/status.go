package api

import (
	"github.com/gin-gonic/gin"
	"joy/World"
	"net/http"
)


/*
{
	user:
	Rooms: [
	]
}
 */
func Status(c *gin.Context) {
	runningRooms := []int{}
	for _, room := range World.Rooms{
		if room.Running{
			runningRooms = append(runningRooms, room.Id);
		}
	}
	c.JSON(http.StatusOK,gin.H{"Rooms":runningRooms, "User":World.NumberOfUser});
}