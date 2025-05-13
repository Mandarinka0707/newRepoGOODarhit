package websocket

import (
	"net/http"

	"github.com/Mandarinka0707/newRepoGOODarhit/chat/internal/entity"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan entity.Message)
