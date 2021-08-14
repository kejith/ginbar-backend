package fiberapi

import (
	"fmt"
	"ginbar/api"
	"ginbar/api/models"
	"ginbar/mysql/db"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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

// GetPost retrieves a singöe Post from the Database
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
