package models

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"ginbar/mysql/db"
)

// CommentJSON is a struct to map Data from the Database to a reduced JSON object
type CommentJSON struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Content   string    `json:"content"`
	Score     int     `json:"score"`
	Username  string    `json:"username"`
	PostID    int32       `json:"post_id"`
	Upvoted   int8      `json:"upvoted"`
}

// PopulateVoted fills the struct with data from the Database Object
func (c *CommentJSON) PopulateVoted(comment db.GetVotedCommentRow) {
	if comment.DeletedAt.Valid {
		c.DeletedAt = comment.DeletedAt.Time
	}

	c.ID = int(comment.ID)
	c.CreatedAt = comment.CreatedAt
	c.UpdatedAt = comment.UpdatedAt
	c.Content = comment.Content
	c.PostID = comment.PostID
	c.Username = comment.UserName
	c.Score = int(comment.Score)
	c.Upvoted = int8(comment.Upvoted)
}

// PopulateVoteds fills the struct with data from the Database Object
// TODO a bit hacky.
// Because we get 2 diffrent but alike structs we need a second method
// Needs to be fixed
func (c *CommentJSON) PopulateVoteds(comment db.GetVotedCommentsRow) {
	if comment.DeletedAt.Valid {
		c.DeletedAt = comment.DeletedAt.Time
	}

	c.ID = int(comment.ID)
	c.CreatedAt = comment.CreatedAt
	c.UpdatedAt = comment.UpdatedAt
	c.Content = comment.Content
	c.PostID = comment.PostID
	c.Username = comment.UserName
	c.Score = int(comment.Score)
	c.Upvoted = int8(comment.Upvoted.(int32))
}

// PopulateComment fills the struct with data from the Database Object
func (c *CommentJSON) PopulateComment(comment db.Comment) {
	if comment.DeletedAt.Valid {
		c.DeletedAt = comment.DeletedAt.Time
	}

	fmt.Println("PopulateComment")

	c.ID = int(comment.ID)
	c.CreatedAt = comment.CreatedAt
	c.UpdatedAt = comment.UpdatedAt
	c.Content = comment.Content
	c.Username = comment.UserName
	c.PostID = comment.PostID
	c.Score = int(comment.Score)
}

// GetVotedCommentsByPost returns Comments with Vote Information
func GetVotedCommentsByPost(store db.Store, context *gin.Context, params db.GetVotedCommentsParams) ([]CommentJSON, error) {
	comments, err := store.GetVotedComments(context, params)
	if err != nil {
		return nil, err
	}

	// Populate Standardized JSON Array from DatabaseResult
	var commentsJSON []CommentJSON
	for _, comment := range comments {
		var c CommentJSON = CommentJSON{}
		c.PopulateVoteds(comment)
		commentsJSON = append(commentsJSON, c)
	}

	return commentsJSON, nil
}

// GetCommentsByPost returns Comments
func GetCommentsByPost(store db.Store, context *gin.Context, postID int) ([]CommentJSON, error) {
	comments, err := store.GetCommentsByPost(context, int32(postID))
	if err != nil {
		return nil, err
	}

	// Populate Standardized JSON Array from DatabaseResult
	var commentsJSON []CommentJSON
	for _, comment := range comments {
		var c CommentJSON = CommentJSON{}

		c.PopulateComment(comment)
		commentsJSON = append(commentsJSON, c)
	}

	return commentsJSON, nil
}
