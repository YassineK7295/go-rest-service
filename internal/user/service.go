package user

import (
	"context"
	"database/sql"

	"github.com/yassinekhaliqui/go-rest-service/internal/membership"
	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

type Service interface {
	GetWithGroup(ctx context.Context, userId string) (model.User, *[]model.Group, error)
	InsertTx(ctx context.Context, user model.User, groupNames *[]string) error
	Delete(ctx context.Context, userId string) error
	UpdateTx(ctx context.Context, user model.User, groupNames *[]string) error
}

type service struct {
	repo              Repository
	membershipService membership.Service
	db                *sql.DB
}

// Creates a new instance of the user service
func NewService(db *sql.DB) Service {
	return service{NewRepository(db), membership.NewService(db), db}
}

// Gets the user and their groups
func (s service) GetWithGroup(ctx context.Context, userId string) (model.User, *[]model.Group, error) {
	user, err := s.repo.Get(ctx, userId)
	if err != nil {
		return model.User{}, &[]model.Group{}, nil
	}

	groups, err := s.membershipService.GetGroupsForUser(ctx, user.Id)
	if err != nil {
		return model.User{}, &[]model.Group{}, nil
	}

	return user, groups, nil
}

// Inserts the user and their links to groups in a transaction
func (s service) InsertTx(ctx context.Context, user model.User, groupNames *[]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = func() error {
		userId, err := s.repo.InsertTx(ctx, tx, user)
		if err != nil {
			return err
		}

		if groupNames != nil && len(*groupNames) != 0 {
			err = s.membershipService.InsertTx(ctx, tx, userId, groupNames)
			if err != nil {
				return err
			}
		}

		return tx.Commit()
	}()

	if err != nil {
		tx.Rollback()
	}

	return err
}

// Deletes a user and their links to groups
func (s service) Delete(ctx context.Context, userId string) error {
	return s.repo.Delete(ctx, userId)
}

// Updates the user and their links to groups in a transaction
func (s service) UpdateTx(ctx context.Context, user model.User, groupNames *[]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = func() error {
		userId, err := s.repo.UpdateTx(ctx, tx, user)
		if err != nil {
			return err
		}

		if err = s.membershipService.UpdateTx(ctx, tx, userId, groupNames); err != nil {
			return err
		}

		return tx.Commit()
	}()

	if err != nil {
		tx.Rollback()
	}

	return err
}
