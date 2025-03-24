package models

import "time"

type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
	StatusRegistered Status = "REGISTERED"
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
	Accrual    float64   `db:"accrual" json:"accrual,omitempty"`
	UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
}

type Withdrawal struct {
	ID          int       `db:"id" json:"-"`
	UserID      string    `db:"user_id" json:"-"`
	OrderID     int       `db:"order_id" json:"order,string"`
	Sum         float64   `db:"sum" json:"sum"`
	ProcessedAt time.Time `db:"processed_at" json:"processed_at,omitempty"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
