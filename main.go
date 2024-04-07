package main

import (
	"database/sql"
	"dcard2024/internal/get_ads"
	"dcard2024/internal/post_ads"

	"github.com/gin-gonic/gin"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SetUpAds() {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=ad port=5432 sslmode=disable")
	checkErr(err)

	_, err = db.Exec("DROP TABLE IF EXISTS ad;")
	checkErr(err)

	_, err = db.Exec(`CREATE TABLE ad (
		title text,
		start_at timestamp,
		end_at timestamp,
		age_start int,
		age_end int,
		gender text[],
		country text[],
		platform text[]
	  );`)
	checkErr(err)
	db.Close()
}

func setupRouter() *gin.Engine {
	SetUpAds()
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	v1 := r.Group("/api/v1")
	v1.GET("/ad", get_ads.GetAds)
	v1.POST("/ad", post_ads.PostAds)
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
