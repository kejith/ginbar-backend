package db

/*
import (
	"context"
)

const getVotedPosts = `-- name: GetVotedPosts :many
SELECT
	p.*,
	IFNULL(pv.upvoted, 0) as upvoted
FROM
	posts p
	LEFT JOIN post_votes AS pv ON pv.post_id = p.id
	AND pv.user_id = ?
WHERE
	p.deleted_at IS NULL
`

// GetVotedPosts Returns Posts with user associated vote
func (q *Queries) GetVotedPosts(ctx context.Context, userID uint) ([]VotedPost, error) {
	rows, err := q.db.QueryContext(ctx, getVotedPosts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []VotedPost{}
	for rows.Next() {
		var i VotedPost
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.Image,
			&i.UserName,
			&i.Upvoted,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVotedPost = `-- name: GetVotedPost :one
SELECT
	p.*,
	IFNULL(pv.upvoted, 0) as upvoted
FROM
	posts p
	LEFT JOIN post_votes AS pv ON pv.post_id = p.id
	AND pv.user_id = ?
WHERE
	p.deleted_at IS NULL AND
	p.id = ?;
`

// GetVotedPost Returns Posts with user associated vote
func (q *Queries) GetVotedPost(ctx context.Context, postID int64, userID int64) (VotedPost, error) {
	row := q.db.QueryRowContext(ctx, getVotedPost, postID, userID)
	var i VotedPost
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Url,
		&i.Image,
		&i.UserName,
		&i.Upvoted,
	)
	return i, err
}
*/
