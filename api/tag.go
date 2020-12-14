package api

import (
	"errors"
	"fmt"
	"ginbar/api/models"
	"ginbar/mysql/db"
	"net/http"

	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
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
func (server *Server) CreatePostTag(context *gin.Context) {
	// Read Data from Form
	var form createPostTagForm
	err := context.ShouldBind(&form)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	if form.PostID <= 0 {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New("PostID invalid"))
		return
	}

	// Get Session Information
	session := sessions.Default(context)
	fmt.Println(session.Get("userid"))
	userID, ok := session.Get("userid").(int32)

	if !ok {
		context.Status(http.StatusInternalServerError)
		context.Error(errors.New("UserID Type Assertion failed"))
		return
	}

	res, err := server.store.CreateTag(context, form.Name)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	tagID, err := res.LastInsertId()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	params := db.AddTagToPostParams{
		TagID:  int32(tagID),
		PostID: form.PostID,
		UserID: userID,
	}

	res, err = server.store.AddTagToPost(context, params)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	postTagID, err := res.LastInsertId()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	tag := models.PostTagJSON{
		ID:     int32(postTagID),
		Score:  0,
		Name:   form.Name,
		PostID: form.PostID,
		UserID: userID,
	}

	// post mutated we need to recache the post response
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v", tag.PostID)))

	context.JSON(http.StatusOK, tag)
}

// VotePostTag updates voting information
func (server *Server) VotePostTag(context *gin.Context) {
	// Read Data from Form
	var form postTagVoteForm
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
		params := db.UpsertPostTagVoteParams{
			UserID:    userID,
			PostTagID: form.PostTagID,
			Upvoted:   form.VoteState,
		}

		err = server.store.UpsertPostTagVote(context, params)
	} else {
		params := db.DeletePostTagVoteParams{
			UserID:    userID,
			PostTagID: form.PostTagID,
		}

		err = server.store.DeletePostTagVote(context, params)
	}

	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	tag, err := server.store.GetPostTag(context, form.PostTagID)
	if err != nil {
		context.Error(err)
		return
	}

	// post mutated we need to recache the post response
	server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", tag.PostID, userID)))

	context.Status(http.StatusOK)

}
