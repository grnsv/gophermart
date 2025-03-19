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
	db                         *sql.DB
	isLoginExistsStmt          *sql.Stmt
	createUserStmt             *sql.Stmt
	findUserByLoginStmt        *sql.Stmt
	createOrderStmt            *sql.Stmt
	findOrderByIDStmt          *sql.Stmt
	getOrdersByUserIDStmt      *sql.Stmt
	updateOrderStmt            *sql.Stmt
	getBalanceStmt             *sql.Stmt
	createWithdrawalStmt       *sql.Stmt
	getWithdrawalsByUserIDStmt *sql.Stmt
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
		SELECT * FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`); err != nil {
		return err
	}
	if s.updateOrderStmt, err = s.db.PrepareContext(ctx, `
		UPDATE orders
		SET status = $1, accrual = $2
		WHERE id = $3
	`); err != nil {
		return err
	}
	if s.getBalanceStmt, err = s.db.PrepareContext(ctx, `
		WITH
			orders_total AS (
				SELECT COALESCE(SUM(orders.accrual), 0) AS accrued
				FROM orders
				WHERE orders.user_id = $1
			),
			withdrawals_total AS (
				SELECT COALESCE(SUM(withdrawals.sum), 0) AS withdrawn
				FROM withdrawals
				WHERE withdrawals.user_id = $2
			)
		SELECT
			orders_total.accrued - withdrawals_total.withdrawn AS current,
			withdrawals_total.withdrawn
		FROM
			orders_total, withdrawals_total;
	`); err != nil {
		return err
	}
	if s.createWithdrawalStmt, err = s.db.PrepareContext(ctx, `
		INSERT INTO withdrawals (user_id, order_id, sum)
		VALUES ($1, $2, $3)
	`); err != nil {
		return err
	}
	if s.getWithdrawalsByUserIDStmt, err = s.db.PrepareContext(ctx, `
		SELECT order_id, sum, processed_at FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
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
	if err := s.updateOrderStmt.Close(); err != nil {
		return err
	}
	if err := s.getBalanceStmt.Close(); err != nil {
		return err
	}
	if err := s.createWithdrawalStmt.Close(); err != nil {
		return err
	}
	if err := s.getWithdrawalsByUserIDStmt.Close(); err != nil {
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
		&order.Accrual,
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
			&order.Accrual,
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

func (s *DBStorage) UpdateOrder(ctx context.Context, order *models.Order) error {
	_, err := s.updateOrderStmt.ExecContext(ctx, order.Status, order.Accrual, order.ID)
	return err
}

func (s *DBStorage) GetBalance(ctx context.Context, userID string) (*models.Balance, error) {
	var balance models.Balance
	if err := s.getBalanceStmt.QueryRowContext(ctx, userID, userID).Scan(
		&balance.Current,
		&balance.Withdrawn,
	); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (s *DBStorage) CreateWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) error {
	_, err := s.createWithdrawalStmt.ExecContext(ctx, withdrawal.UserID, withdrawal.OrderID, withdrawal.Sum)
	return err
}

func (s *DBStorage) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]*models.Withdrawal, error) {
	var withdrawals []*models.Withdrawal
	rows, err := s.getWithdrawalsByUserIDStmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal models.Withdrawal
		if err := rows.Scan(
			&withdrawal.OrderID,
			&withdrawal.Sum,
			&withdrawal.ProcessedAt,
		); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, &withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, ErrNotFound
	}

	return withdrawals, nil
}
