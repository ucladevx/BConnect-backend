package websockets

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ucladevx/BConnect-backend/errors"
	"github.com/ucladevx/BConnect-backend/models"
)

/*
Package for handling websockets and subscriptions to a connection
*/

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingTime       = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Websocket directly interfaces with websocket
type Websocket struct {
	Conn           *websocket.Conn
	MessageChannel chan []byte
}

//Subscription subscribes to a room
type Subscription struct {
	Socket *Websocket
	RoomID string
}

//NewWebsocket creates new interface with websocket
func NewWebsocket(conn *websocket.Conn) *Websocket {
	return &Websocket{
		Conn:           conn,
		MessageChannel: make(chan []byte),
	}
}

// Write writes a message
func (ws *Websocket) Write(messageType int, buff []byte) error {
	ws.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.Conn.WriteMessage(messageType, buff)
}

//Upgrade creates new interface with websocket on upgraded endpoint
func Upgrade(w http.ResponseWriter, r *http.Request) (*Websocket, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return NewWebsocket(ws), &errors.UpgradeError{
			Issue: err,
		}
	}
	return NewWebsocket(ws), nil
}

//NewSubscription creates new subscription
func NewSubscription(socket *Websocket, roomID string) Subscription {
	return Subscription{
		Socket: socket,
		RoomID: roomID,
	}
}

//SocketReader reads from websocket connection
func (s Subscription) SocketReader(broadcast chan models.Message, unregister chan Subscription) {
	sock := s.Socket
	defer func() {
		unregister <- s
		sock.Conn.Close()
	}()
	sock.Conn.SetReadLimit(maxMessageSize)
	sock.Conn.SetReadDeadline(time.Now().Add(pongWait))
	sock.Conn.SetPongHandler(func(string) error { sock.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, buf, err := s.Socket.Conn.ReadMessage()
		if err != nil {
			break
		}
		msg := models.Message{
			Message:     buf,
			MessageRoom: s.RoomID,
		}
		broadcast <- msg
	}
}

//SocketWriter writes to websocket connection
func (s Subscription) SocketWriter() {
	sock := s.Socket
	ticker := time.NewTicker(pingTime)
	defer func() {
		ticker.Stop()
		sock.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-sock.MessageChannel:
			if !ok {
				sock.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := sock.Write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := sock.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
