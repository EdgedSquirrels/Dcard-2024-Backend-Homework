package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"
)

type Advertisement struct {
	Title      string    `binding:"required"`
	StartAt    time.Time `binding:"required"`
	EndAt      time.Time `binding:"required"`
	Conditions struct {
		AgeStart null.Int `binding:"min=1,max=100"`
		AgeEnd   null.Int `binding:"min=1,max=100"`
		Gender   []string `binding:"dive,oneof= M F"`
		Country  []string `binding:"dive,iso3166_1_alpha2"`
		Platform []string `binding:"dive,oneof= android ios web"`
	} `binding:"required"`
}

type ReqAd struct {
	Offset   int         `form:"offset" binding:"min=0,max=100"`
	Limit    int         `form:"limit,default=5" binding:"min=1,max=100"`
	Age      null.Int    `form:"age" binding:"min=0,max=100"` // 0 as null
	Gender   string `form:"gender"`
	Country  string `form:"country"`
	Platform string `form:"platform"`
}

type ResAd struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}

func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/ad", GetAds)
	v1.POST("/ad", PostAds)
	r.Run(":8080")
}

func GetAds(c *gin.Context) {
	db, err := sql.Open("postgres", "dbname=ad sslmode=disable")
	checkErr(err)

	var data ReqAd
	if err = c.BindQuery(&data); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(data)

	fmt.Println("Successfully created connection to database")

	condition := `
		($1::int IS NULL OR age_start IS NULL OR age_start <= $1)
		AND ($1::int IS NULL OR age_end IS NULL OR $1 <= age_end)
		AND start_at <= NOW() AND NOW() <= end_at
		AND ($2 = '' OR gender IS NULL OR $2 = ANY(gender))
		AND ($3 = '' OR country IS NULL OR $3 = ANY(country))
		AND ($4 = '' OR platform IS NULL OR $4 = ANY(platform))
	`

	sqlStatement := fmt.Sprintf("SELECT title, end_at FROM ad WHERE %s ORDER BY end_at ASC LIMIT $6 OFFSET $5", condition)

	rows, err := db.Query(sqlStatement, data.Age, data.Gender, data.Country, data.Platform, data.Offset, data.Limit)
	checkErr(err)

	resData := []ResAd{}

	for rows.Next() {
		var resAd ResAd
		err = rows.Scan(&resAd.Title, &resAd.EndAt)
		checkErr(err)
		resData = append(resData, resAd)
	}

	// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	c.AsciiJSON(http.StatusOK, resData)
}

func PostAds(c *gin.Context) {
	var data Advertisement
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(data)

	db, err := sql.Open("postgres", "dbname=ad sslmode=disable")
	checkErr(err)

	sqlStatement := `
		INSERT INTO ad (title, start_at, end_at, age_start, age_end, gender, country, platform)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	conditions := data.Conditions
	_, err = db.Exec(sqlStatement, data.Title, data.StartAt, data.EndAt,
		conditions.AgeStart, conditions.AgeEnd,
		pq.Array(conditions.Gender), pq.Array(conditions.Country), pq.Array(conditions.Platform))
	checkErr(err)

	// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	c.AsciiJSON(http.StatusOK, data)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}