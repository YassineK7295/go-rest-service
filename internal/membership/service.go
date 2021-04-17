package membership

import (
	"context"
	"database/sql"

	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

type Service interface {
	GetGroupsForUser(ctx context.Context, userId uint64) (*[]model.Group, error)
	GetUsersForGroup(ctx context.Context, groupId uint64) (*[]model.User, error)
	InsertTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error
	UpdateTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error
	UpdateGroupMembership(ctx context.Context, groupName string, userIds *[]string) error
}

type service struct {
	repo Repository
}

// Creates a new membership service instance
func NewService(db *sql.DB) Service {
	return service{NewRepository(db)}
}

// Gets groups for a user
func (s service) GetGroupsForUser(ctx context.Context, userId uint64) (*[]model.Group, error) {
	return s.repo.GetGroupsForUser(ctx, userId)
}

// Gets users for a group
func (s service) GetUsersForGroup(ctx context.Context, groupId uint64) (*[]model.User, error) {
	return s.repo.GetUsersForGroup(ctx, groupId)
}

// Inserts user to groups linkage as part of a transaction
func (s service) InsertTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error {
	return s.repo.InsertTx(ctx, tx, userId, groupNames)
}

// Updates a user to groups linkage as part of a transaction
func (s service) UpdateTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error {
	return s.repo.UpdateTx(ctx, tx, userId, groupNames)
}

// Updates group membership
func (s service) UpdateGroupMembership(ctx context.Context, groupName string, userIds *[]string) error {
	return s.repo.UpdateGroupMembership(ctx, groupName, userIds)
}
