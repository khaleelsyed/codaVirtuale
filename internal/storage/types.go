package storage

import "time"

type Desk struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Label      string `json:"label"`
}

type Ticket struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	SubURL      string    `json:"sub_url"`
	QueueNumber int       `json:"queue_number"` // Can be reset anytime (daily, every 12 hours ...)
	DeskID      int       `json:"desk_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
