package api

import (
	"database/sql"
	"errors"
	"fmt"
	"ginbar/api/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

type FileUploadForm struct {
	fileData *multipart.FileHeader `form:"file_data" binding:"required"`
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

// UploadPost handles uploads from files and creates a post with them
func (server *Server) UploadPost(context *gin.Context) {
	file, handler, err := context.Request.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	mimeHeader := strings.Split(handler.Header.Get("Content-Type"), "/")
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", mimeHeader)

	cwd, err := os.Getwd()
	if err != nil {
		context.Error(err)
		return
	}

	imageDir := filepath.Join(cwd, "public", "images")
	thumbnailDir := filepath.Join(imageDir, "thumbnails")
	videoDir := filepath.Join(cwd, "public", "videos")
	tmpDir := filepath.Join(cwd, "tmp")
	_ = thumbnailDir

	fmt.Println(thumbnailDir)

	fileType := mimeHeader[0]
	fileFormat := mimeHeader[1]
	time := time.Now().UnixNano()
	fileName := fmt.Sprintf("%v.%s", time, fileFormat)
	thumbnailFileName := fmt.Sprintf("%v.%s", time, "png")

	switch fileType {
	case "video":
		filePath := filepath.Join(videoDir, fileName)

		localFile, err := os.Create(filePath)
		if err != nil {
			context.Error(err)
			return
		}
		defer file.Close()

		_, err = io.Copy(localFile, file)
		if err != nil {

			context.Error(err)
			return
		}

		tmpThumbnailFilePath := filepath.Join(tmpDir, thumbnailFileName)
		commandArgs := fmt.Sprintf("-i %s -ss 00:00:01.000 -vframes 1 %s -hide_banner -loglevel panic", filePath, tmpThumbnailFilePath)
		cmd := exec.Command("ffmpeg", strings.Split(commandArgs, " ")...)
		//cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

		fmt.Println(filepath.Join(thumbnailDir, thumbnailFileName))
		err = utils.CreateThumbnailFromFile(tmpThumbnailFilePath, filepath.Join(thumbnailDir, thumbnailFileName))

		if err != nil {
			fmt.Println("Thumbnail creation from Temporary file failed")
			fmt.Println(err)
			log.Fatalf("Thumbnail creation from Temporary file failed with %s\n", err)
		}

		session := sessions.Default(context)
		userName, ok := session.Get("user").(string)

		if !ok {
			context.Status(http.StatusInternalServerError)
			context.Error(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
			return
		}

		parameters := db.CreatePostParams{
			Url:         "",
			Filename:    fileName,
			UserName:    userName,
			ContentType: handler.Header.Get("Content-Type"),
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

	// Create file
	//dst, err := os.Create()

	/*
		// The file cannot be received.
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "File could not be received",
			})
			return
		}
	*/
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
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)

	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
		return
	}

	parameters := db.CreatePostParams{
		Url:         form.URL,
		Filename:    filepath.Base(fileName),
		UserName:    userName,
		ContentType: "image",
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
