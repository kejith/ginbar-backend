package models

type PostsQueries struct {
	MinimumID   string `query:"lowestID" form:"lowestID"`
	MaximumID   string `query:"highestID" form:"highestID"`
	PostsPerRow string `query:"postsPerRow" form:"postsPerRow"`
}

type PostQueries struct {
	ID string `query:"post_id" form:"post_id"`
}
