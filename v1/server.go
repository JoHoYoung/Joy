package v1

import (
	"github.com/gin-gonic/gin"
	"joy/Config"
	"joy/v1/api"
	"net/http"
	"strconv"
)


var config = Config.Get()
func Start() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*")

	r.GET("/status", api.Status)

	r.GET("/echo", Connect)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"template/index.tmpl",gin.H{
			"addr": "ws://" + config.HOST + ":" + strconv.Itoa(config.PORT) + "/echo",})
	})
	r.Run(config.HOST + ":" + strconv.Itoa(config.PORT)) // listen and serve on 0.0.0.0:8080

}