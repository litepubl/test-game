package controller

import (
	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Code    int               `json:"code" example:"3"`
	Message string            `json:"message" example:"Item not found"`
	Details map[string]string `json:"details" example:"{}"`
}

func errorResponse(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, errResponse{
		Code:    code,
		Message: msg,
		Details: map[string]string{},
	})
}
