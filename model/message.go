package model

type Message struct {
	Time    int64  `json:"time"`
	Src     string `json:"src"`
	Dest    string `json:"dest"`
	Content string `json:"content"`
}
