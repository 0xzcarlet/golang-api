package response

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Error: message})
}

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, SuccessResponse{Data: data})
}

func SuccessMessage(c *gin.Context, code int, message string) {
	c.JSON(code, SuccessResponse{Message: message})
}
