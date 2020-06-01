package chat

import (
	"sync"

	"github.com/ucladevx/BConnect-backend/models"
	"github.com/ucladevx/BConnect-backend/utils/websockets"
)

/*
Chats employs a subscription model for users and rooms; users subscribe to rooms.
Chats is (testing) concurrency safe.
For now, the way to open chat rooms is through test.html, which was ported from an external
source
*/

//Rooms chat
type Rooms struct {
	roomLock   sync.RWMutex
	Rooms      map[string]map[*websockets.Websocket]bool
	Broadcast  chan models.Message
	Register   chan websockets.Subscription
	Unregister chan websockets.Subscription
}

//NewRooms creates chat room object
func NewRooms() Rooms {
	return Rooms{
		Broadcast:  make(chan models.Message),
		Register:   make(chan websockets.Subscription),
		Unregister: make(chan websockets.Subscription),
		Rooms:      make(map[string]map[*websockets.Websocket]bool),
	}
}

//Run runs rooms
func (r *Rooms) Run() {
	for {
		select {
		case sub := <-r.Register:
			r.roomLock.RLock()
			connections := r.Rooms[sub.RoomID]
			r.roomLock.RUnlock()
			if connections == nil {
				connections = make(map[*websockets.Websocket]bool)
				r.roomLock.Lock()
				r.Rooms[sub.RoomID] = connections
				r.roomLock.Unlock()
			}
			r.roomLock.Lock()
			r.Rooms[sub.RoomID][sub.Socket] = true
			r.roomLock.Unlock()
		case sub := <-r.Unregister:
			r.roomLock.RLock()
			connections := r.Rooms[sub.RoomID]
			r.roomLock.RUnlock()
			if connections != nil {
				if _, ok := connections[sub.Socket]; ok {
					delete(connections, sub.Socket)
					close(sub.Socket.MessageChannel)
					if len(connections) == 0 {
						delete(r.Rooms, sub.RoomID)
					}
				}
			}
		case message := <-r.Broadcast:
			r.roomLock.RLock()
			connections := r.Rooms[message.MessageRoom]
			r.roomLock.RUnlock()
			for c := range connections {
				select {
				case c.MessageChannel <- message.Message:
				default:
					close(c.MessageChannel)
					delete(connections, c)
					if len(connections) == 0 {
						delete(r.Rooms, message.MessageRoom)
					}
				}
			}
		}
	}
}
