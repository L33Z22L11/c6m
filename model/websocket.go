package model

import "github.com/gorilla/websocket"

var Connections map[string]*websocket.Conn
