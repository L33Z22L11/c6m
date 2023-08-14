package model

type Group struct {
	Gid     string   `json:"gid"`
	Gname   string   `json:"gname"`
	Owner   string   `json:"owner"`
	Admins  []string `json:"admins"`
	Members []string `json:"members"`
}
