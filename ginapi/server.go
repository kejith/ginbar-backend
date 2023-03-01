package ginapi

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"ginbar/models"
	"ginbar/mysql/db"
	"ginbar/utils"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var secret = "IX~|xTE@4*v@e95sLll4g`#6G288be"

type BaseServer struct {
	store              db.Store
	sessions           sessions.CookieStore
	directories        utils.Directories
	postsResponseCache *persistence.InMemoryStore
}

// Server serves HTTP requests and Stores Connections, Sessions and State
type Server struct {
	BaseServer
	router *gin.Engine
}

type PostsJson struct {
	Posts *[]models.PostJSON `json:"posts"`
}

func SetupDirectories() utils.Directories {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	directories := utils.Directories{
		CWD:       cwd,
		Image:     filepath.Join(cwd, "public", "images"),
		Thumbnail: filepath.Join("public", "images", "thumbnails"),
		Video:     filepath.Join(cwd, "public", "videos"),
		Tmp:       filepath.Join(cwd, "tmp"),
		Upload:    filepath.Join(cwd, "public", "upload"),
	}

	return directories
}

// NewServer creates a new HTTP server and sets up routing.
func NewServer(store db.Store) (*Server, error) {

	directories := SetupDirectories()

	// create Server
	server := &Server{
		BaseServer: BaseServer{
			store:       store,
			sessions:    sessions.NewCookieStore([]byte(secret)),
			directories: directories,
		},
		router: gin.New(),
	}

	//
	// SETUP SERVER
	//

	// Middleware

	server.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://kejith.de", "http://kejith.de:3000", "https://ginbar.kejith.de"},
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
		groupUtils.GET("/regenerate/thumbnails", server.RegenerateThumbnails)
		groupUtils.GET("/regenerate/hashes", server.RecalculateHashes)
		groupUtils.GET("/reprocess/urls", server.RedownloadAndCompressImages)
		groupUtils.GET("/check/urls", server.CheckIfExternalMediaIsAvailable)
	}

	// API
	groupAPI := server.router.Group("/api")
	groupAPI.Use(gin.Logger())
	groupAPI.Use(gin.Recovery())

	// APi/POST
	groupPost := groupAPI.Group("/post")
	{
		// TODO fix cache by user
		// groupPost.GET("/", _cache.PageByUser(server.postsResponseCache, time.Minute, server.GetAll))
		// groupPost.GET("/:post_id", _cache.PageByUser(server.postsResponseCache, time.Minute*30, server.Get))
		groupPost.GET("/", cache.CachePage(server.postsResponseCache, time.Minute, server.GetAll))
		groupPost.GET("/:post_id", server.Get)
		// groupPost.GET("/", server.GetAll)
		// groupPost.GET("/:post_id", server.Get)
		groupPost.GET("/:post_id/comments", server.GetComments)
		groupPost.POST("/create", server.CreatePost)
		groupPost.POST("/create/multiple", server.CreateMultiplePosts)
		groupPost.POST("/upload", server.UploadPost)
		groupPost.POST("/delete", server.DeletePost)
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
	user, okUser := session.Get("user").(string)
	userID, okUserID := session.Get("userid").(uint)
	if !okUser || !okUserID || len(user) < 3 || userID <= 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

// Keys for Accessing Data from the Context
const (
	KeyUsername     = "username"
	KeyUserID       = "userid"
	KeyIsAuthorized = "isAuthorized"
)

// AuthenticationRequired checks if a user session is valid and sets the flag
// "username" and "userid" in the context
func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user, okUser := session.Get("user").(string)
		userID, okUserID := session.Get("userid").(uint)
		if !okUser || !okUserID || len(user) < 3 || userID <= 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(KeyUsername, user)
		c.Set(KeyUserID, userID)
	}
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
	// return server.router.RunTLS(":443", "./fullchain.pem", "./privkey.pem")

}

// func errorResponse(err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }
