package fiberapi

import (
	"errors"
	"fmt"

	"ginbar/api/models"
	"ginbar/mysql/db"

	"github.com/gofiber/fiber/v2"
)

type postTagVoteForm struct {
	PostTagID int32 `form:"post_tag_id"`
	VoteState int32 `form:"vote_state"`
}

type createPostTagForm struct {
	Name   string `form:"name"`
	PostID int32  `form:"post_id"`
}

// CreatePostTag creates a Tag
func (server *FiberServer) CreatePostTag(c *fiber.Ctx) error {
	// Read Data from Form
	form := new(createPostTagForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	if form.PostID <= 0 {
		return errors.New("PostID invalid")
	}

	// Get Session Information
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID <= 0 {
		return fmt.Errorf("create post tag: user information could not be loaded from session")
	}

	res, err := server.store.CreateTag(c.Context(), form.Name)
	if err != nil {
		return err
	}

	tagID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	params := db.AddTagToPostParams{
		TagID:  int32(tagID),
		PostID: form.PostID,
		UserID: user.ID,
	}

	res, err = server.store.AddTagToPost(c.Context(), params)
	if err != nil {
		return err
	}

	postTagID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	tag := models.PostTagJSON{
		ID:     int32(postTagID),
		Score:  0,
		Name:   form.Name,
		PostID: form.PostID,
		UserID: user.ID,
	}

	// post mutated we need to recache the post response
	//server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v", tag.PostID)))
	// TODO: caching

	return c.JSON(tag)
}

// VotePostTag updates voting information
func (server *FiberServer) VotePostTag(c *fiber.Ctx) error {
	// Read Data from Form
	form := new(postTagVoteForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	// Get Session Information
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID <= 0 {
		return fmt.Errorf("create post tag: user information could not be loaded from session")
	}

	if form.VoteState != 0 {
		params := db.UpsertPostTagVoteParams{
			UserID:    user.ID,
			PostTagID: form.PostTagID,
			Upvoted:   form.VoteState,
		}

		err = server.store.UpsertPostTagVote(c.Context(), params)
	} else {
		params := db.DeletePostTagVoteParams{
			UserID:    user.ID,
			PostTagID: form.PostTagID,
		}

		err = server.store.DeletePostTagVote(c.Context(), params)
	}

	if err != nil {
		return err
	}

	// tag, err := server.store.GetPostTag(c.Context(), form.PostTagID)
	// if err != nil {
	// 	return err
	// }

	// post mutated we need to recache the post response
	//server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", tag.PostID, userID)))
	// TODO: Caching

	return nil
}
