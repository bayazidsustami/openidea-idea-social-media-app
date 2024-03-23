package friend_model

import (
	"fmt"
	"strings"
)

type FriendRequest struct {
	UserId string `json:"userId" validate:"required"`
}

type FilterFriends struct {
	SortBy   string `json:"sortBy" validate:"oneof=friendCount createdAt ''"`
	OrderBy  string `json:"orderBy" validate:"oneof=asc desc ''"`
	UserOnly bool   `json:"onlyFriend" validate:"boolean"`
	Limit    int    `json:"limit" validate:"number,gte=0"`
	Offset   int    `json:"offset" validate:"number,gte=0"`
	Search   string `json:"search"`
}

func (ff *FilterFriends) BuildQuery(userId string) string {
	query := "SELECT u.user_id, u.name, u.image_url, (SELECT COUNT(*) AS friends_count FROM friends f WHERE f.user_id_requester = u.user_id), u.created_at, count(*) over() AS total_item " +
		"FROM users u "

	condition := []string{}

	if ff.UserOnly {
		query += "JOIN friends f2 ON u.user_id = f2.user_id_requester "
		query += fmt.Sprintf("WHERE f2.user_id_accepter IN (SELECT u2.user_id FROM users u2 WHERE u2.user_id = '%s') ", userId)
	}

	if ff.Search != "" {
		condition = append(condition, fmt.Sprintf("u.name LIKE '%%%s%%' ", ff.Search))
	}

	if len(condition) > 0 {
		if ff.UserOnly {
			query += " AND " + strings.Join(condition, " AND ")
		} else {
			query += " WHERE " + strings.Join(condition, " AND ")
		}
	}

	// Add sorting and ordering
	if ff.SortBy != "" {
		orderBy := "ASC"
		if ff.OrderBy == "desc" {
			orderBy = "DESC"
		}
		mappedSortBy := "friends_count"
		if ff.SortBy == "friendCount" {
			mappedSortBy = "friends_count"
		} else {
			mappedSortBy = "u.created_at"
		}
		query += fmt.Sprintf(" ORDER BY %s %s ", mappedSortBy, orderBy)
	}

	// Add limit and offset
	if ff.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", ff.Limit)
		query += fmt.Sprintf(" OFFSET %d", ff.Offset*ff.Limit)
	}
	return query
}
