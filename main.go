package main

import (
	"github.com/gin-gonic/gin"
	"./get_ads"
	"./internal/post_ads"
)

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/ad", GetAds)
	v1.POST("/ad", PostAds)
	r.Run(":8080")
}
