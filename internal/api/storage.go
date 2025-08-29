package api

import "github.com/khaleelsyed/codaVirtuale/internal/storage"

type Storage interface {
	CallNextTicket(desk storage.Desk) (storage.Ticket, error)
	LastCalled(category storage.Category, positions int) ([]storage.Ticket, error)
	SeeNext(category storage.Category) (storage.Ticket, error)
	SeeQueue(category storage.Category) ([]storage.Ticket, error)

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
