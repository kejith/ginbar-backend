package api

import (
	"fmt"
	"ginbar/api/utils"
	"ginbar/mysql/db"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/gin-gonic/gin"
)

// RegenerateThumbnails regenerates all the thumbnails from the images
// saved in the umage directory
func (server *Server) RegenerateThumbnails(context *gin.Context) {
	posts, err := server.store.GetAllPosts(context)
	if err != nil {
		log.Print(err)
		return
	}
	imageDir := server.directories.Image
	thumbDir := server.directories.Thumbnail
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
		fileName := post.Filename

		i = i + 1
		if post.ContentType == "image" {
			fmt.Println(i, length, filepath.Join(imageDir, fileName))
			err = utils.CreateThumbnailFromFile(
				filepath.Join(imageDir, fileName),
				filepath.Join(thumbDir, fileName))

			if err != nil {
				fmt.Println(err)
			}
		} else {
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
	posts, err := server.store.GetImagePosts(context)
	if err != nil {
		log.Print(err)
		return
	}

	length := len(posts)
	count := 0
	for _, post := range posts {
		url := post.Url

		count = 1 + count
		fmt.Println(count, length, url)

		if post.ContentType == "image" && url != "" {
			response, _, fileFormat, err := utils.LoadFileFromURL(url)
			processResult, err := utils.ProcessImageFromURL(
				response,
				fileFormat,
				server.directories,
			)

			if err != nil {
				fmt.Println(err)
			}

			params := db.UpdatePostFilesParams{
				Filename:          filepath.Base(processResult.Filename),
				ThumbnailFilename: filepath.Base(processResult.ThumbnailFilename),
				ID:                int32(post.ID),
			}

			err = server.store.UpdatePostFiles(context, params)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

// RecalculateHashes iterates over every existing Image Post and calculates
// the hashes and updates its value in the storage
func (server *Server) RecalculateHashes(context *gin.Context) {
	posts, err := server.store.GetImagePosts(context)
	if err != nil {
		log.Print(err)
		return
	}

	length := len(posts)
	count := 0
	for _, post := range posts {
		count = 1 + count

		hash, _ := PerceptionHashFromFile(filepath.Join(server.directories.Image, post.Filename))
		fmt.Println(count, length, post.Filename, hash)

		params := db.UpdatePostHashesParams{
			ID:     int32(post.ID),
			PHash0: hash[0],
			PHash1: hash[1],
			PHash2: hash[2],
			PHash3: hash[3],
		}

		err = server.store.UpdatePostHashes(context, params)
		if err != nil {
			fmt.Println(err)
		}

	}
}

// PerceptionHashFromFile calculates a perception hash from a given image
// file
func PerceptionHashFromFile(filepath string) (hashes []uint64, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	imageHash, err := goimagehash.ExtPerceptionHash(img, 16, 16)
	if err != nil {
		return nil, err
	}

	return imageHash.GetHash(), nil
}
