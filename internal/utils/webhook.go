package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func RegisterHandlers(r *echo.Group, errorLog *log.Logger, infoLog *log.Logger, env map[string]string, redisClient *redis.Client) {
	res := Resource{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Env:           env,
		RedisClient:   redisClient,
	}

	r.POST("/utils/webhook", res.webhook)
}

type Resource struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Env           map[string]string
	RedisClient   *redis.Client
}

func (r Resource) webhook(c echo.Context) error {
	fmt.Println("-------------------------------------------------------------------------------------------------")
	fmt.Println("Webhook triggered")
	fmt.Println("--------------------------------------------------------------------------------------------------")

	// Calculate the time until 00:00 PST
	now := time.Now()
	pst, _ := time.LoadLocation("America/Los_Angeles")
	desiredTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pst)
	timeDifference := desiredTime.Sub(now)

	// Set the key with expiration time
	err := r.RedisClient.Set(c.Request().Context(), "quotaExceeded", "true" , timeDifference).Err()
	if err != nil {
		r.ErrorLog.Println("Error setting key in Redis:", err)
		// Handle the error as needed
	}
	
	return c.JSON(http.StatusOK, echo.Map{
		"token": "tok",
	})
}
