// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int32        `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	Content   string       `json:"content"`
	Score     int32        `json:"score"`
	UserName  string       `json:"user_name"`
	PostID    int32        `json:"post_id"`
}

type CommentVote struct {
	ID        int32        `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	UserID    int32        `json:"user_id"`
	CommentID int32        `json:"comment_id"`
	Upvoted   int32        `json:"upvoted"`
}

type Post struct {
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
}

type PostTag struct {
	ID     int32 `json:"id"`
	Score  int32 `json:"score"`
	TagID  int32 `json:"tag_id"`
	PostID int32 `json:"post_id"`
	UserID int32 `json:"user_id"`
}

type PostTagVote struct {
	ID        int32 `json:"id"`
	UserID    int32 `json:"user_id"`
	PostTagID int32 `json:"post_tag_id"`
	Upvoted   int32 `json:"upvoted"`
}

type PostVote struct {
	ID        int32        `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	Upvoted   int32        `json:"upvoted"`
	UserID    int32        `json:"user_id"`
	PostID    int32        `json:"post_id"`
}

type Tag struct {
	ID        int32  `json:"id"`
	UserLevel int32  `json:"user_level"`
	Name      string `json:"name"`
}

type User struct {
	ID        int32        `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Password  string       `json:"password"`
	Level     int32        `json:"level"`
}
