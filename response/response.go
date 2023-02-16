package response

import (
	"github.com/labstack/echo/v4"
)

type JSONSuccessResult struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data"`
}

// SuccessResponse => for custom response
func SuccessResponse(c echo.Context, data interface{}, status int) error {
	c.JSON(status, JSONSuccessResult{
		Code:    status,
		Message: "Success",
		Data:    data,
	})
	return nil
}
