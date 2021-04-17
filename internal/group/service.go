package group

import (
	"context"
	"database/sql"

	"github.com/yassinekhaliqui/go-rest-service/internal/membership"
	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

type Service interface {
	GetWithUsers(ctx context.Context, groupName string) (model.Group, *[]model.User, error)
	Insert(ctx context.Context, group model.Group) (uint64, error)
	Delete(ctx context.Context, groupName string) error
	UpdateGroupMembership(ctx context.Context, groupName string, userIds *[]string) error
}

type service struct {
	repo              Repository
	membershipService membership.Service
}

// Creates a new group service instance
func NewService(db *sql.DB) Service {
	return service{NewRepository(db), membership.NewService(db)}
}

// Gets the group and the linked users
func (s service) GetWithUsers(ctx context.Context, groupName string) (model.Group, *[]model.User, error) {
	group, err := s.repo.Get(ctx, groupName)
	if err != nil {
		return model.Group{}, nil, err
	}

	users, err := s.membershipService.GetUsersForGroup(ctx, group.Id)
	if err != nil {
		return model.Group{}, nil, err
	}

	return group, users, nil
}

// Inserts a new group
func (s service) Insert(ctx context.Context, group model.Group) (uint64, error) {
	return s.repo.Insert(ctx, group)
}

// Deletes the group
func (s service) Delete(ctx context.Context, groupName string) error {
	return s.repo.Delete(ctx, groupName)
}

// Updates the membership of the group
func (s service) UpdateGroupMembership(ctx context.Context, groupName string, userIds *[]string) error {
	return s.membershipService.UpdateGroupMembership(ctx, groupName, userIds)
}
