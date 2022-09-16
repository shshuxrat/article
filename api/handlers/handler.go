package handlers

import (
	"article/config"
	"article/storage"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	strg storage.StorageI
	cfg  config.Config
}

func NewHandler(strg storage.StorageI, cfg config.Config) Handler {
	return Handler{
		strg: strg,
		cfg:  cfg,
	}
}

func (h *Handler) getOffsetParam(c *gin.Context) (int, error) {
	offsetStr := c.DefaultQuery("offset", h.cfg.DefaultOffset)
	return strconv.Atoi(offsetStr)
}

func (h *Handler) getLimitParam(c *gin.Context) (int, error) {
	limitStr := c.DefaultQuery("offset", h.cfg.DefaultLimit)
	return strconv.Atoi(limitStr)
}
