package user

import (
	"context"
	"database/sql"

	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

type Repository interface {
	Get(ctx context.Context, userId string) (model.User, error)
	InsertTx(ctx context.Context, tx *sql.Tx, user model.User) (uint64, error)
	Delete(ctx context.Context, userId string) error
	UpdateTx(ctx context.Context, tx *sql.Tx, user model.User) (uint64, error)
}

type repository struct {
	db *sql.DB
}

// Creates a new instance of the user repo
func NewRepository(db *sql.DB) Repository {
	return repository{
		db: db,
	}
}

// Calls get_user and returns a User object
func (r repository) Get(ctx context.Context, userId string) (model.User, error) {
	rows, err := r.db.QueryContext(ctx, "call get_user(?)", userId)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	var user model.User
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.UserId); err != nil {
			return model.User{}, err
		}
	}

	return user, nil
}

// Inserts a user as part of a transaction
func (r repository) InsertTx(ctx context.Context, tx *sql.Tx, user model.User) (uint64, error) {
	rows, err := tx.QueryContext(ctx, "call ins_user(?, ?, ?)", user.FirstName, user.LastName, user.UserId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id uint64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}

// Deletes a user
func (r repository) Delete(ctx context.Context, userId string) error {
	_, err := r.db.QueryContext(ctx, "call del_user(?)", userId)
	return err
}

// Updates a user as part of a transacion
func (r repository) UpdateTx(ctx context.Context, tx *sql.Tx, user model.User) (uint64, error) {
	rows, err := tx.QueryContext(ctx, "call upd_user(?, ?, ?)", user.FirstName, user.LastName, user.UserId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id uint64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}

	return id, nil
}
