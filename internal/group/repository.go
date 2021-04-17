package group

import (
	"context"
	"database/sql"

	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

type Repository interface {
	Get(ctx context.Context, groupName string) (model.Group, error)
	Insert(ctx context.Context, group model.Group) (uint64, error)
	Delete(ctx context.Context, groupName string) error
}

type repository struct {
	db *sql.DB
}

// Creates a new instance of group repository
func NewRepository(db *sql.DB) Repository {
	return repository{
		db: db,
	}
}

// Calls the get_group sp and returns a Group object
func (r repository) Get(ctx context.Context, groupName string) (model.Group, error) {
	rows, err := r.db.QueryContext(ctx, "call get_group(?)", groupName)
	if err != nil {
		return model.Group{}, err
	}
	defer rows.Close()

	var group model.Group
	for rows.Next() {
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			return model.Group{}, err
		}
	}

	return group, nil
}

// Calls ins_group sp and returns the id of that row
func (r repository) Insert(ctx context.Context, group model.Group) (uint64, error) {
	rows, err := r.db.QueryContext(ctx, "call ins_group(?)", group.Name)
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

// Calls the del_group sp
func (r repository) Delete(ctx context.Context, groupName string) error {
	_, err := r.db.QueryContext(ctx, "call del_group(?)", groupName)
	return err
}
