package api

import "github.com/khaleelsyed/codaVirtuale/internal/storage"

type Storage interface {
	CallNextTicket(deskID int) (storage.Ticket, error)
	LastCalled(categoryID int, positions int) ([]storage.Ticket, error)
	SeeNext(categoryID int) (storage.Ticket, error)
	SeeQueue() ([]int, error)

	CreateTicket(category storage.Category) (storage.Ticket, error)

	CreateCategory(name string) (storage.Category, error)
	GetCategory(id int) (storage.Category, error)
	UpdateCategory(id int, name string) (storage.Category, error)
	DeleteCategory(id int) error

	CreateDesk(label string, category storage.Category) (storage.Desk, error)
	GetDesk(id int) (storage.Desk, error)
	UpdateDesk(id int, label string) (storage.Desk, error)
	DeleteDesk(id int) error
}
