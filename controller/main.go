package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var (
	webPort = flag.String("port", "8080", "The port for the web server")
	local   = flag.Bool("local", true, "Use local or cloud worker")
)

func main() {
	flag.Parse()

	router := gin.Default()
	router.GET("/", pingApi)
	router.POST("/job", postJob)
	router.GET("/job/:id", getJob)

	webAddress := ":" + *webPort

	router.Run(webAddress)
}
