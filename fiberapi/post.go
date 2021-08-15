package fiberapi

import (
	"fmt"
	"ginbar/api"
	"ginbar/api/models"
	"ginbar/api/utils"
	"ginbar/mysql/db"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type postForm struct {
	URL string `form:"URL" binding:"required"`
}

type postVoteForm struct {
	PostID    int32 `form:"post_id"`
	VoteState int32 `form:"vote_state"`
}

func (server *FiberServer) GetPosts(c *fiber.Ctx) error {
	var posts *[]models.PostJSON
	var err error

	// Parse HTML Queries
	queries := &models.PostsQueries{}
	if err := c.QueryParser(queries); err != nil {
		return err
	}

	getPostParams := new(models.GetPostsParams)
	getPostParams.SetFromQuery(queries)

	// Retrieve Posts from Store
	posts, err = models.GetPosts(server.store, *getPostParams, c.Context())
	if err != nil {
		return err
	}

	return c.JSON(api.PostsJson{Posts: posts})
}

// GetPost retrieves a single Post from the Database
func (server *FiberServer) GetPost(c *fiber.Ctx) error {
	// parse Parameter
	paramPostID := c.Params("post_id", "0")
	if paramPostID == "0" {
		return fmt.Errorf("post: id can not be 0")
	}
	postID, _ := strconv.ParseInt(paramPostID, 10, 32)

	// user
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	// if user id is 0 we only send public informartion of a post
	if user.ID == 0 {
		params := db.GetPostParams{
			ID:        int32(postID),
			UserLevel: 0,
		}

		post, err := server.store.GetPost(c.Context(), params)
		if err != nil {
			return err
		}

		var p models.PostJSON
		p.PopulatePost(post)

		tagsParams := db.GetTagsByPostParams{
			UserID: user.ID,
			PostID: post.ID,
		}

		tags, _ := server.store.GetTagsByPost(c.Context(), tagsParams)

		var tagsJSON []models.PostTagJSON
		for _, tag := range tags {
			t := models.PostTagJSON{}
			t.Populate(tag)
			tagsJSON = append(tagsJSON, t)
		}

		p.Tags = tagsJSON

		return c.JSON(fiber.Map{
			"status": http.StatusOK,
			"data":   p,
		})

	} else { // if user id is not null sent non-public post information
		postParams := db.GetVotedPostParams{
			ID:        int32(postID),
			UserID:    int32(user.ID),
			UserLevel: int32(user.Level),
		}
		post, err := server.store.GetVotedPost(c.Context(), postParams)
		if err != nil {
			return err
		}

		p := models.PostJSON{}
		p.PopulateVoted(post)

		commentParams := db.GetVotedCommentsParams{
			UserID: user.ID,
			PostID: post.ID,
		}
		p.Comments, err = server.store.GetVotedComments(c.Context(), commentParams)

		if err != nil {
			return err
		}

		tagsParams := db.GetTagsByPostParams{
			UserID: user.ID,
			PostID: post.ID,
		}

		tags, _ := server.store.GetTagsByPost(c.Context(), tagsParams)

		var tagsJSON []models.PostTagJSON
		for _, tag := range tags {
			t := models.PostTagJSON{}
			t.Populate(tag)
			tagsJSON = append(tagsJSON, t)
		}

		p.Tags = tagsJSON

		return c.JSON(fiber.Map{
			"status": http.StatusOK,
			"data":   p,
		})
	}
}

func (server *FiberServer) createPostFromFile(c *fiber.Ctx, url, inputFile string) (*db.Post, []db.GetPossibleDuplicatePostsRow) {
	var err error

	user, err := server.GetUserFromSession(c)
	if err != nil {
		panic(err)
	}

	if user.Name == "" {
		panic(fmt.Errorf("createPostFromFile: could not get Userdata[User Name]"))
	}

	parameters := db.CreatePostParams{Url: url, UserName: user.Name}

	mimeType := mime.TypeByExtension(filepath.Ext(inputFile))
	mimeComponents := strings.Split(mimeType, "/")
	fileType := mimeComponents[0]

	var duplicatePosts []db.GetPossibleDuplicatePostsRow
	processResult := &utils.ImageProcessResult{}
	switch fileType {
	case "image":
		processResult, err = utils.ProcessImage(inputFile, server.directories)
		if err != nil {
			panic(err)
		}

		hash := processResult.PerceptionHash

		duplicatePosts, err = models.GetDuplicatePosts(server.store, c.Context(), hash)
		if err != nil {
			fmt.Println(fmt.Errorf("couldnt get duplicate posts: %w", err))
		}

		if len(duplicatePosts) > 0 {
			return nil, duplicatePosts
		}

		parameters.PHash0 = hash.GetHash()[0]
		parameters.PHash1 = hash.GetHash()[1]
		parameters.PHash2 = hash.GetHash()[2]
		parameters.PHash3 = hash.GetHash()[3]

		parameters.ContentType = "image"

	case "video":
		processResult.Filename, processResult.ThumbnailFilename, err = utils.ProcessVideo(inputFile, fileType, server.directories)

		if err != nil {
			panic(err)
		}

		parameters.ContentType = mimeType
	}

	parameters.Filename = filepath.Base(processResult.Filename)
	parameters.ThumbnailFilename = filepath.Base(processResult.ThumbnailFilename)

	res, err := server.store.CreatePost(c.Context(), parameters)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	postID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	post, _ := server.store.GetPost(c.Context(), db.GetPostParams{
		ID:        int32(postID),
		UserLevel: user.Level,
	})

	return &post, duplicatePosts
}

func (server *FiberServer) CreatePost(c *fiber.Ctx) error {
	postForm := new(postForm)

	if err := c.BodyParser(postForm); err != nil {
		return err
	}

	filePath, err := utils.DownloadFile(postForm.URL, server.directories.Tmp)
	if err != nil {
		panic(fmt.Errorf("could not download the File from the provided URL: %w", err))
	}

	post, reposts := server.createPostFromFile(c, postForm.URL, filePath)
	if len(reposts) > 0 {
		return c.JSON(fiber.Map{
			"status": "possibleDuplicatesFound",
			"posts":  reposts,
		})
	}

	// we mutated posts so we need to recache the getPosts response
	//server.postsResponseCache.Flush()
	// TODO: Caching

	// everything worked fine so we send a Status code 200
	if post != nil { // if we didnt create a Post
		return c.JSON(fiber.Map{
			"status": "postCreated",
			"posts":  []db.Post{*post},
		})
	} else {
		c.JSON(fiber.Map{
			"status": "postCreated",
			"posts":  []db.Post{},
		})
	}

	return nil
}

func (server *FiberServer) UploadPost(c *fiber.Ctx) error {
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID <= 0 {
		return fmt.Errorf("upload post: user data could not be loaded from session")
	}

	// parse uploaded file from form
	file, err := c.FormFile("file")
	if err != nil {
		return fmt.Errorf("failed to parse uploaded file: %w", err)
	}

	fp := filepath.Base(file.Filename)
	tmpFilePath := filepath.Join(server.directories.Tmp, fp)
	if err := c.SaveFile(file, tmpFilePath); err != nil {
		panic(fmt.Errorf("could not save uploaded File: %w", err))
	}

	post, reposts := server.createPostFromFile(c, "", tmpFilePath)
	if len(reposts) > 0 {
		return c.JSON(fiber.Map{
			"status": "possibleDuplicatesFound",
			"posts":  reposts,
		})
	}

	// we mutated posts so we need to recache the getPosts response
	//server.postsResponseCache.Delete(fmt.Sprintf("/api/post/#%v", userID))
	// TODO: Caching

	if post != nil { // if we didnt create a post
		return c.JSON(fiber.Map{
			"status": "postCreated",
			"posts":  []db.Post{*post},
		})
	} else {
		return c.JSON(fiber.Map{
			"status": "postCreated",
			"posts":  []db.Post{},
		})
	}
}

func (server *FiberServer) VotePost(c *fiber.Ctx) error {
	voteForm := new(postVoteForm)
	if err := c.BodyParser(voteForm); err != nil {
		return err
	}

	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID == 0 {
		return fmt.Errorf("votePost: can not load user information from sesstion")
	}

	if voteForm.VoteState != 0 {
		params := db.UpsertPostVoteParams{
			UserID:  user.ID,
			PostID:  voteForm.PostID,
			Upvoted: voteForm.VoteState,
		}

		err = server.store.UpsertPostVote(c.Context(), params)
	} else {
		params := db.DeletePostVoteParams{
			UserID: user.ID,
			PostID: voteForm.PostID,
		}

		err = server.store.DeletePostVote(c.Context(), params)
	}

	if err != nil {
		return err
	}

	// post mutated we need to recache the post response
	//server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", form.PostID, userID)))
	//context.Status(http.StatusOK)

	return c.SendStatus(fiber.StatusOK)
}
