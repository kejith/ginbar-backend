package api

import (
	"fmt"
	"ginbar/api/models"
	"ginbar/api/utils"
	"ginbar/mysql/db"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// RegenerateThumbnails regenerates all the thumbnails from the images
// saved in the umage directory
func (server *Server) RegenerateThumbnails(context *gin.Context) {
	posts, err := models.GetPosts(server.store, context)
	if err != nil {
		log.Print(err)
		return
	}
	imageDir := server.directories.Image
	//thumbDir := server.directories.Thumbnail
	videoDir := server.directories.Video

	imageFiles, err := ioutil.ReadDir(imageDir)
	if err != nil {
		log.Print(err)
		return
	}
	_ = imageFiles

	//filesInDatabase := make(map[string]int)
	length := len(posts)
	i := 1
	for _, post := range posts {
		fileName := post.FileName

		if post.ContentType == "video/mp4" {
			fmt.Println(i, length, filepath.Join(imageDir, fileName))
			// 	i = i + 1
			// 	err = utils.CreateThumbnailFromFile(
			// 		filepath.Join(imageDir, fileName),
			// 		filepath.Join(thumbDir, fileName))

			// 	if err != nil {
			// 		fmt.Println(err)
			// 	}
			// } else {
			videoFilePath := filepath.Join(videoDir, fileName)
			//fmt.Println(videoFilePath, strings.TrimSuffix(fileName, path.Ext(fileName)))
			_, err := utils.CreateVideoThumbnail(
				videoFilePath,
				strings.TrimSuffix(fileName, path.Ext(fileName)),
				server.directories)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

// RedownloadAndCompressImages ...
func (server *Server) RedownloadAndCompressImages(context *gin.Context) {
	posts, err := models.GetPosts(server.store, context)
	if err != nil {
		log.Print(err)
		return
	}

	length := len(posts)
	count := 0
	for _, post := range posts {
		url := post.URL

		count = 1 + count
		fmt.Println(count, length, url)

		if post.ContentType == "image" {
			response, _, fileFormat, err := utils.LoadFileFromURL(url)
			filename, thumbnailFilename, err := utils.ProcessImageFromURL(
				response,
				fileFormat,
				server.directories,
			)

			if err != nil {
				fmt.Println(err)
			}

			params := db.UpdatePostFilesParams{
				Filename:          filepath.Base(filename),
				ThumbnailFilename: filepath.Base(thumbnailFilename),
				ID:                int32(post.ID),
			}

			err = server.store.UpdatePostFiles(context, params)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
