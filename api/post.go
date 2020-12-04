package api

import (
	"errors"
	"fmt"
	"ginbar/api/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"ginbar/api/models"
	"ginbar/mysql/db"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// --------------------
// FORMS
// --------------------

type postVoteForm struct {
	PostID    int32 `form:"post_id"`
	VoteState int32 `form:"vote_state"`
}

type postForm struct {
	URL string `form:"URL" binding:"required"`
}

// --------------------
// Handlers
// --------------------

// CreatePost inserts a user into the database
func (server *Server) CreatePost(context *gin.Context) {
	var err error
	var form postForm

	// Set Status Codes 500 for failed service, if we get to the end
	// completly Status Code 204 will be set
	context.Status(http.StatusInternalServerError)
	err = context.ShouldBind(&form)
	if err != nil {
		context.Error(err)
		return
	}

	fileName, _, err := utils.ProcessUploadedImage(form.URL, server.directories)
	if err != nil {
		fmt.Println(err)
		return
	}

	// read data from session
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
		return
	}

	parameters := db.CreatePostParams{
		Url:               form.URL,
		Filename:          filepath.Base(fileName),
		ThumbnailFilename: filepath.Base(fileName),
		UserName:          userName,
		ContentType:       "image",
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

// UploadPost handles uploads from files and creates a post with them
func (server *Server) UploadPost(context *gin.Context) {
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
		return
	}

	// Limit File Size => 25 << 20 is 25MB
	context.Request.ParseMultipartForm(25 << 20)
	file, handler, err := context.Request.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	mimeType := handler.Header.Get("Content-Type")
	mimeComponents := strings.Split(mimeType, "/")
	fileType, fileFormat := mimeComponents[0], mimeComponents[1]

	switch fileType {
	case "video":
		fileName, thumbnailFilename, err := utils.ProcessUploadedVideo(file, fileFormat, server.directories)

		parameters := db.CreatePostParams{
			Url:               "",
			Filename:          fileName,
			ThumbnailFilename: thumbnailFilename,
			UserName:          userName,
			ContentType:       handler.Header.Get("Content-Type"),
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
}

// GetAll retrives all users from the database and returns these users as
// JSON Data
func (server *Server) GetAll(context *gin.Context) {
	// read data from session
	// session := sessions.Default(context)

	// var userID int32 = 0
	// if res := session.Get("userid"); res != nil {
	// 	userID = res.(int32)
	// }

	// TODO: check if CORS is needed
	//context.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var posts []models.PostJSON
	var err error
	// IF userID == 0 the client is not logged in so we send Posts without
	// voting information
	// ELSE sending Posts data with voting information
	// if userID == 0 {
	posts, err = models.GetPosts(server.store, context)
	// } else {
	// 	posts, err = models.GetVotedPosts(server.store, context, userID)
	// }

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
