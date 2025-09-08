package api

import "github.com/khaleelsyed/codaVirtuale/internal/types"

type Storage interface {
	CallNextTicket(deskID int) (types.Ticket, error)
	LastCalled(categoryID int, positions int) ([]types.Ticket, error)
	SeeNext(categoryID int) (types.Ticket, error)
	SeeQueue() ([]int, error)

	CreateTicket(categoryID int) (types.Ticket, error)
	GetTicket(ticketID int) (types.Ticket, error)
	DeleteTicket(ticketID int) error

	CreateCategory(name string) (types.Category, error)
	GetCategory(id int) (types.Category, error)
	UpdateCategory(id int, name string) (types.Category, error)
	DeleteCategory(id int) error

	CreateDesk(label string, categoryID int) (types.Desk, error)
	GetDesk(id int) (types.Desk, error)
	UpdateDesk(id int, deskUpdate struct {
		CategoryID int
		Label      string
	}) (types.Desk, error)
	DeleteDesk(id int) error
}
