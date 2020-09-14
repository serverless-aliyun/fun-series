package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type searchQuery struct {
	Keyword string `form:"keyword"`
	Details bool   `form:"details,default=false"`
}

type controller struct {
	svc *service
}

func (ctrl *controller) search(c *gin.Context) {
	var query searchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, failed(err.Error()))
	}

	series, err := ctrl.svc.search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(series))
}

func (ctrl *controller) detail(c *gin.Context) {
	seriesID := c.Param("seriesId")
	if seriesID == "" {
		c.JSON(http.StatusBadRequest, failed("missing path param 'seriesId'"))
		return
	}

	series, err := ctrl.svc.detail(seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(series))
}

func (ctrl *controller) episodes(c *gin.Context) {
	seriesID := c.Param("seriesId")
	if seriesID == "" {
		c.JSON(http.StatusBadRequest, failed("missing path param 'seriesId'"))
		return
	}

	episodes, err := ctrl.svc.episodes(seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(episodes))
}

func failed(msg string) gin.H {
	return gin.H{
		"msg":       msg,
		"timestamp": time.Now().Unix(),
	}
}

func data(data interface{}) gin.H {
	return gin.H{
		"msg":       "success",
		"data":      data,
		"timestamp": time.Now().Unix(),
	}
}
