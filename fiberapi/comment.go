package fiberapi

import (
	"fmt"
	"strconv"

	"ginbar/api/models"
	"ginbar/mysql/db"

	"github.com/gofiber/fiber/v2"
)

// --------------------
// FORMS
// --------------------

type commentWriteForm struct {
	Content string `form:"content"`
	PostID  int32  `form:"post_id"`
}

type commentVoteForm struct {
	CommentID int32 `form:"comment_id"`
	VoteState int32 `form:"vote_state"`
}

// --------------------
// Handlers
// --------------------

// GetComments retrives all comments from the database
func (server *FiberServer) GetComments(c *fiber.Ctx) error {

	paramPostID := c.Params("post_id", "0")
	postID, _ := strconv.ParseInt(paramPostID, 10, 32)
	if postID <= 0 {
		return fmt.Errorf("post id invalid")
	}

	// user
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	params := db.GetVotedCommentsParams{
		PostID: int32(postID),
		UserID: int32(user.ID),
	}

	// We serve Comments with voting information when
	comments, err := server.store.GetVotedComments(c.Context(), params)

	if err != nil {
		return err
	}

	return c.JSON(comments)
}

// CreateComment inserts a user into the database
func (server *FiberServer) CreateComment(c *fiber.Ctx) error {
	// Read Data from Form
	form := new(commentWriteForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	if form.PostID <= 0 {
		return fmt.Errorf("post id invalid")
	}

	// Get Session Information
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID <= 0 {
		return fmt.Errorf("upload post: user data could not be loaded from session [userid]")
	}

	if user.Name == "" {
		return fmt.Errorf("upload post: user data could not be loaded from session [username]")
	}

	comment, err := models.NewComment(form.Content, form.PostID, user.Name)
	if err != nil {
		return err
	}

	err = comment.Save(&server.store, c.Context())
	if err != nil {
		return err
	}

	// post mutated we need to recache the post response
	// err = server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", form.PostID, userID)))
	// if err != nil {
	// 	context.Error(err)
	// }
	// TODO: Caching

	return c.JSON(comment)
}

// VoteComment upserts vote information into the database
func (server *FiberServer) VoteComment(c *fiber.Ctx) error {
	// Read Data from Form
	form := new(commentVoteForm)
	if err := c.BodyParser(form); err != nil {
		return err
	}

	if form.CommentID <= 0 {
		return fmt.Errorf("comment id invalid")
	}

	// Get Session Information
	user, err := server.GetUserFromSession(c)
	if err != nil {
		return err
	}

	if user.ID <= 0 {
		return fmt.Errorf("upload post: user data could not be loaded from session [userid]")
	}

	if user.Name == "" {
		return fmt.Errorf("upload post: user data could not be loaded from session [username]")
	}

	if form.VoteState != 0 {
		params := db.UpsertCommentVoteParams{
			UserID:    user.ID,
			CommentID: form.CommentID,
			Upvoted:   form.VoteState,
		}

		err = server.store.UpsertCommentVote(c.Context(), params)
	} else {
		params := db.DeleteCommentVoteParams{
			UserID:    user.ID,
			CommentID: form.CommentID,
		}

		err = server.store.DeleteCommentVote(c.Context(), params)
	}

	if err != nil {
		return err
	}

	// comment, err := server.store.GetComment(c.Context(), form.CommentID)
	// if err != nil {
	// 	return err
	// }

	// post mutated we need to recache the post response
	// server.postsResponseCache.Delete(cache.CreateKey(fmt.Sprintf("/api/post/%v#%v", comment.PostID, userID)))
	// TODO: Caching

	return nil
}

/*
// GetComment retrieves a post from the database
func (server *Server) GetComment(context *gin.Context) {
	commentID, err := strconv.ParseInt(context.Param("id"), 10, 32)

	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	session := sessions.Default(context)
	userID, ok := session.Get("userid").(int64)
	if !ok {
		userID = 0
	}

	if userID == 0 {
		post, err := server.store.GetComment(context, int(commentID))
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
		post, err := server.store.GetVotedPost(context, postID, userID)
		if err != nil {
			context.Error(err)
			context.Status(http.StatusInternalServerError)
			return
		}

		p := models.PostJSON{}
		p.PopulateVoted(post)

		context.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   p,
		})
	}
}

// Create inserts a user into the database
func (server *Server) Create(context *gin.Context) {
	var comment *entity.Comment = &entity.Comment{}
	var form commentWriteForm

	err := context.ShouldBind(&form)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	session := sessions.Default(context)
	userName := session.Get("user")

	comment.Content = form.Content
	comment.UserName = userName.(string)
	comment.PostID = form.PostID

	if len(comment.Content) < 4 {
		context.Status(http.StatusUnprocessableEntity)
		fmt.Println("Length of the Comment is too short to create")
		return
	}

	comment, err = c.application.Create(comment)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		context.Error(err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":         http.StatusOK,
		"createdComment": comment,
	})
}

// CommentUpdate updates the data of a user in the database
func (server *Server) CommentUpdate(context *gin.Context) {
	var post *entity.Comment

	post = loadComment(c, context)

	// TODO update Comment data when relevant Comment Data is added
	_ = post

	context.Status(http.StatusAccepted)
}

// CommentShow renders a view for a post
func (server *Server) CommentShow(context *gin.Context) {
	post := loadComment(c, context)
	_ = post

}

func loadComment(server *Server, context *gin.Context) (post *entity.Comment) {
	commentID, err := strconv.ParseInt(context.Param("id"), 10, 32)
	if err != nil {
		context.Error(err)
		context.Status(http.StatusInternalServerError)

		return nil
	}

	fmt.Println(commentID)

	post, err = c.application.Get(commentID)
	if err != nil {
		context.Error(err)
		switch err.(type) {
		case *application.ParametersNotValidError:
			context.Status(http.StatusUnprocessableEntity)
		default:
			context.Status(http.StatusInternalServerError)
		}

		return nil
	}

	return post
}
*/
