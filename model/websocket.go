package model

import "github.com/gorilla/websocket"

type Message struct {
	Type    string `json:"type"`
	Time    string `json:"time"`
	Src     string `json:"src"`
	Dest    string `json:"dest"`
	Content string `json:"content"`
}

var Connections map[string]*websocket.Conn
