package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Domain-level errors
var (
	ErrDuplicateUser = errors.New(
		"user with this email or phone already exists",
	)
	ErrUserNotFound = errors.New(
		"user not found",
	)

	ErrForeignKey = errors.New(
		"foreign key constraint violated",
	)
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	StoreRefreshToken(ctx context.Context, userId int64, tokenHash string, ExpiresAt time.Time) error
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresRepository{
		pool: pool,
	}
}

func handleQueryError(err error) error {

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("no rows returned: %w", err)
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		// Unique constraint violation
		case "23505":
			return ErrDuplicateUser

		// Foreign key constraint violation
		case "23503":
			return ErrForeignKey

		// Other PostgreSQL errors
		default:
			return fmt.Errorf(
				"postgres error [%s]: %s",
				pgErr.Code,
				pgErr.Message,
			)
		}
	}

	// Non-PostgreSQL error
	return fmt.Errorf("database error: %w", err)
}

// 5. Create User
func (r *postgresRepository) CreateUser(
	ctx context.Context,
	user *User,
) error {
	query := `
		INSERT INTO users (
			full_name,
			email,
			phone,
			password_hash
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, role
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role,
	)

	if err != nil {
		return handleQueryError(err)
	}

	return nil
}
func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `
		SELECT 
			id,
			full_name,
			email,
			phone,
			password_hash,
			role,
			created_at,
			updated_at
		
			FROM users where email=$1
	`
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *postgresRepository) StoreRefreshToken(ctx context.Context, userId int64, tokenHash string, expiresAt time.Time) error {
	query := `
	
	INSERT INTO refresh_tokens(
		user_id,
		token_hash,
		expires_at
		)
		VALUES($1,$2,$3)
	`
	_, err := r.pool.Exec(ctx, query, userId, tokenHash, expiresAt)
	return err
}
