package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, itemsController *ItemsController) {
	router.GET("/items/list", itemsController.List)
	router.POST("/item/create", itemsController.Create)
	router.PATCH("/item/update", itemsController.Update)
	router.DELETE("/item/remove", itemsController.Delete)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
