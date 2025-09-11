package storage

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/khaleelsyed/codaVirtuale/internal/types"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *types.SugarWithTrace
}

func (s *PostgresStorage) CallNextTicket(deskID int) (types.Ticket, error) {
	query := `UPDATE ticket t
	SET desk_id = $1
	FROM desk d
	WHERE d.id = $1
	  AND t.category_id = d.category_id
	  AND t.closed = FALSE
	  AND t.desk_id IS NULL
	  AND t.id = (
	    SELECT id
	    FROM ticket
	    WHERE category_id = d.category_id
	      AND closed = FALSE
	      AND desk_id IS NULL
	    ORDER BY created_at
	    LIMIT 1
	  )
	RETURNING t.id, t.category_id, t.sub_url, t.desk_id, t.closed, t.created_at;`

	result, err := s.db.Exec(query, deskID)
	if err != nil {
		s.logger.Warnw("error with CallNextTicket", "desk_id", deskID, "error", err)
	}

	if err = checkSingleRowAffected(result, deskID, "CallNextTicket", s.logger); err != nil {
		if err == ErrNoRowsAffected {
			return types.Ticket{}, types.ErrnotFound
		}

		return types.Ticket{}, err
	}

	result.RowsAffected()

	return types.Ticket{}, types.ErrNotImplemented
}

func (s *PostgresStorage) SeeNext(categoryID int) (types.Ticket, error) {

	return types.Ticket{}, types.ErrNotImplemented
}

func (s *PostgresStorage) SeeQueue() ([]int, error) {
	return []int{}, types.ErrNotImplemented
}

func (s *PostgresStorage) CreateTicket(ticketCreate types.TicketCreate) (types.Ticket, error) {
	result, err := s.db.Query("INSERT INTO ticket (category_id, sub_url) VALUES ($1, $2) RETURNING id, category_id, sub_url, closed, created_at", ticketCreate.CategoryID, ticketCreate.SubURL)
	if err != nil {
		s.logger.Warnw("could not create ticket", "error", err)
		return types.Ticket{}, err
	}
	defer result.Close()

	var ticket types.Ticket

	for result.Next() {
		if err = result.Scan(&ticket.ID, &ticket.CategoryID, &ticket.SubURL, &ticket.Closed, &ticket.CreatedAt); err != nil {
			return types.Ticket{}, err
		}
		ticket.DeskID = -1
	}

	return ticket, nil
}

func (s *PostgresStorage) GetTicket(id int) (types.Ticket, error) {
	result, err := s.db.Query("SELECT id, category_id, sub_url, desk_id, closed, created_at FROM ticket WHERE id = $1", id)
	if err != nil {
		s.logger.Warnw("error with GetCategory", "error", err)
	}
	defer result.Close()

	var ticket types.Ticket

	if result.Next() {
		err = result.Scan(&ticket.ID, &ticket.CategoryID, &ticket.SubURL, &ticket.DeskID, &ticket.Closed, &ticket.CreatedAt)
		if err != nil {
			if err.Error() == errDeskNull.Error() {
				ticket.DeskID = -1
				return ticket, nil
			}
			return types.Ticket{}, err
		}
		return ticket, nil
	}

	return types.Ticket{}, types.ErrnotFound
}

func (s *PostgresStorage) DeleteTicket(id int) error {
	query := `DELETE FROM ticket
	WHERE id = $1;`

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Tracew("error deleting ticket", "id", id, "error", err)
		return err
	}

	return checkSingleRowAffected(result, id, "DeleteTicket", s.logger)
}

func (s *PostgresStorage) CreateCategory(name string) (types.Category, error) {
	result, err := s.db.Query("INSERT INTO category (name) VALUES ($1) RETURNING id, name", name)
	if err != nil {
		s.logger.Warnw("could not create category", "error", err)
		return types.Category{}, err
	}
	defer result.Close()

	var category types.Category

	for result.Next() {
		if err = result.Scan(&category.ID, &category.Name); err != nil {
			return types.Category{}, err
		}
	}

	return category, nil
}

func (s *PostgresStorage) GetCategory(id int) (types.Category, error) {
	result, err := s.db.Query("SELECT * FROM category WHERE id = $1", id)
	if err != nil {
		s.logger.Warnw("error with GetCategory", "error", err)
	}
	defer result.Close()

	var category types.Category

	if result.Next() {
		if err = result.Scan(&category.ID, &category.Name); err != nil {
			return types.Category{}, err
		}
		return category, nil
	}

	return types.Category{}, types.ErrnotFound
}

func (s *PostgresStorage) UpdateCategory(id int, name string) (types.Category, error) {
	var err error

	query := `UPDATE category
	SET name = $1
	WHERE id = $2;`

	result, err := s.db.Exec(query, name, id)
	if err != nil {
		s.logger.Tracew("error updating category", "id", id, "error", err)
		return types.Category{}, err
	}

	if err = checkSingleRowAffected(result, id, "UpdateCategory", s.logger); err != nil {
		return types.Category{}, err
	}
	return types.Category{ID: id, Name: name}, nil
}

func (s *PostgresStorage) DeleteCategory(id int) error {
	query := `DELETE FROM category
	WHERE id = $1;`

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Tracew("error deleting category", "id", id, "error", err)
		return err
	}

	return checkSingleRowAffected(result, id, "DeleteCategory", s.logger)
}

