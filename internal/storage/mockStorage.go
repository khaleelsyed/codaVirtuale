package storage

import (
	"fmt"
	"time"

	"github.com/khaleelsyed/codaVirtuale/internal/types"
)

type MockStorage struct{}

func (s MockStorage) CallNextTicket(deskID int) (types.Ticket, error) {

	return types.Ticket{
		ID:          2,
		CategoryID:  4,
		SubURL:      "frjikll23l",
		QueueNumber: 2,
		DeskID:      deskID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) LastCalled(categoryID int, positions int) ([]types.Ticket, error) {
	tickets := []types.Ticket{
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

func (s MockStorage) SeeNext(categoryID int) (types.Ticket, error) {

	return types.Ticket{
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

func (s MockStorage) CreateTicket(categoryID int) (types.Ticket, error) {
	return types.Ticket{
		ID:          8,
		CategoryID:  categoryID,
		SubURL:      "hjkl8",
		QueueNumber: 8,
		DeskID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (s MockStorage) GetTicket(ticketID int) (types.Ticket, error) {
	return types.Ticket{
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

func (s MockStorage) CreateCategory(name string) (types.Category, error) {
	return types.Category{ID: 1, Name: name}, nil
}

func (s MockStorage) GetCategory(id int) (types.Category, error) {
	return types.Category{
		ID:   id,
		Name: fmt.Sprintf("Desk %d", id),
	}, nil
}

func (s MockStorage) UpdateCategory(id int, name string) (types.Category, error) {
	return types.Category{ID: id, Name: name}, nil
}

func (s MockStorage) DeleteCategory(id int) error {
	return nil
}

func (s MockStorage) CreateDesk(label string, categoryID int) (types.Desk, error) {
	return types.Desk{
		ID:         1,
		CategoryID: categoryID,
		Label:      label,
	}, nil
}

func (s MockStorage) GetDesk(id int) (types.Desk, error) {
	return types.Desk{
		ID:         id,
		CategoryID: 1,
		Label:      fmt.Sprintf("desk %d", id),
	}, nil
}

func (s MockStorage) UpdateDesk(id int, deskUpdate struct {
	CategoryID int
	Label      string
}) (types.Desk, error) {
	if deskUpdate.CategoryID == 0 {
		deskUpdate.CategoryID = 1
	}

	if deskUpdate.Label == "" {
		deskUpdate.Label = fmt.Sprintf("desk %d", id)
	}
	return types.Desk{
		ID:         id,
		CategoryID: deskUpdate.CategoryID,
		Label:      deskUpdate.Label,
	}, nil
}

func (s MockStorage) DeleteDesk(id int) error {
	return nil
}

func NewMockStorage() (MockStorage, error) {
	return MockStorage{}, nil
}

func (s MockStorage) Init() error {
	return nil
}
