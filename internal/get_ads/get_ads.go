package get_ads

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-retry"
	"gopkg.in/guregu/null.v4"
)

type AdQuery struct {
	Offset   int      `form:"offset" binding:"min=0,max=100"`
	Limit    int      `form:"limit,default=5" binding:"min=1,max=100"`
	Age      null.Int `form:"age" binding:"min=1,max=100"`
	Gender   string   `form:"gender" binding:"omitempty,oneof= M F"`
	Country  string   `form:"country" binding:"omitempty,iso3166_1_alpha2"`
	Platform string   `form:"platform" binding:"omitempty,oneof= android ios web"`
}

type Ad struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}

type AdResponse struct {
	Items []Ad `json:"items"`
}

func GetAds(c *gin.Context) {
	var q AdQuery
	if err := c.BindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := checkAge(q.Age); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=ad port=5432 sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	condition := `
		($1::int IS NULL OR age_start IS NULL OR age_start <= $1)
		AND ($1::int IS NULL OR age_end IS NULL OR $1 <= age_end)
		AND start_at <= NOW() AND NOW() <= end_at
		AND ($2 = '' OR gender IS NULL OR $2 = ANY(gender))
		AND ($3 = '' OR country IS NULL OR $3 = ANY(country))
		AND ($4 = '' OR platform IS NULL OR $4 = ANY(platform))`
	sqlStatement := fmt.Sprintf("SELECT title, end_at FROM ad WHERE %s ORDER BY end_at ASC LIMIT $6 OFFSET $5", condition)

	// try to query db with maximum 3 trials
	ctx := context.Background()
	b := retry.WithMaxRetries(3, retry.NewFibonacci(100 * time.Millisecond))
	if err := retry.Do(ctx, b, func(ctx context.Context) error {
		rows, err := db.Query(sqlStatement, q.Age, q.Gender, q.Country, q.Platform, q.Offset, q.Limit)
		if err != nil {
			return retry.RetryableError(err)
		}
		defer rows.Close()

		ads := []Ad{}
		for rows.Next() {
			var ad Ad
			if err = rows.Scan(&ad.Title, &ad.EndAt); err != nil {
				return retry.RetryableError(err)
			}
			ads = append(ads, ad)
		}
		adRsp := AdResponse{Items: ads}

		c.AsciiJSON(http.StatusOK, adRsp)
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

}

func checkAge(age null.Int) error {
	if age.Valid && (age.Int64 < 1 || age.Int64 > 100) {
		return errors.New("age should be between 1 and 100")
	}
	return nil
}
