package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grnsv/gophermart/internal/models"
	_ "github.com/lib/pq"
)

type DBStorage struct {
	db                    *sql.DB
	isLoginExistsStmt     *sql.Stmt
	createUserStmt        *sql.Stmt
	findUserByLoginStmt   *sql.Stmt
	createOrderStmt       *sql.Stmt
	findOrderByIDStmt     *sql.Stmt
	getOrdersByUserIDStmt *sql.Stmt
}

func New(ctx context.Context, dsn string) (Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	storage := &DBStorage{db: db}
	err = storage.initStmt(ctx)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *DBStorage) initStmt(ctx context.Context) error {
	var err error
	if s.isLoginExistsStmt, err = s.db.PrepareContext(ctx, `
		SELECT EXISTS(SELECT * FROM users WHERE login = $1) AS exists
	`); err != nil {
		return err
	}
	if s.createUserStmt, err = s.db.PrepareContext(ctx, `
		INSERT INTO users (id, login, password)
		VALUES ($1, $2, $3)
	`); err != nil {
		return err
	}
	if s.findUserByLoginStmt, err = s.db.PrepareContext(ctx, `
		SELECT * FROM users WHERE login = $1 LIMIT 1
	`); err != nil {
		return err
	}
	if s.createOrderStmt, err = s.db.PrepareContext(ctx, `
		INSERT INTO orders (id, user_id, status)
		VALUES ($1, $2, $3)
	`); err != nil {
		return err
	}
	if s.findOrderByIDStmt, err = s.db.PrepareContext(ctx, `
		SELECT * FROM orders WHERE id = $1 LIMIT 1
	`); err != nil {
		return err
	}
	if s.getOrdersByUserIDStmt, err = s.db.PrepareContext(ctx, `
		SELECT * FROM orders WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`); err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) Close() error {
	if err := s.isLoginExistsStmt.Close(); err != nil {
		return err
	}
	if err := s.createUserStmt.Close(); err != nil {
		return err
	}
	if err := s.findUserByLoginStmt.Close(); err != nil {
		return err
	}
	if err := s.createOrderStmt.Close(); err != nil {
		return err
	}
	if err := s.findOrderByIDStmt.Close(); err != nil {
		return err
	}
	if err := s.getOrdersByUserIDStmt.Close(); err != nil {
		return err
	}
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) IsLoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	if err := s.isLoginExistsStmt.QueryRowContext(ctx, login).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (s *DBStorage) CreateUser(ctx context.Context, user *models.User) error {
	if _, err := s.createUserStmt.ExecContext(ctx, user.ID, user.Login, user.Password); err != nil {
		return err
	}
	return nil
}

func (s *DBStorage) FindUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	if err := s.findUserByLoginStmt.QueryRowContext(ctx, login).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *DBStorage) CreateOrder(ctx context.Context, order *models.Order) error {
	_, err := s.createOrderStmt.ExecContext(ctx, order.ID, order.UserID, order.Status)
	return err
}

func (s *DBStorage) FindOrderByID(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	if err := s.findOrderByIDStmt.QueryRowContext(ctx, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.UploadedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (s *DBStorage) GetOrdersByUserID(ctx context.Context, userID string) ([]*models.Order, error) {
	var orders []*models.Order
	rows, err := s.getOrdersByUserIDStmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.UploadedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, ErrNotFound
	}

	return orders, nil
}

func (s *DBStorage) UpdateBalance(userID string, balance float64) error {
	panic("unimplemented")
}
