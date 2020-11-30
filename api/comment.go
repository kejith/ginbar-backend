package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ginbar/mysql/db"
	"ginbar/api/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type commentWriteForm struct {
	Content string `form:"content"`
	PostID  int32  `form:"post_id"`
}

type commentVoteForm struct {
	CommentID int32 `form:"comment_id"`
	VoteState int32 `form:"vote_state"` 
}

// we probably dont need this
// we populate the posts with comments in the post request

// GetComments retrives all comments from the database
func (server *Server) GetComments(context *gin.Context) {
	postID, err := strconv.ParseInt(context.Param("post_id"), 10, 32)
	if err != nil {
		context.Error(errors.New("Post ID is not valid"))
	}
	// read data from session
	session := sessions.Default(context)

	var userID uint = 0
	if res := session.Get("userid"); res != nil {
		userID = res.(uint)
	}

	//var comments []models.CommentJSON
	if userID != 0 {
		// 	// We serve Comments without Vote information when we do not have a User
		// 	// ID
		// 	// TODO: If we dont have a User ID we dont want to serve comments at all
		// 	fmt.Println("Comments without UserID")

		// 	dbComments, err := server.store.GetCommentsByPost(context, int32(postID))
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		context.Error(err)
		// 		context.Status(http.StatusInternalServerError)
		// 		return
		// 	}

		// 	fmt.Println("dbComments length: ", len(dbComments))

		// 	comments := make([]models.CommentJSON, len(dbComments))
		// 	for i, comment := range dbComments {

		// 		var deletedAt time.Time
		// 		if comment.DeletedAt.Valid {
		// 			deletedAt = comment.DeletedAt.Time
		// 		}

		// 		comments[i] = models.CommentJSON{
		// 			ID:        int(comment.ID),
		// 			CreatedAt: comment.CreatedAt,
		// 			UpdatedAt: comment.UpdatedAt,
		// 			DeletedAt: deletedAt,
		// 			Content:   comment.Content,
		// 			Username:  comment.UserName,
		// 			Ups:       int(comment.Ups),
		// 			Downs:     int(comment.Downs),
		// 		}

		// 	}

		// 	// comments, err = models.GetCommentsByPost(server.store, context, int(postID))

		// } else {

		params := db.GetVotedCommentsParams{
			PostID: int32(postID),
			UserID: int32(userID),
		}
		// We serve Comments with voting information when we have user information
		comments, err := server.store.GetVotedComments(context, params)
		//comments, err = models.GetVotedCommentsByPost(server.store, context, params)

		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.JSON(http.StatusOK, comments)
		return
	}

	context.Status(http.StatusInternalServerError)
}

// CreateComment inserts a user into the database
func (server *Server) CreateComment(context *gin.Context) {
	// Read Data from Form
	var form commentWriteForm
	err := context.ShouldBind(&form)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	if form.PostID <= 0 {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New("PostID invalid"))
		return
	}

	// Get Session Information
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New("Username Type Assertion failed"))
		return
	}

	params := db.CreateCommentParams{
		Content:  form.Content,
		UserName: userName,
		PostID:   form.PostID,
	}

	if len(params.Content) < 4 {
		context.Status(http.StatusUnprocessableEntity)
		fmt.Println("Length of the Comment is too short to create")
		return
	}

	err = server.store.CreateComment(context, params)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	comment, err := server.store.GetLatestComment(context, userName)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}
	
	commentJSON := models.CommentJSON{}
	commentJSON.PopulateComment(comment)
	context.JSON(http.StatusOK, commentJSON)
}

// VoteComment upserts vote information into the database
func (server *Server) VoteComment(context *gin.Context) {
	// Read Data from Form
	var form commentVoteForm
	err := context.ShouldBind(&form)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	// read userID from session
	session := sessions.Default(context)
	userID, ok := session.Get("userid").(int32)
	if !ok {
		context.Status(http.StatusInternalServerError)
		return
	}

	fmt.Println(userID)

	if form.VoteState != 0 {
		params := db.UpsertCommentVoteParams{
			UserID: userID,
			CommentID: form.CommentID,
			Upvoted: form.VoteState,
		}

		err = server.store.UpsertCommentVote(context, params);
	} else {
		params := db.DeleteCommentVoteParams{
			UserID: userID,
			CommentID: form.CommentID,
		}

		err = server.store.DeleteCommentVote(context, params);
	}

	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	context.Status(http.StatusOK)
}
/*
// GetComment retrieves a post from the database
func (server *Server) GetComment(context *gin.Context) {
	commentID, err := strconv.ParseInt(context.Param("id"), 10, 32)

	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	session := sessions.Default(context)
	userID, ok := session.Get("userid").(int64)
	if !ok {
		userID = 0
	}

	if userID == 0 {
		post, err := server.store.GetComment(context, int(commentID))
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}
		var p models.PostJSON
		p.PopulatePost(post)

		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   p,
		})
	} else {
		post, err := server.store.GetVotedPost(context, postID, userID)
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		p := models.PostJSON{}
		p.PopulateVoted(post)

		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   p,
		})
	}
}

// Create inserts a user into the database
func (server *Server) Create(context *gin.Context) {
	var comment *entity.Comment = &entity.Comment{}
	var form commentWriteForm

	err := context.ShouldBind(&form)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	session := sessions.Default(context)
	userName := session.Get("user")

	comment.Content = form.Content
	comment.UserName = userName.(string)
	comment.PostID = form.PostID

	if len(comment.Content) < 4 {
		context.Status(http.StatusUnprocessableEntity)
		fmt.Println("Length of the Comment is too short to create")
		return
	}

	comment, err = c.application.Create(comment)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":         http.StatusOK,
		"createdComment": comment,
	})
}

// CommentUpdate updates the data of a user in the database
func (server *Server) CommentUpdate(context *gin.Context) {
	var post *entity.Comment

	post = loadComment(c, context)

	// TODO update Comment data when relevant Comment Data is added
	_ = post

	context.Status(http.StatusAccepted)
}

// CommentShow renders a view for a post
func (server *Server) CommentShow(context *gin.Context) {
	post := loadComment(c, context)
	_ = post

}

func loadComment(server *Server, context *gin.Context) (post *entity.Comment) {
	commentID, err := strconv.ParseInt(context.Param("id"), 10, 32)
	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)

		return nil
	}

	fmt.Println(commentID)

	post, err = c.application.Get(commentID)
	if err != nil {
		context.Error(err)
		switch err.(type) {
		case *application.ParametersNotValidError:
			context.Status(http.StatusUnprocessableEntity)
		default:
			context.Status(http.StatusInternalServerError)
		}

		return nil
	}

	return post
}
*/
