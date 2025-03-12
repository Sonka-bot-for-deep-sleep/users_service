package repository

import (
	"context"
	"fmt"

	"github.com/Sonka-bot-for-deep-sleep/user_service/application/models"
	"github.com/jackc/pgx/v5"
)

type user struct {
	db *pgx.Conn
}

func NewUser(db *pgx.Conn) *user {
	return &user{
		db: db,
	}
}

func (u *user) GetByTgID(ctx context.Context, tgID string) (*models.User, error) {
	var user models.User
	query := "SELECT ID, Tg_ID, Name, Login FROM users WHERE Tg_ID = @tgID"
	args := pgx.NamedArgs{
		"tgID": tgID,
	}

	row := u.db.QueryRow(ctx, query, args)

	if err := row.Scan(&user.ID, &user.TgId, &user.Name, &user.Login); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("GetUserByTgID: failed get user by id, user not found: %w", err)
		}
		return nil, fmt.Errorf("GetUserByTgID: failed get user by tg id: %w", err)
	}

	return &user, nil
}

func (u *user) CreateUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO users (Tg_ID, Name, Login) 
	VALUES ($1, $2, $3)`
	_, err := u.db.Exec(ctx, query, user.TgId, user.Name, user.Login)
	if err != nil {
		return fmt.Errorf("CreateUser: failed create user: %w", err)
	}

	return nil
}
