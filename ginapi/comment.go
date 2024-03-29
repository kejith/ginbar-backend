package ginapi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ginbar/models"
	"ginbar/mysql/db"

	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// --------------------
// FORMS
// --------------------

type commentWriteForm struct {
	Content string `form:"content"`
	PostID  int32  `form:"post_id"`
}

type commentVoteForm struct {
	CommentID int32 `form:"comment_id"`
	VoteState int32 `form:"vote_state"`
}

// --------------------
// Handlers
// --------------------

// GetComments retrives all comments from the database
func (server *Server) GetComments(context *gin.Context) {

	postID, err := strconv.ParseInt(context.Param("post_id"), 10, 32)
	if err != nil {
		panic(errors.New("post ID is not valid"))
	}
	// read data from session
	session := sessions.Default(context)
	userID := session.Get("userid").(uint)

	params := db.GetVotedCommentsParams{
		PostID: int32(postID),
		UserID: int32(userID),
	}

	// We serve Comments with voting information when
	comments, err := server.store.GetVotedComments(context, params)

	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, comments)
}

// CreateComment inserts a user into the database
func (server *Server) CreateComment(context *gin.Context) {
	// Read Data from Form
	var form commentWriteForm

	err := context.ShouldBind(&form)
	if err != nil {
		panic(err)
	}

	// Get Session Information
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	if !ok {
		panic(errors.New("username Type Assertion failed"))
	}

	comment, err := models.NewComment(form.Content, form.PostID, userName)
	if err != nil {
		panic(err)
	}

	err = comment.Save(&server.store, context)
	if err != nil {
		panic(err)
	}

	userID, ok := session.Get("userid").(int32)
	if !ok {
		userID = 0
	}

	// post mutated we need to recache the post response
	err = server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", form.PostID, userID)))
	if err != nil {
		context.Error(err)
	}

	context.JSON(http.StatusOK, comment)
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

	if form.VoteState != 0 {
		params := db.UpsertCommentVoteParams{
			UserID:    userID,
			CommentID: form.CommentID,
			Upvoted:   form.VoteState,
		}

		err = server.store.UpsertCommentVote(context, params)
	} else {
		params := db.DeleteCommentVoteParams{
			UserID:    userID,
			CommentID: form.CommentID,
		}

		err = server.store.DeleteCommentVote(context, params)
	}

	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	comment, err := server.store.GetComment(context, form.CommentID)

	if err != nil {
		context.Error(err)
		return
	}

	// post mutated we need to recache the post response
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", comment.PostID, userID)))

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
