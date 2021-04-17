package membership

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/yassinekhaliqui/go-rest-service/internal/model"
)

type Repository interface {
	GetGroupsForUser(ctx context.Context, userId uint64) (*[]model.Group, error)
	GetUsersForGroup(ctx context.Context, groupId uint64) (*[]model.User, error)
	InsertTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error
	UpdateTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error
	UpdateGroupMembership(ctx context.Context, groupName string, userIds *[]string) error
}

type repository struct {
	db *sql.DB
}

// Creates a new membership repo instance
func NewRepository(db *sql.DB) Repository {
	return repository{db}
}

// Gets the groups that the user belongs to
func (r repository) GetGroupsForUser(ctx context.Context, userId uint64) (*[]model.Group, error) {
	rows, err := r.db.QueryContext(ctx, "call get_user_membership(?)", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []model.Group

	for rows.Next() {
		var group model.Group
		if err := rows.Scan(&group.Id, &group.Name); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return &groups, nil
}

// Gets the users that are inside of a group
func (r repository) GetUsersForGroup(ctx context.Context, groupId uint64) (*[]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "call get_group_membership(?)", groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.UserId); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

// Inserts a link between a user and an array of groups
// Done in a transaction
func (r repository) InsertTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error {
	groupNamesStr := toDelimitedString(groupNames, ",")
	if groupNamesStr == "" {
		return nil
	}
	_, err := tx.QueryContext(ctx, "call ins_membership(?, ?)", userId, groupNamesStr)
	return err
}

// Removes existing user - group rows and inserts new ones
// Done in a transaction
func (r repository) UpdateTx(ctx context.Context, tx *sql.Tx, userId uint64, groupNames *[]string) error {
	groupNamesStr := toDelimitedString(groupNames, ",")
	if groupNamesStr == "" {
		return nil
	}
	_, err := tx.QueryContext(ctx, "call upd_membership(?, ?)", userId, groupNamesStr)
	return err
}

// Removes existing users of a group, and inserts new users
func (r repository) UpdateGroupMembership(ctx context.Context, groupName string, userIds *[]string) error {
	userIdsStr := toDelimitedString(userIds, ",")
	if userIdsStr == "" {
		return nil
	}
	_, err := r.db.QueryContext(ctx, "call upd_group_membership(?, ?)", groupName, userIdsStr)
	return err
}

// Converts a list to a delimited separated string
func toDelimitedString(strs *[]string, delimiter string) string {
	if strs ==  nil || len(*strs) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, groupName := range *strs {
		fmt.Fprintf(&sb, "\"%s\",", groupName)
	}
	
	return sb.String()[:sb.Len()-1]
}
