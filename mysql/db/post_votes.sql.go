// Code generated by sqlc. DO NOT EDIT.
// source: post_votes.sql

package db

import (
	"context"
)

const createPostVote = `-- name: CreatePostVote :exec
INSERT INTO post_votes 
    (user_id, post_id, upvoted)
VALUES 
    (?, ?, ?)
`

type CreatePostVoteParams struct {
	UserID  int32 `json:"user_id"`
	PostID  int32 `json:"post_id"`
	Upvoted int32 `json:"upvoted"`
}

func (q *Queries) CreatePostVote(ctx context.Context, arg CreatePostVoteParams) error {
	_, err := q.db.ExecContext(ctx, createPostVote, arg.UserID, arg.PostID, arg.Upvoted)
	return err
}

const deletePostVote = `-- name: DeletePostVote :exec
UPDATE post_votes
SET deleted_at = NOW()
WHERE id =
`

func (q *Queries) DeletePostVote(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deletePostVote, id)
	return err
}

const updatePostVote = `-- name: UpdatePostVote :exec
UPDATE post_votes
SET upvoted = ?
WHERE id = ?
`

type UpdatePostVoteParams struct {
	Upvoted int32 `json:"upvoted"`
	ID      int32 `json:"id"`
}

func (q *Queries) UpdatePostVote(ctx context.Context, arg UpdatePostVoteParams) error {
	_, err := q.db.ExecContext(ctx, updatePostVote, arg.Upvoted, arg.ID)
	return err
}
