package api

import (
	"errors"
	"fmt"
	"ginbar/api/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

// CreateMultiplePosts inserts multiple posts into the database
func (server *Server) CreateMultiplePosts(context *gin.Context) {
	var form postForm

	// Set Status Codes 500 for failed service, if we get to the end
	// completly Status Code 204 will be set
	context.Status(http.StatusInternalServerError)
	err := context.ShouldBind(&form)
	if err != nil {
		context.Error(err)
		return
	}

	urls := strings.Split(form.URL, ",")

	if len(urls) < 1 {
		panic(errors.New("No URL transmitted"))
	}

	var posts = make([](*db.Post), len(urls))
	var total = len(urls)
	for i, url := range urls {
		var start = time.Now()
		fmt.Print(fmt.Sprintf("[%v/%v] %s", i, total, url))
		post := server.createPostFromURL(context, url)
		if post != nil {
			posts[i] = post
		}
		fmt.Print(fmt.Sprintf(" %v\n", time.Since(start)))
	}

	// we mutated posts so we need to recache the getPosts response
	session := sessions.Default(context)
	userID, ok := session.Get("userid").(int32)
	if !ok {
		userID = 0
	}
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/#%v", userID)))

	// everything worked fine so we send a Status code 204
	// TODO implement Status 201
	context.JSON(http.StatusOK, gin.H{
		"status": "postCreated",
		"posts":  posts,
	})

}

func (server *Server) createPostFromURL(context *gin.Context, url string) *db.Post {
	var err error

	response, fileType, fileFormat, err := utils.LoadFileFromURL(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer response.Body.Close()

	// read data from session
	session := sessions.Default(context)
	userName, ok := session.Get("user").(string)
	userLevel, ok := session.Get("userlevel").(int32)
	if !ok {
		userLevel = 0
	}

	if !ok {
		fmt.Println(errors.New(" PostHandler.Create => Type Assertion failed on session['user']"))
		return nil
	}

	parameters := db.CreatePostParams{Url: url, UserName: userName}
	var processResult utils.ImageProcessResult

	switch fileType {
	case "image":
		processResult, err = utils.ProcessImageFromURL(response, fileFormat, server.directories)
		params := db.GetPossibleDuplicatePostsParams{
			Column1: processResult.PerceptionHash.GetHash()[0],
			Column2: processResult.PerceptionHash.GetHash()[1],
			Column3: processResult.PerceptionHash.GetHash()[2],
			Column4: processResult.PerceptionHash.GetHash()[3],
		}

		duplicatePosts, err := server.store.GetPossibleDuplicatePosts(context, params)

		if err != nil {
			fmt.Println(err)
			return nil
		}

		if len(duplicatePosts) > 0 {
			context.JSON(http.StatusOK, gin.H{
				"status": "possibleDuplicatesFound",
				"posts":  duplicatePosts,
			})

			return nil
		}

		if err != nil {
			context.Error(err)
			return nil
		}

		parameters.PHash0 = processResult.PerceptionHash.GetHash()[0]
		parameters.PHash1 = processResult.PerceptionHash.GetHash()[1]
		parameters.PHash2 = processResult.PerceptionHash.GetHash()[2]
		parameters.PHash3 = processResult.PerceptionHash.GetHash()[3]

		parameters.ContentType = "image"

		break
	case "video":
		processResult.Filename, processResult.ThumbnailFilename, err = utils.ProcessVideoFromURL(response, fileFormat, server.directories)

		if err != nil {
			context.Error(err)
			return nil
		}

		parameters.ContentType = fmt.Sprintf("%s/%s", fileType, fileFormat)
		break
	}

	parameters.Filename = filepath.Base(processResult.Filename)
	parameters.ThumbnailFilename = filepath.Base(processResult.ThumbnailFilename)

	res, err := server.store.CreatePost(context, parameters)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	postID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	post, err := server.store.GetPost(context, db.GetPostParams{
		ID:        int32(postID),
		UserLevel: userLevel,
	})

	return &post
}

// CreatePost inserts a post into the database
func (server *Server) CreatePost(context *gin.Context) {
	var form postForm

	// Set Status Codes 500 for failed service, if we get to the end
	// completly Status Code 204 will be set
	context.Status(http.StatusInternalServerError)
	err := context.ShouldBind(&form)
	if err != nil {
		context.Error(err)
		return
	}

	post := server.createPostFromURL(context, form.URL)

	// we mutated posts so we need to recache the getPosts response
	session := sessions.Default(context)
	userID, ok := session.Get("userid").(int32)
	if !ok {
		userID = 0
	}
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/#%v", userID)))

	// everything worked fine so we send a Status code 204
	// TODO implement Status 201
	context.JSON(http.StatusOK, gin.H{
		"status": "postCreated",
		"posts":  []db.Post{*post},
	})

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

	parameters := db.CreatePostParams{
		UserName: userName,
	}

	switch fileType {
	case "video":
		parameters.Filename, parameters.ThumbnailFilename, err = utils.ProcessUploadedVideo(&file, fileFormat, server.directories)
		parameters.ContentType = mimeType
		// everything worked fine so we send a Status code 204
		// TODO implement Status 201
		context.Status(http.StatusNoContent)
		break
	case "image":
		processResult, err := utils.ProcessImageFromMultipart(&file, fileFormat, server.directories)

		if err != nil {
			panic(err)
		}

		parameters.Filename = filepath.Base(processResult.Filename)
		parameters.ThumbnailFilename = filepath.Base(processResult.ThumbnailFilename)
		parameters.ContentType = "image"

		parameters.PHash0 = processResult.PerceptionHash.GetHash()[0]
		parameters.PHash1 = processResult.PerceptionHash.GetHash()[1]
		parameters.PHash2 = processResult.PerceptionHash.GetHash()[2]
		parameters.PHash3 = processResult.PerceptionHash.GetHash()[3]

		break
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
