package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	ErrUsernameExists = errors.New("username already exists")
	ErrUserNotFound   = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `
        INSERT INTO users (id, email, password, username, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password, user.Username, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// PostgreSQL error code 23505 = unique_violation
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_email_unique" {
					return ErrEmailExists
				}
				if pgErr.ConstraintName == "users_username_unique" {
					return ErrUsernameExists
				}
			}
		}
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT id, email, username FROM users WHERE id = $1`

	user := &User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("unable to query user: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, email, password, username FROM users WHERE username = $1`

	user := &User{}
	err := r.db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Email, &user.Password, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("unable to query user: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	query := `
        UPDATE users
        SET email = $2, password = $3, updated_at = $4	
        WHERE id = $1
    `
	_, err := r.db.Exec(ctx, query, user.ID, user.Email, user.Password, user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// PostgreSQL error code 23505 = unique_violation
			if pgErr.Code == "23505" {
				if pgErr.ConstraintName == "users_email_unique" {
					return ErrEmailExists
				}
			}
		}
		return fmt.Errorf("unable to update row: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("unable to delete row: %w", err)
	}
	return nil
}
