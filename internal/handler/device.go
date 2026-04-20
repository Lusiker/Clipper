package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lusiker/clipper/internal/middleware"
	"github.com/lusiker/clipper/internal/service"
	"github.com/lusiker/clipper/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type DeviceHandler struct {
	deviceService *service.DeviceService
	hub           *ws.Hub
}

func NewDeviceHandler(deviceService *service.DeviceService, hub *ws.Hub) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService, hub: hub}
}

func (h *DeviceHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", middleware.AuthRequired(), h.List)
	r.GET("/ws", middleware.AuthRequired(), h.WebSocket)
}

func (h *DeviceHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)

	devices, err := h.deviceService.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"devices": devices})
}

func (h *DeviceHandler) WebSocket(c *gin.Context) {
	userID := middleware.GetUserID(c)
	deviceID := c.Query("device_id")
	deviceName := c.Query("device_name")

	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id required"})
		return
	}

	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	ip := c.ClientIP()
	device, err := h.deviceService.Register(userID, deviceID, deviceName, ip)
	if err != nil || device == nil {
		conn.Close()
		return
	}

	client := ws.NewClient(h.hub, conn, userID, deviceID, deviceName)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}