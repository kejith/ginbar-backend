package fiberapi

import (
	"ginbar/api"
	"ginbar/api/models"

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
	// Parse HTML Queries
	queries := &models.PostQueries{}
	if err := c.QueryParser(queries); err != nil {
		return err
	}

	return nil
}
