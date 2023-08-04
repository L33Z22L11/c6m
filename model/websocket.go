package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	connections map[*websocket.Conn]bool
	lock        sync.Mutex
}

var Connections = ConnectionManager{
	connections: make(map[*websocket.Conn]bool),
}
