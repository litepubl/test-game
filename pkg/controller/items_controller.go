package controller

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/litepubl/test-game/pkg/entity"
	"github.com/litepubl/test-game/pkg/items"
)

type ItemsController struct {
	items items.CRUDItems
}

func NewItemsController(i items.CRUDItems) *ItemsController {
	return &ItemsController{
		items: i,
	}
}

func (ic ItemsController) List(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	json, err := ic.items.List(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", json)
}

func (ic ItemsController) Create(c *gin.Context) {
	campaignId, err := ic.campaignId(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "campaignId argument not specified")
		return
	}

	b := createBody{}
	err = c.BindJSON(&b)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" {
		errorResponse(c, http.StatusBadRequest, "Name argument not specified")
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	item, err := ic.items.Create(ctx, campaignId, b.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, "service unavailable")

		return
	}

	c.JSON(http.StatusOK, item)
}

func (ic ItemsController) Update(c *gin.Context) {
	u := entity.UpdateData{}
	var err error
	u.Id, err = ic.id(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "id argument not specified")
		return
	}

	u.CampaignId, err = ic.campaignId(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "campaignId argument not specified")
		return
	}

	err = c.BindJSON(&u)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		errorResponse(c, http.StatusBadRequest, "Name argument not specified")
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	item, err := ic.items.Update(ctx, u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	c.JSON(http.StatusOK, item)
}

func (ic ItemsController) Delete(c *gin.Context) {
	id, err := ic.id(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "id argument not specified")
		return
	}

	campaignId, err := ic.campaignId(c)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "campaignId argument not specified")
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	err = ic.items.Delete(ctx, id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, "service unavailable")
		return
	}

	c.JSON(http.StatusOK, removedResponse{
		Id:         id,
		CampaignId: campaignId,
		Removed:    true,
	})
}

func (ic ItemsController) id(c *gin.Context) (int, error) {
	return strconv.Atoi(c.Query("id"))
}

func (ic ItemsController) campaignId(c *gin.Context) (int, error) {
	return strconv.Atoi(c.Query("campaignId"))
}
