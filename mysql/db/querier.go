// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
)

type Querier interface {
	AddTagToPost(ctx context.Context, arg AddTagToPostParams) error
	CreateComment(ctx context.Context, arg CreateCommentParams) error
	CreateCommentVote(ctx context.Context, arg CreateCommentVoteParams) error
	CreatePost(ctx context.Context, arg CreatePostParams) error
	CreatePostVote(ctx context.Context, arg CreatePostVoteParams) error
	CreateTag(ctx context.Context, name string) error
	CreateUser(ctx context.Context, arg CreateUserParams) error
	DeleteComment(ctx context.Context, id int32) error
	DeleteCommentVote(ctx context.Context, arg DeleteCommentVoteParams) error
	DeletePost(ctx context.Context, id int32) error
	DeletePostVote(ctx context.Context, arg DeletePostVoteParams) error
	DeleteTag(ctx context.Context, id int32) error
	DeleteTagByName(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, id int32) error
	GetComment(ctx context.Context, id int32) (Comment, error)
	GetComments(ctx context.Context) ([]Comment, error)
	GetCommentsByPost(ctx context.Context, postID int32) ([]Comment, error)
	GetLatestComment(ctx context.Context, userName string) (Comment, error)
	GetPost(ctx context.Context, id int32) (Post, error)
	GetPosts(ctx context.Context) ([]Post, error)
	GetPostsByUser(ctx context.Context, userName string) ([]Post, error)
	GetTag(ctx context.Context, id int32) (Tag, error)
	GetTagByName(ctx context.Context, name string) (Tag, error)
	GetTags(ctx context.Context) ([]Tag, error)
	GetTagsByPost(ctx context.Context, postID int32) ([]GetTagsByPostRow, error)
	GetUser(ctx context.Context, id int32) (GetUserRow, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	GetUsers(ctx context.Context) ([]GetUsersRow, error)
	GetVotedComment(ctx context.Context, arg GetVotedCommentParams) (GetVotedCommentRow, error)
	GetVotedComments(ctx context.Context, arg GetVotedCommentsParams) ([]GetVotedCommentsRow, error)
	GetVotedPost(ctx context.Context, arg GetVotedPostParams) (GetVotedPostRow, error)
	GetVotedPosts(ctx context.Context, userID int32) ([]GetVotedPostsRow, error)
	RemoveTagFromPost(ctx context.Context) error
	UpdateCommentVote(ctx context.Context, arg UpdateCommentVoteParams) error
	UpdatePostVote(ctx context.Context, arg UpdatePostVoteParams) error
	UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) error
	UpsertCommentVote(ctx context.Context, arg UpsertCommentVoteParams) error
	UpsertPostVote(ctx context.Context, arg UpsertPostVoteParams) error
}

var _ Querier = (*Queries)(nil)
