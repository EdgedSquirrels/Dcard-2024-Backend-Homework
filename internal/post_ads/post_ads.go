package post_ads

import (
	"database/sql"
	"errors"
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
		AgeStart null.Int
		AgeEnd   null.Int
		Gender   []string `binding:"dive,oneof= M F"`
		Country  []string `binding:"dive,iso3166_1_alpha2"`
		Platform []string `binding:"dive,oneof= android ios web"`
	} `binding:"required"`
}

func checkAge(ageStart, ageEnd null.Int) error {
	if ageStart.Valid && ageEnd.Valid && ageStart.Int64 > ageEnd.Int64 {
		return errors.New("age start should be less than age end")
	}
	if ageStart.Valid && (ageStart.Int64 < 1 || ageStart.Int64 > 100) {
		return errors.New("age start should be between 1 and 100")
	}
	if ageEnd.Valid && (ageEnd.Int64 < 1 || ageEnd.Int64 > 100) {
		return errors.New("age end should be between 1 and 100")
	}
	return nil
}

func PostAds(c *gin.Context) {
	var data Advertisement
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := checkAge(data.Conditions.AgeStart, data.Conditions.AgeEnd); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=ad port=5432 sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	sqlStatement := `
		INSERT INTO ad (title, start_at, end_at, age_start, age_end, gender, country, platform)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	conditions := data.Conditions

	_, err = db.Exec(sqlStatement, data.Title, data.StartAt, data.EndAt,
		conditions.AgeStart, conditions.AgeEnd,
		pq.Array(conditions.Gender), pq.Array(conditions.Country), pq.Array(conditions.Platform))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
