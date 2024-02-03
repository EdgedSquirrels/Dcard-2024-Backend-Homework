package main

import (
	"dcard2024/internal/get_ads"
	"dcard2024/internal/post_ads"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/ad", get_ads.GetAds)
	v1.POST("/ad", post_ads.PostAds)
	r.Run(":8080")
}
