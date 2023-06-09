package main

import (
	"torch/src/endpoints"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/status/java/:ip", endpoints.FetchJavaHandler)
	router.GET("/status/bedrock/:ip", endpoints.FetchBedrockHandler)
	router.GET("/srv/:host", endpoints.SrvHandler)
	router.GET("/icon/:ip", endpoints.IconHandler)
	router.GET("/ping", endpoints.PingHandler)

	router.Run(":8000")

}
