package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ProcessUploadedVideo saves the uploaded video to disk and creates a thumbnail
func ProcessUploadedVideo(file multipart.File, format string, dirs Directories) (fileName string, thumbnailFilename string, err error) {
	name := fmt.Sprintf("%v", time.Now().UnixNano())

	// save uploaded video file into video directory
	videoFilePath, err := SaveMultipartFile(file, name, format, dirs.Video)
	if err != nil {
		return "", "", err
	}

	thumbnailFilename, err = CreateVideoThumbnail(videoFilePath, name, dirs)

	return filepath.Base(videoFilePath), thumbnailFilename, nil
}

// SaveMultipartFile takes a multipart File and saves it to the disk and returns
// a filepath to the saved file
func SaveMultipartFile(file multipart.File, name string, format string, directory string) (filePath string, err error) {
	fileName := fmt.Sprintf("%v.%s", name, format)
	filePath = filepath.Join(directory, fileName)

	localFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// CreateVideoThumbnail creates a Thumbnail from a Video File
func CreateVideoThumbnail(inputFilePath string, name string, dirs Directories) (filename string, err error) {
	filename = fmt.Sprintf("%s.jpeg", name)
	tmpThumbnailFilePath := filepath.Join(dirs.Tmp, filename)
	commandArgs := fmt.Sprintf("-i %s -ss 00:00:01.000 -vframes 1 %s -hide_banner -loglevel panic", inputFilePath, tmpThumbnailFilePath)
	cmd := exec.Command("ffmpeg", strings.Split(commandArgs, " ")...)
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	err = CreateThumbnailFromFile(tmpThumbnailFilePath, filepath.Join(dirs.Thumbnail, filename))

	if err != nil {
		return "", err
	}
	return "", nil
}
