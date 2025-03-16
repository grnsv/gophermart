package models

import "time"

type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

type User struct {
	ID       string `db:"id" json:"id"`
	Login    string `db:"login" json:"login"`
	Password string `db:"password" json:"-"`
}

type Order struct {
	ID         int       `db:"id" json:"number,string"`
	UserID     string    `db:"user_id" json:"-"`
	Status     Status    `db:"status" json:"status"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}