func (s *PostgresStorage) CreateDesk(label string, categoryID int) (types.Desk, error) {
	result, err := s.db.Query("INSERT INTO desk (label, category_id) VALUES ($1, $2) RETURNING id, label, category_id", label, categoryID)
	if err != nil {
		s.logger.Warnw("could not create desk", "error", err)
		return types.Desk{}, err
	}
	defer result.Close()

	var desk types.Desk

	for result.Next() {
		if err = result.Scan(&desk.ID, &desk.Label, &desk.CategoryID); err != nil {
			return types.Desk{}, err
		}
	}

	return desk, nil
}

func (s *PostgresStorage) GetDesk(id int) (types.Desk, error) {
	result, err := s.db.Query("SELECT id, category_id, label FROM desk WHERE id = $1", id)
	if err != nil {
		s.logger.Warnw("error with GetDesk Query", "error", err)
	}
	defer result.Close()

	var desk types.Desk

	if result.Next() {
		if err = result.Scan(&desk.ID, &desk.CategoryID, &desk.Label); err != nil {
			s.logger.Tracew("error with GetDesk Scanner", "id", id, "error", err)
			return types.Desk{}, err
		}
		return desk, nil
	}

	return types.Desk{}, types.ErrnotFound
}

func (s *PostgresStorage) UpdateDesk(id int, deskUpdate struct {
	CategoryID int
	Label      string
}) (types.Desk, error) {
	var err error

	query := `UPDATE desk
	SET category_id = $1, label = $2
	WHERE id = $3;`

	result, err := s.db.Exec(query, deskUpdate.CategoryID, deskUpdate.Label, id)
	if err != nil {
		s.logger.Tracew("error updating desk", "id", id, "category_id", deskUpdate.CategoryID, "error", err)
		return types.Desk{}, err
	}

	if err = checkSingleRowAffected(result, id, "UpdateDesk", s.logger); err != nil {
		return types.Desk{}, err
	}

	return types.Desk{ID: id, Label: deskUpdate.Label, CategoryID: deskUpdate.CategoryID}, nil

}

func (s *PostgresStorage) DeleteDesk(id int) error {
	query := `DELETE FROM desk
	WHERE id = $1;`

	result, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Tracew("error deleting desk", "id", id, "error", err)
		return err
	}

	return checkSingleRowAffected(result, id, "DeleteDesk", s.logger)
}

func checkSingleRowAffected(result sql.Result, id int, operation string, logger *types.SugarWithTrace) error {
	var err error

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Warnw(fmt.Sprintf("err while counting rows affected during %s", operation), "id", id, "error", err)
		return err
	}

	switch rowsAffected {
	case 1:
		return nil
	case 0:
		logger.Warnw(fmt.Sprintf("failed to perform %s", operation), "id", id, "RowsAffected", rowsAffected, "error", nil)
		return ErrNoRowsAffected
	default:
		err = errAffectedMultipleRows(operation)
		logger.Errorw(err.Error(), "id", id, "RowsAffected", rowsAffected)
		return err
	}
}

func (s *PostgresStorage) createCategoryTable() error {
	query := `CREATE TABLE IF NOT EXISTS category(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) createDeskTable() error {
	query := `CREATE TABLE IF NOT EXISTS desk(
	id SERIAL PRIMARY KEY,
	label VARCHAR(50) NOT NULL,
	category_id INT REFERENCES category(id) ON DELETE RESTRICT
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) createTicketTable() error {
	query := `CREATE TABLE IF NOT EXISTS ticket(
	id SERIAL PRIMARY KEY,
	category_id INT REFERENCES category(id),
	sub_url TEXT UNIQUE,
	desk_id INT REFERENCES desk(id),
	closed BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) createPreventCategoryDeleteOnOpenTickets() error {
	query := `CREATE OR REPLACE FUNCTION prevent_category_delete_on_open_tickets()
	RETURNS trigger AS $$
	BEGIN
	  -- Check if there are any tickets linked to this category that are not closed
	  IF EXISTS (
	    SELECT 1
	    FROM ticket t
	    WHERE t.category_id = OLD.id
	      AND t.closed = false
	  ) THEN
	    RAISE EXCEPTION 'Cannot delete category %: open tickets exist', OLD.id;
	  END IF;

	  RETURN OLD;
	END;
	$$ LANGUAGE plpgsql;

	CREATE OR REPLACE TRIGGER trg_prevent_category_delete
	BEFORE DELETE ON category
	FOR EACH ROW
	EXECUTE FUNCTION prevent_category_delete_on_open_tickets();`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) Init() error {
	var err error

	if err = s.createCategoryTable(); err != nil {
		s.logger.Errorw("unable to create `category` table", "error", err)
		return err
	}

	if err = s.createDeskTable(); err != nil {
		s.logger.Errorw("unable to create `desk` table", "error", err)
		return err
	}

	if err = s.createTicketTable(); err != nil {
		s.logger.Errorw("unable to create `ticket` table", "error", err)
		return err
	}

	if err = s.createPreventCategoryDeleteOnOpenTickets(); err != nil {
		s.logger.Errorw("unable to add function/trigger `prevent_category_delete_on_open_tickets`", "error", err)
		return err
	}

	return nil
}

func NewPostgresStorage(logger *types.SugarWithTrace) (*PostgresStorage, error) {

	connStr := os.Getenv("POSTGRES_CONN_STRING")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Errorw("failed to open db connection", "error", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Errorw("failed to ping db", "error", err)
		return nil, err
	}

	return &PostgresStorage{db: db, logger: logger}, nil
}
