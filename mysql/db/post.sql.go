// Code generated by sqlc. DO NOT EDIT.
// source: post.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createPost = `-- name: CreatePost :execresult
INSERT INTO posts 
    (
		url, 
		filename, 
		thumbnail_filename, 
		user_name, 
		content_type, 
		p_hash_0, 
		p_hash_1, 
		p_hash_2, 
		p_hash_3, 
		uploaded_filename
	)
VALUES 
    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type CreatePostParams struct {
	Url               string `json:"url"`
	Filename          string `json:"filename"`
	ThumbnailFilename string `json:"thumbnail_filename"`
	UserName          string `json:"user_name"`
	ContentType       string `json:"content_type"`
	PHash0            uint64 `json:"p_hash_0"`
	PHash1            uint64 `json:"p_hash_1"`
	PHash2            uint64 `json:"p_hash_2"`
	PHash3            uint64 `json:"p_hash_3"`
	UploadedFilename  string `json:"uploaded_filename"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createPost,
		arg.Url,
		arg.Filename,
		arg.ThumbnailFilename,
		arg.UserName,
		arg.ContentType,
		arg.PHash0,
		arg.PHash1,
		arg.PHash2,
		arg.PHash3,
		arg.UploadedFilename,
	)
}

const deletePost = `-- name: DeletePost :exec
UPDATE posts 
SET deleted_at = NOW() 
WHERE
	id = ?
`

func (q *Queries) DeletePost(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getAllPosts = `-- name: GetAllPosts :many
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts
`

func (q *Queries) GetAllPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getAllPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
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

const getImagePosts = `-- name: GetImagePosts :many
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts 
WHERE
	content_type = "image"
ORDER BY
	posts.id DESC
`

func (q *Queries) GetImagePosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getImagePosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
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

const getNewerPosts = `-- name: GetNewerPosts :many
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts 
WHERE
	deleted_at IS NULL AND
	posts.id >= ? AND
	posts.user_level <= ?
ORDER BY
	posts.id
LIMIT ?
`

type GetNewerPostsParams struct {
	ID        int32 `json:"id"`
	UserLevel int32 `json:"user_level"`
	Limit     int32 `json:"limit"`
}

func (q *Queries) GetNewerPosts(ctx context.Context, arg GetNewerPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getNewerPosts, arg.ID, arg.UserLevel, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
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

const getOlderPosts = `-- name: GetOlderPosts :many
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts 
WHERE
	deleted_at IS NULL AND
	posts.id <= ? AND
	posts.user_level <= ?
ORDER BY
	posts.id DESC
LIMIT ?
`

type GetOlderPostsParams struct {
	ID        int32 `json:"id"`
	UserLevel int32 `json:"user_level"`
	Limit     int32 `json:"limit"`
}

func (q *Queries) GetOlderPosts(ctx context.Context, arg GetOlderPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getOlderPosts, arg.ID, arg.UserLevel, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
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

const getPossibleDuplicatePosts = `-- name: GetPossibleDuplicatePosts :many
SELECT 
	posts.id, posts.created_at, posts.updated_at, posts.deleted_at, posts.url, posts.uploaded_filename, posts.filename, posts.thumbnail_filename, posts.content_type, posts.score, posts.user_level, posts.p_hash_0, posts.p_hash_1, posts.p_hash_2, posts.p_hash_3, posts.user_name, 
    (
		bit_count(? ^ p_hash_0) +
        bit_count(? ^ p_hash_1) +
        bit_count(? ^ p_hash_2) +
        bit_count(? ^ p_hash_3) 
        
	) as hamming_distance
    from posts
    having hamming_distance < 50
    order by hamming_distance desc
`

type GetPossibleDuplicatePostsParams struct {
	Column1 interface{} `json:"column_1"`
	Column2 interface{} `json:"column_2"`
	Column3 interface{} `json:"column_3"`
	Column4 interface{} `json:"column_4"`
}

type GetPossibleDuplicatePostsRow struct {
	ID                int32        `json:"id"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at"`
	Url               string       `json:"url"`
	UploadedFilename  string       `json:"uploaded_filename"`
	Filename          string       `json:"filename"`
	ThumbnailFilename string       `json:"thumbnail_filename"`
	ContentType       string       `json:"content_type"`
	Score             int32        `json:"score"`
	UserLevel         int32        `json:"user_level"`
	PHash0            uint64       `json:"p_hash_0"`
	PHash1            uint64       `json:"p_hash_1"`
	PHash2            uint64       `json:"p_hash_2"`
	PHash3            uint64       `json:"p_hash_3"`
	UserName          string       `json:"user_name"`
	HammingDistance   int32        `json:"hamming_distance"`
}

