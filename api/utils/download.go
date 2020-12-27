package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// LoadFileFromURL loads a file from an URL and returns a Reponse, FileType and
// FileFormat
func LoadFileFromURL(url string) (response *http.Response, fileType string, fileFormat string, err error) {
	response, err = http.Get(url)
	if err != nil {
		return nil, "", "", err
	}

	if response.StatusCode != 200 {

		return nil, "", "", errors.New("Received non 200 response code")
	}

	contentType := strings.Split(response.Header.Get("Content-Type"), "/")
	//contentLength := response.Header.Get("Content-Length")
	fileType, fileFormat = contentType[0], contentType[1]

	return response, fileType, fileFormat, nil
}

// DownloadFile downloads a file and saves it to the disk into the folder dir
// and returns the path to the file it created
func DownloadFile(url string, dir string) (filePath string, err error) {
	dst := filepath.Join(dir, path.Base(url))

	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("http get failed with status %v : %w", response.StatusCode, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("received a non-200 Status Code while Get Request %v : %w", response.StatusCode, err)
	}

	file, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("File could not be created: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", fmt.Errorf("Could not copy Response Body into Destination File: %w", err)
	}

	return dst, nil
}
