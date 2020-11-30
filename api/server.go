package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"ginbar/mysql/db"
)

// Server serves HTTP requests
type Server struct {
	store    db.Store
	router   *gin.Engine
	sessions sessions.CookieStore
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) *Server {
	var secret = "IX~|xTE@4*v@e95sLll4g`#6G288be"
	server := &Server{store: store}
	server.router = gin.Default()
	server.router.Use(cors.Default())

	server.sessions = sessions.NewCookieStore([]byte(secret))
	server.router.Use(sessions.Sessions("gbsession", server.sessions))
	//router.Use(static.Serve("/", static.LocalFile("./public", true)))

	groupAPI := server.router.Group("/api")
	groupPost := groupAPI.Group("/post")
	{
		groupPost.GET("/", server.GetAll)
		groupPost.GET("/:post_id", server.Get)
		groupPost.GET("/:post_id/comments", server.GetComments)
		//groupPost.POST("/:post_id/comments/create", server.CreateComment)
		groupPost.POST("/create", server.CreatePost)
	}

	groupUser := groupAPI.Group("/user")
	{
		groupUser.GET("/", server.GetUsers)
		groupUser.GET("/:id", server.GetUser)
		groupUser.POST("/create", server.CreateUser)
		groupUser.POST("/login", server.Login)
		groupUser.POST("/logout", server.Logout)
	}

	groupComments := groupAPI.Group("/comment")
	{
		//groupComments.GET("/", comments.ShowAll)
		//groupComments.GET("/:id", comments.Get)
		groupComments.POST("/create", server.CreateComment)
		groupComments.POST("/vote", server.VoteComment)
	}

	groupCheck := groupAPI.Group("/check")
	groupCheck.Use(AuthRequired)
	{
		groupCheck.GET("/me", server.Me)
	}

	server.router.Use(static.Serve("/", static.LocalFile("./public", true)))

	return server
}

// AuthRequired ...
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
