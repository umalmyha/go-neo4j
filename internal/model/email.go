package model

type Email struct {
	From    string
	To      []string
	Cc      []string
	Content string
}
