package model

type Email struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Cc      []string `json:"cc"`
	Content string   `json:"content"`
}
