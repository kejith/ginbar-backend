// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AddTagToPost(ctx context.Context, arg AddTagToPostParams) (sql.Result, error)
	CreateComment(ctx context.Context, arg CreateCommentParams) error
	CreateCommentVote(ctx context.Context, arg CreateCommentVoteParams) error
	CreatePost(ctx context.Context, arg CreatePostParams) (sql.Result, error)
	CreatePostTagVote(ctx context.Context, arg CreatePostTagVoteParams) error
	CreatePostVote(ctx context.Context, arg CreatePostVoteParams) error
	CreateTag(ctx context.Context, name string) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) error
	DeleteComment(ctx context.Context, id int32) error
	DeleteCommentVote(ctx context.Context, arg DeleteCommentVoteParams) error
	DeletePost(ctx context.Context, id int32) error
	DeletePostTagVote(ctx context.Context, arg DeletePostTagVoteParams) error
	DeletePostVote(ctx context.Context, arg DeletePostVoteParams) error
	DeleteTag(ctx context.Context, id int32) error
	DeleteTagByName(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, id int32) error
	GetAllPosts(ctx context.Context) ([]Post, error)
	GetComment(ctx context.Context, id int32) (Comment, error)
	GetComments(ctx context.Context) ([]Comment, error)
	GetCommentsByPost(ctx context.Context, postID int32) ([]Comment, error)
	GetImagePosts(ctx context.Context) ([]Post, error)
	GetLatestComment(ctx context.Context, userName string) (Comment, error)
	GetNewerPosts(ctx context.Context, arg GetNewerPostsParams) ([]Post, error)
	GetOlderPosts(ctx context.Context, arg GetOlderPostsParams) ([]Post, error)
	GetPossibleDuplicatePosts(ctx context.Context, arg GetPossibleDuplicatePostsParams) ([]GetPossibleDuplicatePostsRow, error)
	GetPost(ctx context.Context, arg GetPostParams) (Post, error)
	GetPostTag(ctx context.Context, id int32) (PostTag, error)
	GetPosts(ctx context.Context, userLevel int32) ([]Post, error)
	GetPostsByUser(ctx context.Context, arg GetPostsByUserParams) ([]Post, error)
	GetTag(ctx context.Context, id int32) (Tag, error)
	GetTagByName(ctx context.Context, name string) (Tag, error)
	GetTags(ctx context.Context) ([]Tag, error)
	GetTagsByPost(ctx context.Context, arg GetTagsByPostParams) ([]GetTagsByPostRow, error)
	GetUser(ctx context.Context, id int32) (GetUserRow, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	GetUsers(ctx context.Context) ([]GetUsersRow, error)
	GetVotedComment(ctx context.Context, arg GetVotedCommentParams) (GetVotedCommentRow, error)
	GetVotedComments(ctx context.Context, arg GetVotedCommentsParams) ([]GetVotedCommentsRow, error)
	GetVotedPost(ctx context.Context, arg GetVotedPostParams) (GetVotedPostRow, error)
	GetVotedPosts(ctx context.Context, arg GetVotedPostsParams) ([]GetVotedPostsRow, error)
	RemoveTagFromPost(ctx context.Context) error
	UpdateCommentVote(ctx context.Context, arg UpdateCommentVoteParams) error
	UpdatePostFiles(ctx context.Context, arg UpdatePostFilesParams) error
	UpdatePostHashes(ctx context.Context, arg UpdatePostHashesParams) error
	UpdatePostVote(ctx context.Context, arg UpdatePostVoteParams) error
	UpdateUserEmail(ctx context.Context, arg UpdateUserEmailParams) error
	UpsertCommentVote(ctx context.Context, arg UpsertCommentVoteParams) error
	UpsertPostTagVote(ctx context.Context, arg UpsertPostTagVoteParams) error
	UpsertPostVote(ctx context.Context, arg UpsertPostVoteParams) error
}

var _ Querier = (*Queries)(nil)
