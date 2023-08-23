package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Group, errorLog *log.Logger, infoLog *log.Logger, env map[string]string) {
	res := Resource{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Env:           env,
	}

	r.POST("/utils/webhook", res.webhook)
}

type Resource struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Env           map[string]string
}

func (r Resource) webhook(c echo.Context) error {
	fmt.Println("-------------------------------------------------------------------------------------------------")
	fmt.Println("Webhook triggered")
	fmt.Println("--------------------------------------------------------------------------------------------------")
	return c.JSON(http.StatusOK, echo.Map{
		"token": "tok",
	})
}
