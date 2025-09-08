package storage

import (
	"database/sql"
	"os"

	"github.com/khaleelsyed/codaVirtuale/internal/types"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func (s *PostgresStorage) CallNextTicket(deskID int) (types.Ticket, error) {

	return types.Ticket{}, types.ErrNotImplemented
}

func (s *PostgresStorage) LastCalled(categoryID int, positions int) ([]types.Ticket, error) {
	tickets := []types.Ticket{}

	return tickets[:positions], types.ErrNotImplemented
}

func (s *PostgresStorage) SeeNext(categoryID int) (types.Ticket, error) {

	return types.Ticket{}, types.ErrNotImplemented
}

func (s *PostgresStorage) SeeQueue() ([]int, error) {
	return []int{}, types.ErrNotImplemented
}

func (s *PostgresStorage) CreateTicket(categoryID int) (types.Ticket, error) {
	return types.Ticket{}, types.ErrNotImplemented
}

func (s *PostgresStorage) GetTicket(ticketID int) (types.Ticket, error) {
	return types.Ticket{}, types.ErrNotImplemented
}

func (s *PostgresStorage) DeleteTicket(ticketID int) error {
	return types.ErrNotImplemented
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
	return types.Category{}, types.ErrNotImplemented
}

func (s *PostgresStorage) DeleteCategory(id int) error {
	return types.ErrNotImplemented
}

func (s *PostgresStorage) CreateDesk(label string, categoryID int) (types.Desk, error) {
	return types.Desk{}, types.ErrNotImplemented
}

func (s *PostgresStorage) GetDesk(id int) (types.Desk, error) {
	return types.Desk{}, types.ErrNotImplemented
}

func (s *PostgresStorage) UpdateDesk(id int, deskUpdate struct {
	CategoryID int
	Label      string
}) (types.Desk, error) {
	return types.Desk{}, types.ErrNotImplemented
}

func (s *PostgresStorage) DeleteDesk(id int) error {
	return types.ErrNotImplemented
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
	category_id SERIAL REFERENCES category(id) ON DELETE RESTRICT
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) createTicketTable() error {
	query := `CREATE TABLE IF NOT EXISTS ticket(
	id SERIAL PRIMARY KEY,
	category_id SERIAL REFERENCES category(id),
	sub_url TEXT UNIQUE,
	queue_number SERIAL,
	desk_id SERIAL REFERENCES desk(id),
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

func NewPostgresStorage(logger *zap.SugaredLogger) (*PostgresStorage, error) {

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
