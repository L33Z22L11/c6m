package models

type Message struct {
	Type    string `json:"type"`
	Time    string `json:"time"`
	Src     string `json:"src"`
	Dest    string `json:"dest"`
	Content string `json:"content"`
}

func ParseMsg(msg *Message) {
	switch msg.Type {
	case "":
	}
}
