package post_ads

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
