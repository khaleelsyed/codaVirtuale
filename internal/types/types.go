package types

import (
	"time"
)

type Desk struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Label      string `json:"label"`
}

type Ticket struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	SubURL     string    `json:"sub_url"`
	DeskID     int       `json:"desk_id"`
	Closed     bool      `json:"closed"`
	CreatedAt  time.Time `json:"created_at"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TicketCreate struct {
	CategoryID int
	SubURL     string
}
