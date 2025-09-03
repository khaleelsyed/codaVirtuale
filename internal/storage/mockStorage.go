package storage

import (
	"errors"
	"fmt"
	"time"
)

var errNotImplemented = errors.New("feature not implemented yet")

type MockStorage struct{}

func (s MockStorage) CallNextTicket(deskID int) (Ticket, error) {

	return Ticket{
		ID:          2,
		CategoryID:  4,
		SubURL:      "frjikll23l",
		QueueNumber: 2,
		DeskID:      deskID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) LastCalled(categoryID int, positions int) ([]Ticket, error) {
	tickets := []Ticket{
		{
			ID:          1,
			CategoryID:  categoryID,
			SubURL:      "hjkl1",
			QueueNumber: 1,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			CategoryID:  categoryID,
			SubURL:      "hjkl2",
			QueueNumber: 2,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          3,
			CategoryID:  categoryID,
			SubURL:      "hjkl1",
			QueueNumber: 3,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          4,
			CategoryID:  categoryID,
			SubURL:      "hjkl4",
			QueueNumber: 4,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          5,
			CategoryID:  categoryID,
			SubURL:      "hjkl5",
			QueueNumber: 5,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          6,
			CategoryID:  categoryID,
			SubURL:      "hjkl6",
			QueueNumber: 6,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          7,
			CategoryID:  categoryID,
			SubURL:      "hjkl7",
			QueueNumber: 1,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          8,
			CategoryID:  categoryID,
			SubURL:      "hjkl8",
			QueueNumber: 8,
			DeskID:      1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return tickets[:positions], nil
}

func (s MockStorage) SeeNext(categoryID int) (Ticket, error) {

	return Ticket{
		ID:          2,
		CategoryID:  categoryID,
		SubURL:      "frjikll23l",
		QueueNumber: 2,
		DeskID:      4,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) SeeQueue() ([]int, error) {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, nil
}

func (s MockStorage) CreateTicket(categoryID int) (Ticket, error) {
	return Ticket{
		ID:          8,
		CategoryID:  categoryID,
		SubURL:      "hjkl8",
		QueueNumber: 8,
		DeskID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) GetTicket(ticketID int) (Ticket, error) {
	return Ticket{
		ID:          ticketID,
		CategoryID:  4,
		SubURL:      "hjkl8",
		QueueNumber: ticketID,
		DeskID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) DeleteTicket(ticketID int) error {
	return nil
}

func (s MockStorage) CreateCategory(name string) (Category, error) {
	return Category{}, errNotImplemented
}

func (s MockStorage) GetCategory(id int) (Category, error) {
	return Category{
		ID:   id,
		Name: fmt.Sprintf("Desk %d", id),
	}, nil
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
		Label:      fmt.Sprintf("desk %d", id),
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
