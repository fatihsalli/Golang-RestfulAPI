package response

import (
	"github.com/labstack/echo/v4"
)

type JSONSuccessResult struct {
	TotalItemCount int         `json:"totalitemcount"`
	Data           interface{} `json:"data"`
}

// array için totalItem count kullanılmadılır!!!

// SuccessArrayResponse => for custom response
func SuccessArrayResponse(c echo.Context, data interface{}, totalItemCount int) error {
	c.JSON(200, JSONSuccessResult{
		TotalItemCount: totalItemCount,
		Data:           data,
	})
	return nil
}
