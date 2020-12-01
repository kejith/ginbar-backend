package api

import (
	"net/http"

	"ginbar/mysql/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests and Stores Connections, Sessions and State
type Server struct {
	store    db.Store
	router   *gin.Engine
	sessions sessions.CookieStore
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) *Server {
	var secret = "IX~|xTE@4*v@e95sLll4g`#6G288be"

	// create a Store
	server := &Server{store: store}

	// Create and configurate the Router
	server.router = gin.Default()

	// CORS
	// TODO: maybe not needed, needs to be tested
	server.router.Use(cors.Default())

	// Create and Store a store for Sessions
	server.sessions = sessions.NewCookieStore([]byte(secret))
	server.router.Use(sessions.Sessions("gbsession", server.sessions))

	// API
	groupAPI := server.router.Group("/api")

	// APi/POST
	groupPost := groupAPI.Group("/post")
	{
		groupPost.GET("/", server.GetAll)
		groupPost.GET("/:post_id", server.Get)
		groupPost.GET("/:post_id/comments", server.GetComments)
		groupPost.POST("/create", server.CreatePost)
		groupPost.POST("/vote", server.VotePost)
	}

	// API/USER
	groupUser := groupAPI.Group("/user")
	{
		groupUser.GET("/", server.GetUsers)
		groupUser.GET("/:id", server.GetUser)
		groupUser.POST("/create", server.CreateUser)
		groupUser.POST("/login", server.Login)
		groupUser.POST("/logout", server.Logout)
	}

	// API/COMMENT
	groupComments := groupAPI.Group("/comment")
	{
		groupComments.POST("/create", server.CreateComment)
		groupComments.POST("/vote", server.VoteComment)
	}

	groupTags := groupAPI.Group("/tag")
	{
		groupTags.POST("/create", server.CreatePostTag)
		groupTags.POST("/vote", server.VotePostTag)
	}

	// API/CHECK
	groupCheck := groupAPI.Group("/check")
	groupCheck.Use(AuthRequired)
	{
		groupCheck.GET("/me", server.Me)
	}

	// Serve static files from ./public/
	server.router.Use(static.Serve("/", static.LocalFile("./public", true)))

	return server
}

// AuthRequired checks if the current Session is valid and the User is
// authenticated
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Next()
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
