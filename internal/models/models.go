package models

import "time"

type User struct {
	ID       int
	Login    string
	Password string
}

type Chat struct {
	ID    int
	Owner int
}

type Message struct {
	ID       int
	ChatID   int
	UserID   int
	Text     string
	TimeSent time.Time
}
