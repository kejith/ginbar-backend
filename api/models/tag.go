package models

import "github.com/kejith/ginbar-backend/com/kejith/ginbar-backend/mysql/db"

// PostTagJSON is a representation in JSON for the Database Object PostTag
type PostTagJSON struct {
	ID      int32  `json:"id"`
	Score   int32  `json:"score"`
	Name    string `json:"name"`
	PostID  int32  `json:"post_id"`
	UserID  int32  `json:"user_id"`
	Upvoted int8   `json:"upvoted"`
}

// Populate ...
func (p *PostTagJSON) Populate(post db.GetTagsByPostRow) {
	p.ID = post.ID
	p.Score = post.Score
	p.Name = post.Name
	p.PostID = post.PostID
	p.UserID = post.UserID
	upvoted, ok := post.Upvoted.(int64)

	if !ok {
		p.Upvoted = 0
	} else {
		p.Upvoted = int8(upvoted)
	}
}
