package models

import (
	"time"

	"ginbar/mysql/db"

	"github.com/gin-gonic/gin"
)

// PostJSON is a struct to map Data from the Database to a reduced JSON object
type PostJSON struct {
	ID          int64                    `json:"id"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	DeletedAt   time.Time                `json:"deleted_at"`
	URL         string                   `json:"url"`
	FileName    string                   `json:"filename"`
	ContentType string                   `json:"content_type"`
	UserName    string                   `json:"user"`
	Upvoted     int8                     `json:"upvoted"`
	Score       int32                    `json:"score"`
	Comments    []db.GetVotedCommentsRow `json:"comments"`
	Tags        []PostTagJSON            `json:"tags"`
}

// PopulateVoteds fills the struct with data from the Database Object
func (p *PostJSON) PopulateVoteds(post db.GetVotedPostsRow) {
	if post.DeletedAt.Valid {
		p.DeletedAt = post.DeletedAt.Time
	}

	p.CreatedAt = post.CreatedAt
	p.UpdatedAt = post.UpdatedAt
	p.ID = int64(post.ID)
	p.URL = post.Url
	p.FileName = post.Filename
	p.ContentType = post.ContentType
	p.UserName = post.UserName
	p.Score = post.Score
	p.Upvoted = int8(post.Upvoted.(int64))
}

// PopulateVoted fills the struct with data from the Database Object
func (p *PostJSON) PopulateVoted(post db.GetVotedPostRow) {
	if post.DeletedAt.Valid {
		p.DeletedAt = post.DeletedAt.Time
	}

	p.CreatedAt = post.CreatedAt
	p.UpdatedAt = post.UpdatedAt
	p.ID = int64(post.ID)
	p.URL = post.Url
	p.FileName = post.Filename
	p.ContentType = post.ContentType
	p.UserName = post.UserName
	p.Upvoted = int8(post.Upvoted.(int64))
	p.Score = post.Score
}

// PopulatePost fills the struct with data from the Database Object
func (p *PostJSON) PopulatePost(post db.Post) {
	if post.DeletedAt.Valid {
		p.DeletedAt = post.DeletedAt.Time
	}

	p.CreatedAt = post.CreatedAt
	p.UpdatedAt = post.UpdatedAt
	p.ID = int64(post.ID)
	p.URL = post.Url
	p.FileName = post.Filename
	p.ContentType = post.ContentType
	p.UserName = post.UserName
	p.Score = post.Score
}

// GetVotedPosts returns Posts with voting information
func GetVotedPosts(store db.Store, context *gin.Context, userID int32) ([]PostJSON, error) {
	posts, err := store.GetVotedPosts(context, userID)
	if err != nil {
		return nil, err
	}

	var postsJSON []PostJSON
	for _, post := range posts {
		var p PostJSON = PostJSON{}
		p.PopulateVoteds(post)
		postsJSON = append(postsJSON, p)
	}

	return postsJSON, nil
}

// GetPosts returns Posts with voting information
func GetPosts(store db.Store, context *gin.Context) ([]PostJSON, error) {
	posts, err := store.GetPosts(context)
	if err != nil {
		return nil, err
	}

	var postsJSON []PostJSON
	for _, post := range posts {
		var p PostJSON = PostJSON{}
		p.PopulatePost(post)
		postsJSON = append(postsJSON, p)
	}

	return postsJSON, nil
}
