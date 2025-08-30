package storage

import (
	"errors"
	"fmt"
	"time"
)

var errNotImplemented = errors.New("feature not implemented yet")

type MockStorage struct{}

func (s MockStorage) CallNextTicket(desk Desk) (Ticket, error) {

	return Ticket{
		ID:          2,
		CategoryID:  desk.CategoryID,
		SubURL:      "frjikll23l",
		QueueNumber: 2,
		DeskID:      desk.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) LastCalled(category Category, positions int) ([]Ticket, error) {
	return []Ticket{}, errNotImplemented
}

func (s MockStorage) SeeNext(category Category) (Ticket, error) {
	return Ticket{}, errNotImplemented
}

func (s MockStorage) SeeQueue(category Category) ([]Ticket, error) {
	return []Ticket{}, errNotImplemented
}

func (s MockStorage) CreateTicket(category Category) (Ticket, error) {
	return Ticket{}, errNotImplemented
}

func (s MockStorage) CreateCategory(name string) (Category, error) {
	return Category{}, errNotImplemented
}

func (s MockStorage) GetCategory(id int) (Category, error) {
	return Category{}, errNotImplemented
}

func (s MockStorage) UpdateCategory(id int, name string) (Category, error) {
	return Category{}, errNotImplemented
}

func (s MockStorage) DeleteCategory(id int) error {
	return errNotImplemented
}

func (s MockStorage) CreateDesk(label string, category Category) (Desk, error) {
	return Desk{}, errNotImplemented
}

func (s MockStorage) GetDesk(id int) (Desk, error) {
	return Desk{
		ID:         id,
		CategoryID: 1,
		Label:      fmt.Sprintf("desk %i", id),
	}, nil
}

func (s MockStorage) UpdateDesk(id int, label string) (Desk, error) {
	return Desk{}, errNotImplemented
}

func (s MockStorage) DeleteDesk(id int) error {
	return errNotImplemented
}

func NewMockStorage() (MockStorage, error) {
	return MockStorage{}, nil
}
