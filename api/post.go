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

	"github.com/gin-contrib/cache"
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

	response, fileType, fileFormat, err := utils.LoadFileFromURL(form.URL)
	if err != nil {
		context.Error(err)
		return
	}
	defer response.Body.Close()

	var filePath string
	var contentType string
	var thumbnailFilePath string
	switch fileType {
	case "image":
		filePath, thumbnailFilePath, err = utils.ProcessImageFromURL(response, fileFormat, server.directories)
		if err != nil {
			context.Error(err)
			return
		}

		contentType = "image"

		break
	case "video":
		filePath, thumbnailFilePath, err = utils.ProcessVideoFromURL(response, fileFormat, server.directories)
		if err != nil {
			context.Error(err)
			return
		}

		contentType = fmt.Sprintf("%s/%s", fileType, fileFormat)
		break
	}

	// read data from session
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	userLevel, ok := session.Get("userlevel").(int32)
	if !ok {
		userLevel = 0
	}

	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
		return
	}

	parameters := db.CreatePostParams{
		Url:               form.URL,
		Filename:          filepath.Base(filePath),
		ThumbnailFilename: filepath.Base(thumbnailFilePath),
		UserName:          userName,
		ContentType:       contentType,
	}

	res, err := server.store.CreatePost(context, parameters)
	if err != nil {
		panic(err)
	}

	postID, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	post, err := server.store.GetPost(context, db.GetPostParams{
		ID:        int32(postID),
		UserLevel: userLevel,
	})

	userID, ok := session.Get("userid").(int32)
	if !ok {
		userID = 0
	}

	// we mutated posts so we need to recache the getPosts response
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/#%v", userID)))

	// everything worked fine so we send a Status code 204
	// TODO implement Status 201
	context.JSON(http.StatusOK, post)

}

// UploadPost handles uploads from files and creates a post with them
func (server *Server) UploadPost(context *gin.Context) {
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)
	userLevel, ok := session.Get("userlevel").(int32)
	if !ok {
		userLevel = 0
	}

	if !ok {

		panic(
			errors.New(
				" PostHandler.Create => Type Assertion failed on session['user']"))
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

	var fileName string
	var thumbnailFilename string
	var contentType string
	switch fileType {
	case "video":
		fileName, thumbnailFilename, err = utils.ProcessUploadedVideo(&file, fileFormat, server.directories)
		contentType = mimeType
		// everything worked fine so we send a Status code 204
		// TODO implement Status 201
		context.Status(http.StatusNoContent)
		break
	case "image":
		fileName, thumbnailFilename, err = utils.ProcessImageFromMultipart(&file, fileFormat, server.directories)
		contentType = "image"
		break
	}

	parameters := db.CreatePostParams{
		Url:               "",
		Filename:          fileName,
		ThumbnailFilename: thumbnailFilename,
		UserName:          userName,
		ContentType:       contentType,
	}

	res, err := server.store.CreatePost(context, parameters)
	if err != nil {
		panic(err)
	}

	postID, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	post, err := server.store.GetPost(context, db.GetPostParams{
		ID:        int32(postID),
		UserLevel: userLevel,
	})

	userID, ok := session.Get("userid").(int32)
	if !ok {
		userID = 0
	}

	// we mutated posts so we need to recache the getPosts response
	server.postsResponseCache.Delete(fmt.Sprintf("/api/post/#%v", userID))
	context.JSON(http.StatusOK, post)
}

// GetAll retrives all users from the database and returns these users as
// JSON Data
func (server *Server) GetAll(context *gin.Context) {

	var posts *[]models.PostJSON
	var err error

	posts, err = models.GetPosts(server.store, context)

	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"posts": *posts})

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

	userLevel, ok := session.Get("userlevel").(int32)
	if !ok {
		userLevel = 0
	}

	if userID == 0 {
		params := db.GetPostParams{
			ID:        int32(postID),
			UserLevel: 0,
		}

		post, err := server.store.GetPost(context, params)
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
			ID:        int32(postID),
			UserID:    int32(userID),
			UserLevel: int32(userLevel),
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

	// post mutated we need to recache the post response
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", form.PostID, userID)))
	context.Status(http.StatusOK)

}
