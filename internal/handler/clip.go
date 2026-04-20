package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lusiker/clipper/internal/middleware"
	"github.com/lusiker/clipper/internal/model"
	"github.com/lusiker/clipper/internal/service"
	"github.com/lusiker/clipper/internal/ws"
)

type ClipHandler struct {
	clipService *service.ClipService
	hub         *ws.Hub
}

func NewClipHandler(clipService *service.ClipService, hub *ws.Hub) *ClipHandler {
	return &ClipHandler{clipService: clipService, hub: hub}
}

func (h *ClipHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", middleware.AuthRequired(), h.List)
	r.POST("", middleware.AuthRequired(), h.Create)
	r.POST("/upload", middleware.AuthRequired(), h.Upload)
	r.GET("/:id", middleware.AuthRequired(), h.Get)
	r.DELETE("/:id", middleware.AuthRequired(), h.Delete)
}

func (h *ClipHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	clips, err := h.clipService.List(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clips": clips})
}

func (h *ClipHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID := c.Query("device_id")

	var req model.ClipCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clip, err := h.clipService.Create(userID, deviceID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.NotifyClipCreated(userID, clip, deviceID)

	c.JSON(http.StatusOK, gin.H{"clip": clip})
}

func (h *ClipHandler) Get(c *gin.Context) {
	userID := middleware.GetUserID(c)
	clipID := c.Param("id")

	clip, err := h.clipService.Get(userID, clipID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if clip == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clip not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clip": clip})
}

func (h *ClipHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	clipID := c.Param("id")
	deviceID := c.Query("device_id")

	if err := h.clipService.Delete(userID, clipID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.NotifyClipDeleted(userID, clipID, deviceID)

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ClipHandler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID := c.Query("device_id")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image file provided"})
		return
	}

	clip, err := h.clipService.UploadImage(userID, deviceID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.NotifyClipCreated(userID, clip, deviceID)

	c.JSON(http.StatusOK, gin.H{"clip": clip})
}