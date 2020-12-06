package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ginbar/api/utils"
	"ginbar/mysql/db"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests and Stores Connections, Sessions and State
type Server struct {
	store              db.Store
	router             *gin.Engine
	sessions           sessions.CookieStore
	directories        utils.Directories
	postsResponseCache *persistence.InMemoryStore
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) (*Server, error) {
	var secret = "IX~|xTE@4*v@e95sLll4g`#6G288be"
	// setup directory paths
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	directories := utils.Directories{
		CWD:       cwd,
		Image:     filepath.Join(cwd, "public", "images"),
		Thumbnail: filepath.Join("public", "images", "thumbnails"),
		Video:     filepath.Join(cwd, "public", "videos"),
		Tmp:       filepath.Join(cwd, "tmp"),
	}

	// create Server
	server := &Server{
		store:       store,
		router:      gin.New(),
		sessions:    sessions.NewCookieStore([]byte(secret)),
		directories: directories,
	}

	//
	// SETUP SERVER
	//

	// Middleware
	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://kejith.de"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	server.router.Use(sessions.Sessions("gbsession", server.sessions))
	server.postsResponseCache = persistence.NewInMemoryStore(time.Second)

	// UTILS
	groupUtils := server.router.Group("/utils")
	groupUtils.Use(gin.Logger())
	{
		groupUtils.GET("/thumbnails/regenerate", server.RegenerateThumbnails)
		groupUtils.GET("/reprocess/urls", server.RedownloadAndCompressImages)
	}

	// API
	groupAPI := server.router.Group("/api")
	groupAPI.Use(gin.Logger())
	groupAPI.Use(gin.Recovery())

	// APi/POST
	groupPost := groupAPI.Group("/post")
	{
		groupPost.GET("/", cache.CachePage(server.postsResponseCache, time.Minute, server.GetAll))
		groupPost.GET("/:post_id", cache.CachePage(server.postsResponseCache, time.Minute*30, server.Get))
		groupPost.GET("/:post_id/comments", server.GetComments)
		groupPost.POST("/create", server.CreatePost)
		groupPost.POST("/upload", server.UploadPost)
		groupPost.POST("/vote", server.VotePost)
	}

	// API/USER
	groupUser := groupAPI.Group("/user")
	{
		groupUser.GET("/", server.GetUsers)
		groupUser.GET("/:id", server.GetUser)
		groupUser.POST("/create", server.CreateUser)
		groupUser.POST("/login", server.Login)
		groupUser.POST("/logout", server.UserLogout)
	}

	// API/COMMENT
	groupComments := groupAPI.Group("/comment")
	{
		groupComments.POST("/create", server.CreateComment)
		groupComments.POST("/vote", server.VoteComment)
	}

	// API/TAG
	groupTags := groupAPI.Group("/tag")
	{
		groupTags.POST("/create", server.CreatePostTag)
		groupTags.POST("/vote", server.VotePostTag)
	}

	// API/CHECK
	groupCheck := groupAPI.Group("/check")
	//groupCheck.Use(AuthRequired)
	{
		groupCheck.GET("/me", server.Me)
	}

	// Serve static files from ./public/
	server.router.Use(static.Serve("/", static.LocalFile("./public", true)))
	server.router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	return server, nil

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

/*
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
*/

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	// server.router.Run(address)
	return server.router.RunTLS(":443", "./kejith.de.pem", "./kejith.de.key")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
