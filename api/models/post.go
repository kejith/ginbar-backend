package models

import (
	"strconv"
	"time"

	"ginbar/mysql/db"

	"github.com/corona10/goimagehash"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// PostJSON is a struct to map Data from the Database to a reduced JSON object
type PostJSON struct {
	ID                int64                    `json:"id"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	DeletedAt         time.Time                `json:"deleted_at"`
	URL               string                   `json:"url"`
	FileName          string                   `json:"filename"`
	ThumbnailFilename string                   `json:"thumbnail_filename"`
	ContentType       string                   `json:"content_type"`
	UserName          string                   `json:"user"`
	Upvoted           int8                     `json:"upvoted"`
	Score             int32                    `json:"score"`
	Comments          []db.GetVotedCommentsRow `json:"comments"`
	Tags              []PostTagJSON            `json:"tags"`
}

// PopulateVoteds fills the struct with data from the Database Object
func (p *PostJSON) PopulateVoteds(post db.GetVotedPostsRow) {
	if post.DeletedAt.Valid {
		p.DeletedAt = post.DeletedAt.Time
	}

	p.CreatedAt = post.CreatedAt
	p.UpdatedAt = post.UpdatedAt
	p.ID = int64(post.ID)
	p.URL = post.Url
	p.FileName = post.Filename
	p.ThumbnailFilename = post.ThumbnailFilename
	p.ContentType = post.ContentType
	p.UserName = post.UserName
	p.Score = post.Score
	p.Upvoted = int8(post.Upvoted.(int64))
}

// PopulateVoted fills the struct with data from the Database Object
func (p *PostJSON) PopulateVoted(post db.GetVotedPostRow) {
	if post.DeletedAt.Valid {
		p.DeletedAt = post.DeletedAt.Time
	}

	p.CreatedAt = post.CreatedAt
	p.UpdatedAt = post.UpdatedAt
	p.ID = int64(post.ID)
	p.URL = post.Url
	p.FileName = post.Filename
	p.ThumbnailFilename = post.ThumbnailFilename
	p.ContentType = post.ContentType
	p.UserName = post.UserName
	p.Upvoted = int8(post.Upvoted.(int64))
	p.Score = post.Score
}

// PopulatePost fills the struct with data from the Database Object
func (p *PostJSON) PopulatePost(post db.Post) {
	if post.DeletedAt.Valid {
		p.DeletedAt = post.DeletedAt.Time
	}

	p.CreatedAt = post.CreatedAt
	p.UpdatedAt = post.UpdatedAt
	p.ID = int64(post.ID)
	p.URL = post.Url
	p.FileName = post.Filename
	p.ThumbnailFilename = post.ThumbnailFilename
	p.ContentType = post.ContentType
	p.UserName = post.UserName
	p.Score = post.Score
}

// GetVotedPosts returns Posts with voting information
func GetVotedPosts(store db.Store, context *gin.Context, params db.GetVotedPostsParams) ([]PostJSON, error) {
	posts, err := store.GetVotedPosts(context, params)
	if err != nil {
		return nil, err
	}

	var postsJSON []PostJSON
	for _, post := range posts {
		var p PostJSON = PostJSON{}
		p.PopulateVoteds(post)
		postsJSON = append(postsJSON, p)
	}

	return postsJSON, nil
}

// GetPosts returns Posts with voting information
func GetPosts(store db.Store, context *gin.Context) (*[]PostJSON, error) {
	lowestIDString, ok := context.GetQuery("lowestID")
	var lowestID int32 = 0
	if ok {
		i, err := strconv.ParseInt(lowestIDString, 10, 32)
		if err == nil {
			lowestID = int32(i)
		}
	}

	highestIDString, ok := context.GetQuery("highestID")
	var highestID int32 = 0
	if ok {
		i, err := strconv.ParseInt(highestIDString, 10, 32)
		if err == nil {
			highestID = int32(i)
		}
	}

	postsPerRowStr, ok := context.GetQuery("postsPerRow")
	var postsPerRow int32 = 12
	if ok {
		i, err := strconv.ParseInt(postsPerRowStr, 10, 32)
		if err == nil {
			postsPerRow = int32(i)
		}
	}

	session := sessions.Default(context)
	userLevel, ok := session.Get("userlevel").(int32)
	if !ok {
		userLevel = 0
	}

	var posts []db.Post
	var err error
	if lowestID != 0 {
		params := db.GetOlderPostsParams{
			ID:        lowestID,
			UserLevel: userLevel,
			Limit:     postsPerRow*10 + 1,
		}

		posts, err = store.GetOlderPosts(context, params)
	}

	if highestID != 0 {
		params := db.GetNewerPostsParams{
			ID:        highestID,
			UserLevel: userLevel,
			Limit:     postsPerRow*10 + 1,
		}

		posts, err = store.GetNewerPosts(context, params)
	}

	if highestID == 0 && lowestID == 0 {
		posts, err = store.GetPosts(context, userLevel)
	}
	if err != nil {
		return nil, err
	}

	var postsJSON []PostJSON
	for _, post := range posts {
		var p PostJSON = PostJSON{}
		p.PopulatePost(post)
		postsJSON = append(postsJSON, p)
	}

	return &postsJSON, nil
}

type PostsQueries struct {
	MinimumID   string `query:"lowestID"`
	MaximumID   string `query:"highestID"`
	PostsPerRow string `query:"postsPerRow"`
}

// GetPosts returns Posts with voting information
func GetPostsFibre(store db.Store, c *fiber.Ctx) (*[]PostJSON, error) {
	queries := &PostsQueries{}

	if err := c.QueryParser(queries); err != nil {
		return nil, err
	}

	lowestID, _ := strconv.ParseInt(queries.MinimumID, 10, 32)
	highestID, _ := strconv.ParseInt(queries.MaximumID, 10, 32)
	postsPerRow, _ := strconv.ParseInt(queries.PostsPerRow, 10, 32)

	_ = postsPerRow
	// session := sessions.Default(c)
	// userLevel, ok := session.Get("userlevel").(int32)
	// if !ok {
	// 	userLevel = 0
	// }

	var userLevel int32 = 0

	var posts []db.Post
	var err error
	// if lowestID != 0 {
	// 	params := db.GetOlderPostsParams{
	// 		ID:        int32(lowestID),
	// 		UserLevel: userLevel,
	// 		Limit:     int32(postsPerRow)*10 + 1,
	// 	}

	// 	posts, err = store.GetOlderPosts(c, params)
	// }

	// if highestID != 0 {
	// 	params := db.GetNewerPostsParams{
	// 		ID:        highestID,
	// 		UserLevel: userLevel,
	// 		Limit:     postsPerRow*10 + 1,
	// 	}

	// 	posts, err = store.GetNewerPosts(c, params)
	// }

	if highestID == 0 && lowestID == 0 {
		posts, err = store.GetPosts(c.Context(), userLevel)
	}
	if err != nil {
		return nil, err
	}

	var postsJSON []PostJSON
	for _, post := range posts {
		var p PostJSON = PostJSON{}
		p.PopulatePost(post)
		postsJSON = append(postsJSON, p)
	}

	return &postsJSON, nil
}

// GetDuplicatePosts retrieves Duplicate Posts from the Storage with a hash
func GetDuplicatePosts(store db.Store, context *gin.Context, hash *goimagehash.ExtImageHash) ([]db.GetPossibleDuplicatePostsRow, error) {
	params := db.GetPossibleDuplicatePostsParams{
		Column1: hash.GetHash()[0],
		Column2: hash.GetHash()[1],
		Column3: hash.GetHash()[2],
		Column4: hash.GetHash()[3],
	}

	return store.GetPossibleDuplicatePosts(context, params)
}
