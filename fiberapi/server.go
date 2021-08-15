package fiberapi

import (
	"crypto/tls"
	"ginbar/mysql/db"
	"ginbar/utils"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// Server serves HTTP requests and Stores Connections, Sessions and State
type FiberServer struct {
	App         *fiber.App
	store       db.Store
	Sessions    *session.Store
	directories utils.Directories
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

func NewFiber(store db.Store) (*FiberServer, error) {
	directories := SetupDirectories()

	s := session.New()

	server := &FiberServer{
		App:         fiber.New(),
		store:       store,
		directories: directories,
		Sessions:    s,
	}

	// Register Middlewars
	//server.App.Use(cache.New())
	// server.App.Use(compress.New(compress.Config{
	// 	Level: compress.LevelBestCompression, // 2
	// }))

	// Register Dynamic Routes

	// Register Groups
	api := server.App.Group("/api")
	api.Use(logger.New())

	// API/USER
	userApi := api.Group("/user")
	userApi.Get("/:id", server.GetUser)
	userApi.Get("/*", server.GetUsers)
	userApi.Post("/login", server.Login)
	userApi.Post("/logout", server.Logout)
	api.Get("/check/me", server.Me)

	// API/POST
	postapi := api.Group("/post")
	// postapi.Use(cache.New(cache.Config{
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		return c.OriginalURL()
	// 	},
	// }))
	postapi.Get("/:post_id", server.GetPost)
	postapi.Get("/*", server.GetPosts)
	postapi.Post("/vote", server.VotePost)
	postapi.Post("/create", server.CreatePost)
	// postapi.Post("/create/multiple", server.CreateMultiplePosts)
	postapi.Post("/upload", server.UploadPost)

	// API/COMMENT
	commentapi := api.Group("/comment")
	commentapi.Post("/create", server.CreateComment)
	commentapi.Post("/vote", server.VoteComment)

	// API/TAG
	tagapi := api.Group("/tag")
	tagapi.Post("/create", server.CreatePostTag)
	tagapi.Post("/vote", server.VotePostTag)

	// Register Static Routes
	server.App.Static("/", "./public")
	server.App.Static("*", "./public/index.html")

	// Create TLS Certificate
	cer, err := tls.LoadX509KeyPair("fullchain.pem", "privkey.pem")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}

	// Create SSL Listener
	ln, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		panic(err)
	}

	log.Fatal(server.App.Listener(ln))

	return server, nil
}
