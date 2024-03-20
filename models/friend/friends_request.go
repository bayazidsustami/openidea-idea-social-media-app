package friend_model

import (
	"fmt"
	"strings"
)

type FriendRequest struct {
	UserId int `json:"userId" validate:"required"`
}

type FilterFriends struct {
	SortBy   string `json:"sortBy" validate:"oneof=friendCount createdAt ''"`
	OrderBy  string `json:"orderBy" validate:"oneof=asc dsc ''"`
	UserOnly bool   `json:"userOnly" validate:"boolean"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Search   string `json:"search"`
}

func (ff *FilterFriends) BuildQuery(userId int) string {
	query := "SELECT f.user_id_accepter, u.name, u.image_url, COUNT(f2.user_id_accepter) AS friends_count, f.created_at, count(*) over() AS total_item " +
		"FROM users u " +
		"JOIN friends f ON u.user_id = f.user_id_requester " +
		"LEFT JOIN friends f2 ON f.user_id_accepter = f2.user_id_accepter "

	condition := []string{}

	if ff.UserOnly {
		condition = append(condition, fmt.Sprintf(" f.user_id_requester = %d ", userId))
	}

	if ff.Search != "" {
		condition = append(condition, fmt.Sprintf("u.name LIKE '%%%s%%' ", ff.Search))
	}

	if len(condition) > 0 {
		query += " WHERE " + strings.Join(condition, " AND ")
	}

	query += "GROUP BY f.user_id_accepter, u.name, u.image_url, f.created_at, u.user_id "

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
			mappedSortBy = "f.created_at"
		}
		query += fmt.Sprintf(" ORDER BY %s %s ", mappedSortBy, orderBy)
	}

	// Add limit and offset
	if ff.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", ff.Limit)
		query += fmt.Sprintf(" OFFSET %d", ff.Offset)
	}
	return query
}
