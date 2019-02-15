package server

import (
	"fmt"
	"github.com/labstack/echo"
	"time"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

func returnError(code int, controller string, err error) error {
	if err != nil {
		fmt.Printf("controller %s error: %s\n", controller, err.Error())
	} else {
		fmt.Printf("controller %s error: null\n", controller)
	}

	return echo.NewHTTPError(code, fmt.Sprintf("something went wrong in [%s]", controller))
}
