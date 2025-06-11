package response

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Data interface{} `json:"data"`
}

func Error(c *gin.Context, status int, err error, message string) {
	errorRes := map[string]interface{}{
		"message": message,
	}

	c.JSON(status, ErrorResponse{Data: errorRes})
}