func (q *Queries) GetPossibleDuplicatePosts(ctx context.Context, arg GetPossibleDuplicatePostsParams) ([]GetPossibleDuplicatePostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPossibleDuplicatePosts,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPossibleDuplicatePostsRow{}
	for rows.Next() {
		var i GetPossibleDuplicatePostsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
			&i.HammingDistance,
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

const getPost = `-- name: GetPost :one
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts 
WHERE
	posts.id = ? AND 
	deleted_at IS NULL AND
	posts.user_level <= ?
`

type GetPostParams struct {
	ID        int32 `json:"id"`
	UserLevel int32 `json:"user_level"`
}

func (q *Queries) GetPost(ctx context.Context, arg GetPostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, arg.ID, arg.UserLevel)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Url,
		&i.UploadedFilename,
		&i.Filename,
		&i.ThumbnailFilename,
		&i.ContentType,
		&i.Score,
		&i.UserLevel,
		&i.PHash0,
		&i.PHash1,
		&i.PHash2,
		&i.PHash3,
		&i.UserName,
	)
	return i, err
}

const getPosts = `-- name: GetPosts :many
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts 
WHERE
	deleted_at IS NULL AND
	posts.user_level <= ?
ORDER BY
	posts.id DESC
LIMIT 50
`

func (q *Queries) GetPosts(ctx context.Context, userLevel int32) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts, userLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
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

const getPostsByUser = `-- name: GetPostsByUser :many
SELECT
	id, created_at, updated_at, deleted_at, url, uploaded_filename, filename, thumbnail_filename, content_type, score, user_level, p_hash_0, p_hash_1, p_hash_2, p_hash_3, user_name 
FROM
	posts 
WHERE
	user_name = ? AND 
	deleted_at IS NULL AND
	posts.user_level <= ?
ORDER BY posts.id DESC
`

type GetPostsByUserParams struct {
	UserName  string `json:"user_name"`
	UserLevel int32  `json:"user_level"`
}

func (q *Queries) GetPostsByUser(ctx context.Context, arg GetPostsByUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsByUser, arg.UserName, arg.UserLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
			&i.UserName,
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
	p.id, p.created_at, p.updated_at, p.deleted_at, p.url, p.uploaded_filename, p.filename, p.thumbnail_filename, p.content_type, p.score, p.user_level, p.p_hash_0, p.p_hash_1, p.p_hash_2, p.p_hash_3, p.user_name, 
	IFNULL(pv.upvoted, 0) as upvoted 
FROM
	posts p
	LEFT JOIN post_votes AS pv ON pv.post_id = p.id 
	AND pv.user_id = ? AND
	p.user_level <= ?
WHERE
	p.deleted_at IS NULL AND
	p.id = ?
`

type GetVotedPostParams struct {
	UserID    int32 `json:"user_id"`
	UserLevel int32 `json:"user_level"`
	ID        int32 `json:"id"`
}

type GetVotedPostRow struct {
	ID                int32        `json:"id"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at"`
	Url               string       `json:"url"`
	UploadedFilename  string       `json:"uploaded_filename"`
	Filename          string       `json:"filename"`
	ThumbnailFilename string       `json:"thumbnail_filename"`
	ContentType       string       `json:"content_type"`
	Score             int32        `json:"score"`
	UserLevel         int32        `json:"user_level"`
	PHash0            uint64       `json:"p_hash_0"`
	PHash1            uint64       `json:"p_hash_1"`
	PHash2            uint64       `json:"p_hash_2"`
	PHash3            uint64       `json:"p_hash_3"`
	UserName          string       `json:"user_name"`
	Upvoted           interface{}  `json:"upvoted"`
}

