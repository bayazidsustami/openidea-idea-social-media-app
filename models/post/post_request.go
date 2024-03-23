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
	SearchTag []string `json:"searchTag"`
}

func (pf *PostFilters) BuildQuery() string {
	query := `
		SELECT p.post_id AS "postId",
			json_build_object(
				'postInHtml', p.post_html,
				'tags', p.tags,
				'createdAt', to_char(p.created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF')
			) AS "post",
			json_agg(
				json_build_object(
					'comment', c.comment,
					'creator', json_build_object(
						'userId', u.user_id,
						'name', u.name,
						'imageUrl', u.image_url,
						'friendCount', (
							SELECT COUNT(*) FROM friends f WHERE f.user_id_requester = u.user_id
						)
					),
					'createdAt', to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF')
				) ORDER BY c.created_at
			) AS "comments",
			json_build_object(
					'userId', u.user_id,
					'name', u.name,
					'imageUrl', u.image_url,
					'friendCount', (
						SELECT COUNT(*) FROM friends f WHERE f.user_id_requester = u.user_id
					)
			) AS "creator"
		FROM 
			posts p
		LEFT JOIN 
			comments c ON p.post_id = c.post_id
		LEFT JOIN 
			users u ON c.user_id = u.user_id
	`

	condition := []string{}

	if pf.Search != "" {
		condition = append(condition, fmt.Sprintf("p.post_html LIKE '%%%s%%' ", pf.Search))
	}

	if pf.SearchTag != nil && len(pf.SearchTag) > 1 {
		tagsCondition := []string{}
		for _, tag := range pf.SearchTag {
			tagsCondition = append(tagsCondition, fmt.Sprintf("'%s' = ANY(p.tags) ", tag))
		}
		condition = append(condition, strings.Join(tagsCondition, " OR "))
	}
	if len(condition) > 0 {
		query += " WHERE " + strings.Join(condition, " AND ")
	}

	query += "GROUP BY p.post_id, p.post_html, p.tags, p.created_at, c.created_at, u.user_id " +
		"ORDER BY p.post_id, c.created_at"

	// Add limit and offset
	if pf.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", pf.Limit)
		query += fmt.Sprintf(" OFFSET %d", pf.Offset)
	}

	return query
}
