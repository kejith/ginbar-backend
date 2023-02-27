package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ProcessVideo processes a video from the disk. It moves the Input File
// to the respective directory and creates a thumbnail
func ProcessVideo(inputFilePath, format string, dirs Directories) (fileName string, thumbnailFilename string, err error) {
	ext := filepath.Ext(inputFilePath)
	dstFileName := GenerateFilename(ext)
	dst := filepath.Join(dirs.Video, dstFileName)

    // move Video from tmp to public
    cmd := exec.Command("mv", inputFilePath, dst)
    if err = cmd.Run(); err != nil {
        return "", "", fmt.Errorf("could not move Video from TMP Dir to Video Dir: %w", err)
    }


	thumbnailFilename = dstFileName[0:len(dstFileName)-len(ext)]
	// create Thumbnail
	thumbnailFilename, err = CreateVideoThumbnail(dst, thumbnailFilename, dirs)
	if err != nil {
		return "", "", fmt.Errorf("could not create a Thumbnail for the Video: %w", err)
	}

	return dstFileName, thumbnailFilename, nil
}

// ProcessUploadedVideo saves the uploaded video to disk and creates a thumbnail
func ProcessUploadedVideo(file *multipart.File, format string, dirs Directories) (fileName string, thumbnailFilename string, err error) {
	name := fmt.Sprintf("%v", time.Now().UnixNano())

	// save uploaded video file into video directory
	videoFilePath, err := SaveMultipartFile(file, name, format, dirs.Video)
	if err != nil {
		return "", "", err
	}

	thumbnailFilename, err = CreateVideoThumbnail(videoFilePath, name, dirs)

	return filepath.Base(videoFilePath), thumbnailFilename, nil
}

// ProcessVideoFromURL saves the uploaded video to disk and creates a thumbnail
func ProcessVideoFromURL(response *http.Response, format string, dirs Directories) (fileName string, thumbnailFilename string, err error) {
	name := fmt.Sprintf("%v", time.Now().UnixNano())

	// save uploaded video file into video directory
	videoFilePath, err := SaveVideoFromURL(response, name, format, dirs.Video)
	if err != nil {
		return "", "", fmt.Errorf("Saving Video from URL: %v", err)
	}

	thumbnailFilename, err = CreateVideoThumbnail(videoFilePath, name, dirs)
	if err != nil {
		return "", "", fmt.Errorf("Creating Video Thumbnail: %v", err)
	}

	return filepath.Base(videoFilePath), thumbnailFilename, nil
}

// SaveMultipartFile takes a multipart File and saves it to the disk and returns
// a filepath to the saved file
func SaveMultipartFile(file *multipart.File, name string, format string, directory string) (filePath string, err error) {
	fileName := fmt.Sprintf("%v.%s", name, format)
	filePath = filepath.Join(directory, fileName)

	localFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, *file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// SaveVideoFromURL ...
func SaveVideoFromURL(response *http.Response, name string, format string, directory string) (filePath string, err error) {
	fileName := fmt.Sprintf("%v.%s", name, format)
	filePath = filepath.Join(directory, fileName)

	localFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, response.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// CreateVideoThumbnail creates a Thumbnail from a Video File
func CreateVideoThumbnail(inputFilePath string, name string, dirs Directories) (filename string, err error) {
	filename = fmt.Sprintf("%s.jpg", name)
	webpFilename := fmt.Sprintf("%s.webp", name)
	tmpThumbnailFilePath := filepath.Join(dirs.Tmp, filename)
	commandArgs := fmt.Sprintf("-i %s -ss 00:00:01.000 -vframes 1 %s -hide_banner -loglevel panic", inputFilePath, tmpThumbnailFilePath)
	cmd := exec.Command("ffmpeg", strings.Split(commandArgs, " ")...)
	
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	
	err = cmd.Run()


	if err != nil {
		return "", err
	}
	
	err = CreateThumbnailFromFile(tmpThumbnailFilePath, filepath.Join(dirs.Thumbnail, webpFilename), dirs)

	if err != nil {
		return "", err
	}
	return webpFilename, nil
}
