package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data interface{} `json:"data"`
}

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{
		Data: data,
	})
}