func (q *Queries) GetVotedPost(ctx context.Context, arg GetVotedPostParams) (GetVotedPostRow, error) {
	row := q.db.QueryRowContext(ctx, getVotedPost, arg.UserID, arg.UserLevel, arg.ID)
	var i GetVotedPostRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.Url,
		&i.UploadedFilename,
		&i.Filename,
		&i.ThumbnailFilename,
		&i.ContentType,
		&i.Score,
		&i.UserLevel,
		&i.PHash0,
		&i.PHash1,
		&i.PHash2,
		&i.PHash3,
		&i.UserName,
		&i.Upvoted,
	)
	return i, err
}

const getVotedPosts = `-- name: GetVotedPosts :many
SELECT
	p.id, p.created_at, p.updated_at, p.deleted_at, p.url, p.uploaded_filename, p.filename, p.thumbnail_filename, p.content_type, p.score, p.user_level, p.p_hash_0, p.p_hash_1, p.p_hash_2, p.p_hash_3, p.user_name,
	IFNULL(pv.upvoted, 0) as upvoted 
FROM
	posts p
	LEFT JOIN ( SELECT id, created_at, updated_at, deleted_at, upvoted, user_id, post_id FROM post_votes WHERE user_id = ? ) AS pv ON pv.post_id = p.id 
WHERE
	p.deleted_at IS NULL AND
	p.user_level <= ?
ORDER BY p.id DESC
`

type GetVotedPostsParams struct {
	UserID    int32 `json:"user_id"`
	UserLevel int32 `json:"user_level"`
}

type GetVotedPostsRow struct {
	ID                int32        `json:"id"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at"`
	Url               string       `json:"url"`
	UploadedFilename  string       `json:"uploaded_filename"`
	Filename          string       `json:"filename"`
	ThumbnailFilename string       `json:"thumbnail_filename"`
	ContentType       string       `json:"content_type"`
	Score             int32        `json:"score"`
	UserLevel         int32        `json:"user_level"`
	PHash0            uint64       `json:"p_hash_0"`
	PHash1            uint64       `json:"p_hash_1"`
	PHash2            uint64       `json:"p_hash_2"`
	PHash3            uint64       `json:"p_hash_3"`
	UserName          string       `json:"user_name"`
	Upvoted           interface{}  `json:"upvoted"`
}

func (q *Queries) GetVotedPosts(ctx context.Context, arg GetVotedPostsParams) ([]GetVotedPostsRow, error) {
	rows, err := q.db.QueryContext(ctx, getVotedPosts, arg.UserID, arg.UserLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetVotedPostsRow{}
	for rows.Next() {
		var i GetVotedPostsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.Url,
			&i.UploadedFilename,
			&i.Filename,
			&i.ThumbnailFilename,
			&i.ContentType,
			&i.Score,
			&i.UserLevel,
			&i.PHash0,
			&i.PHash1,
			&i.PHash2,
			&i.PHash3,
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

const updatePostFiles = `-- name: UpdatePostFiles :exec
UPDATE
	posts
SET
	filename = ?,
	thumbnail_filename = ?
WHERE
	id = ?
`

type UpdatePostFilesParams struct {
	Filename          string `json:"filename"`
	ThumbnailFilename string `json:"thumbnail_filename"`
	ID                int32  `json:"id"`
}

func (q *Queries) UpdatePostFiles(ctx context.Context, arg UpdatePostFilesParams) error {
	_, err := q.db.ExecContext(ctx, updatePostFiles, arg.Filename, arg.ThumbnailFilename, arg.ID)
	return err
}

const updatePostHashes = `-- name: UpdatePostHashes :exec
UPDATE
	posts
SET
	p_hash_0 = ?,
	p_hash_1 = ?,
	p_hash_2 = ?,
	p_hash_3 = ?
WHERE
	id = ?
`

type UpdatePostHashesParams struct {
	PHash0 uint64 `json:"p_hash_0"`
	PHash1 uint64 `json:"p_hash_1"`
	PHash2 uint64 `json:"p_hash_2"`
	PHash3 uint64 `json:"p_hash_3"`
	ID     int32  `json:"id"`
}

func (q *Queries) UpdatePostHashes(ctx context.Context, arg UpdatePostHashesParams) error {
	_, err := q.db.ExecContext(ctx, updatePostHashes,
		arg.PHash0,
		arg.PHash1,
		arg.PHash2,
		arg.PHash3,
		arg.ID,
	)
	return err
}
