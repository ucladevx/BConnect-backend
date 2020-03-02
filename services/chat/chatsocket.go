package chat

import (
	"github.com/gorilla/websocket"
	"github.com/ucladevx/BConnect-backend/models"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan models.Message)    // broadcast channel
