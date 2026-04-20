package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	UserID   string
	DeviceID string
	DeviceName string
	Send     chan []byte
	mu       sync.Mutex
}

func NewClient(hub *Hub, conn *websocket.Conn, userID, deviceID, deviceName string) *Client {
	return &Client{
		Hub:        hub,
		Conn:       conn,
		UserID:     userID,
		DeviceID:   deviceID,
		DeviceName: deviceName,
		Send:       make(chan []byte, 256),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		c.mu.Lock()
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		c.mu.Unlock()

		if err != nil {
			break
		}
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}