package main

import (
	"github.com/gin-gonic/gin"
)

var HTTPServerURL = "0.0.0.0:5000"


func main() {
	defer GetDB().Close()

	router := gin.Default()
	router.POST("/aggregation", Aggregation)
	router.POST("/deaggregation", DeAggregation)

	router.Run("localhost:7777")
}
