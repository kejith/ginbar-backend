// Code generated by sqlc. DO NOT EDIT.
// source: post_tag_votes.sql

package db

import (
	"context"
)

const createPostTagVote = `-- name: CreatePostTagVote :exec
INSERT INTO post_tag_votes 
    (user_id, post_tag_id, upvoted)
VALUES 
    (?, ?, ?)
`

type CreatePostTagVoteParams struct {
	UserID    int32 `json:"user_id"`
	PostTagID int32 `json:"post_tag_id"`
	Upvoted   int32 `json:"upvoted"`
}

func (q *Queries) CreatePostTagVote(ctx context.Context, arg CreatePostTagVoteParams) error {
	_, err := q.db.ExecContext(ctx, createPostTagVote, arg.UserID, arg.PostTagID, arg.Upvoted)
	return err
}

const deletePostTagVote = `-- name: DeletePostTagVote :exec
DELETE FROM post_tag_votes WHERE user_id = ? AND post_tag_id = ?
`

type DeletePostTagVoteParams struct {
	UserID    int32 `json:"user_id"`
	PostTagID int32 `json:"post_tag_id"`
}

func (q *Queries) DeletePostTagVote(ctx context.Context, arg DeletePostTagVoteParams) error {
	_, err := q.db.ExecContext(ctx, deletePostTagVote, arg.UserID, arg.PostTagID)
	return err
}

const upsertPostTagVote = `-- name: UpsertPostTagVote :exec
REPLACE INTO 
    post_tag_votes(user_id, post_tag_id, upvoted)
VALUE 
    (?, ?, ?)
`

type UpsertPostTagVoteParams struct {
	UserID    int32 `json:"user_id"`
	PostTagID int32 `json:"post_tag_id"`
	Upvoted   int32 `json:"upvoted"`
}

func (q *Queries) UpsertPostTagVote(ctx context.Context, arg UpsertPostTagVoteParams) error {
	_, err := q.db.ExecContext(ctx, upsertPostTagVote, arg.UserID, arg.PostTagID, arg.Upvoted)
	return err
}
