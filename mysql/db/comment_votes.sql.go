// Code generated by sqlc. DO NOT EDIT.
// source: comment_votes.sql

package db

import (
	"context"
)

const createCommentVote = `-- name: CreateCommentVote :exec
INSERT INTO comment_votes 
    (user_id, comment_id, upvoted)
VALUES 
    (?, ?, ?)
`

type CreateCommentVoteParams struct {
	UserID    int32 `json:"user_id"`
	CommentID int32 `json:"comment_id"`
	Upvoted   int32 `json:"upvoted"`
}

func (q *Queries) CreateCommentVote(ctx context.Context, arg CreateCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, createCommentVote, arg.UserID, arg.CommentID, arg.Upvoted)
	return err
}

const deleteCommentVote = `-- name: DeleteCommentVote :exec
DELETE FROM comment_votes WHERE user_id = ? AND comment_id = ?
`

type DeleteCommentVoteParams struct {
	UserID    int32 `json:"user_id"`
	CommentID int32 `json:"comment_id"`
}

func (q *Queries) DeleteCommentVote(ctx context.Context, arg DeleteCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, deleteCommentVote, arg.UserID, arg.CommentID)
	return err
}

const updateCommentVote = `-- name: UpdateCommentVote :exec
UPDATE comment_votes
SET upvoted = ?
WHERE id = ?
`

type UpdateCommentVoteParams struct {
	Upvoted int32 `json:"upvoted"`
	ID      int32 `json:"id"`
}

func (q *Queries) UpdateCommentVote(ctx context.Context, arg UpdateCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, updateCommentVote, arg.Upvoted, arg.ID)
	return err
}

const upsertCommentVote = `-- name: UpsertCommentVote :exec
REPLACE INTO 
    comment_votes(user_id, comment_id, upvoted)
VALUE 
    (?, ?, ?)
`

type UpsertCommentVoteParams struct {
	UserID    int32 `json:"user_id"`
	CommentID int32 `json:"comment_id"`
	Upvoted   int32 `json:"upvoted"`
}

func (q *Queries) UpsertCommentVote(ctx context.Context, arg UpsertCommentVoteParams) error {
	_, err := q.db.ExecContext(ctx, upsertCommentVote, arg.UserID, arg.CommentID, arg.Upvoted)
	return err
}
