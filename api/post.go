package api

import (
	"database/sql"
	"errors"
	"fmt"
	"ginbar/api/utils"
	"net/http"
	"path/filepath"
	"strconv"

	"ginbar/api/models"
	"ginbar/mysql/db"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PostsHandler struct defines the Dependencies that will be used
type PostsHandler struct {
	db *sql.DB
}

type postVoteForm struct {
	PostID    int32 `form:"post_id"`
	VoteState int32 `form:"vote_state"`
}

// NewPostsHandler constructor
func NewPostsHandler(db *sql.DB) *PostsHandler {
	return &PostsHandler{db: db}
}

type postForm struct {
	URL string `form:"URL" binding:"required"`
}

var _posts []db.Post = nil

// GetAll retrives all users from the database and returns these users as
// JSON Data
func (server *Server) GetAll(context *gin.Context) {
	// read data from session
	session := sessions.Default(context)

	var userID int32 = 0
	if res := session.Get("userid"); res != nil {
		userID = res.(int32)
	}

	// TODO: check if CORS is needed
	//context.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var posts []models.PostJSON
	var err error
	// IF userID == 0 the client is not logged in so we send Posts without
	// voting information
	// ELSE sending Posts data with voting information
	if userID == 0 {
		posts, err = models.GetPosts(server.store, context)
	} else {
		posts, err = models.GetVotedPosts(server.store, context, userID)
	}

	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	// TODO: remove struct. Not needed anymore but we have to change
	// frontend code too
	type PostsResult struct {
		Posts []models.PostJSON `json:"posts"`
	}

	postsResult := PostsResult{}
	postsResult.Posts = posts

	context.JSON(http.StatusOK, postsResult)

}

// Get retrieves a post from the database
func (server *Server) Get(context *gin.Context) {
	var postID int64
	postID, err := strconv.ParseInt(context.Param("post_id"), 10, 64)
	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	session := sessions.Default(context)
	userID, ok := session.Get("userid").(int32)
	if !ok {
		userID = 0
	}

	if userID == 0 {
		post, err := server.store.GetPost(context, int32(postID))
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
		postParams := db.GetVotedPostParams{
			ID:     int32(postID),
			UserID: int32(userID),
		}
		post, err := server.store.GetVotedPost(context, postParams)
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		p := models.PostJSON{}
		p.PopulateVoted(post)

		commentParams := db.GetVotedCommentsParams{
			UserID: userID,
			PostID: post.ID,
		}
		p.Comments, err = server.store.GetVotedComments(context, commentParams)

		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		tagsParams := db.GetTagsByPostParams{
			UserID: userID,
			PostID: post.ID,
		}

		tags, err := server.store.GetTagsByPost(context, tagsParams)

		var tagsJSON []models.PostTagJSON
		for _, tag := range tags {
			t := models.PostTagJSON{}
			t.Populate(tag)
			tagsJSON = append(tagsJSON, t)
		}

		p.Tags = tagsJSON

		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   p,
		})
	}
}

// CreatePost inserts a user into the database
func (server *Server) CreatePost(context *gin.Context) {
	//var post *db.Post = &db.Post{}
	var err error

	var form postForm

	// context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// Set Status Codes 500 for failed service, if we get to the end
	// completly Status Code 204 will be set
	context.Status(http.StatusInternalServerError)
	err = context.ShouldBind(&form)
	if err != nil {
		context.Error(err)
		return
	}

	fileName, err := utils.ProcessUploadedImage(form.URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// read data from session
	//session := sessions.Default(context)
	//userName, ok := session.Get("user").(string)

	userName, ok := "kejith", true
	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
		return
	}

	parameters := db.CreatePostParams{
		Url:      form.URL,
		Image:    filepath.Base(fileName),
		UserName: userName,
	}

	err = server.store.CreatePost(context, parameters)
	if err != nil {
		context.Error(err)
		return
	}

	// everything worked fine so we send a Status code 204
	// TODO implement Status 201
	context.Status(http.StatusNoContent)
}

// VotePost updates voting information
func (server *Server) VotePost(context *gin.Context) {
	// Read Data from Form
	var form postVoteForm
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
		params := db.UpsertPostVoteParams{
			UserID:  userID,
			PostID:  form.PostID,
			Upvoted: form.VoteState,
		}

		err = server.store.UpsertPostVote(context, params)
	} else {
		params := db.DeletePostVoteParams{
			UserID: userID,
			PostID: form.PostID,
		}

		err = server.store.DeletePostVote(context, params)
	}

	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	context.Status(http.StatusOK)

}

// PostUpdate updates the data of a user in the database
func (server *Server) PostUpdate(context *gin.Context) {
	// TODO: Implement Updates of Posts MAYBE
	// Should Posts be updated?
}
