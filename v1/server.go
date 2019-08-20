package v1

import (
	"github.com/gin-gonic/gin"
	"joy/config"
	"joy/v1/api"
	"net/http"
	"strconv"
)


var conf = config.Get()
func Start() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	r.GET("/status", api.Status)
	r.GET("/echo", Connect)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"template/index.tmpl",gin.H{
			"addr": "ws://" + conf.HOST + ":" + strconv.Itoa(conf.PORT) + "/echo",})
	})
	r.Run(conf.HOST + ":" + strconv.Itoa(conf.PORT)) // listen and serve on 0.0.0.0:8080
}