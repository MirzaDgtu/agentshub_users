package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"agentshub_users/internal/domain/models"
)

// UserRepository реализует интерфейс репозитория для работы с пользователями.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password, first_name, middle_name, last_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.Password,
		user.FirstName,
		user.MiddleName,
		user.LastName,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUser возвращает пользователя по его ID.
func (r *UserRepository) GetUser(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, email, first_name, middle_name, last_name, profile_image_url, is_blocked, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.ProfileImageURL,
		&user.IsBlocked,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// UpdateUser обновляет данные пользователя.
func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET
			first_name = $1,
			middle_name = $2,
			last_name = $3,
			profile_image_url = $4,
			updated_at = $5
		WHERE id = $6
	`

	_, err := r.db.Exec(ctx, query,
		user.FirstName,
		user.MiddleName,
		user.LastName,
		user.ProfileImageURL,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser удаляет пользователя по его ID.
func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers возвращает список пользователей с пагинацией.
func (r *UserRepository) ListUsers(ctx context.Context, limit int, offset int) ([]*models.User, error) {
	query := `
		SELECT id, email, first_name, middle_name, last_name, profile_image_url, is_blocked, created_at, updated_at
		FROM users
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.MiddleName,
			&user.LastName,
			&user.ProfileImageURL,
			&user.IsBlocked,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

// BlockUser блокирует пользователя.
func (r *UserRepository) BlockUser(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_blocked = true WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to block user: %w", err)
	}

	return nil
}

// UnblockUser разблокирует пользователя.
func (r *UserRepository) UnblockUser(ctx context.Context, id int64) error {
	query := `UPDATE users SET is_blocked = false WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to unblock user: %w", err)
	}

	return nil
}
