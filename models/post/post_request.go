package post_model

import (
	"fmt"
	"strings"
)

type PostCreateRequest struct {
	PostHtml string   `json:"postInHtml" validate:"required,min=2,max=500"`
	Tags     []string `json:"tags" validate:"required,min=0,dive,alphanum"`
}

type PostFilters struct {
	Limit     int      `json:"limit" validate:"number,gte=0"`
	Offset    int      `json:"offset" validate:"number,gte=0"`
	Search    string   `json:"search"`
	SearchTag []string `json:"searchTag" validate:"dive,alphanum"`
}

func (pf *PostFilters) BuildQuery() string {
	query := "SELECT p.post_id, p.post_html, p.tags, p.created_at, " +
		"u.user_id, u.name, u.image_url, u.created_at, " +
		"(SELECT COUNT(*) FROM friends f WHERE p.user_id = f.user_id_requester), " +
		"jsonb_agg(jsonb_build_object(" +
		"'comment', c.comment," +
		`'createdAt', to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'),` +
		`'creator', jsonb_build_object('userId', cu.user_id, 'name', cu.name, 'imageUrl', coalesce(cu.image_url, ''), 'friendCount', (SELECT COUNT(*) FROM friends cf WHERE cu.user_id = cf.user_id_requester), 'createdAt', to_char(cu.created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'))` +
		")) AS comments, " +
		"count(*) over() AS total_item " +
		"FROM posts p " +
		"JOIN users u ON p.user_id = u.user_id " +
		"LEFT JOIN comments c ON p.post_id = c.post_id " +
		"LEFT JOIN users cu ON c.user_id = cu.user_id "

	condition := []string{}

	if pf.Search != "" {
		condition = append(condition, fmt.Sprintf("p.post_html LIKE '%%%s%%' ", pf.Search))
	}

	if pf.SearchTag != nil && len(pf.SearchTag) > 0 {
		tagsCondition := []string{}
		for _, tag := range pf.SearchTag {
			tagsCondition = append(tagsCondition, fmt.Sprintf("'%s' = ANY(p.tags)", tag))
		}
		condition = append(condition, strings.Join(tagsCondition, " OR "))
	}

	if len(condition) > 0 {
		query += " WHERE " + strings.Join(condition, " AND ")
	}

	// TODO: Need to adjust this
	query += "GROUP BY p.post_id, u.user_id "
	// query += "GROUP BY p.post_id, u.user_id, c.created_at ORDER BY p.created_at, c.created_at "

	// Add limit and offset
	if pf.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", pf.Limit)
		query += fmt.Sprintf(" OFFSET %d", pf.Offset)
	}

	return query
}
