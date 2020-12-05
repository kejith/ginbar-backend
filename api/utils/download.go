package utils

import (
	"errors"
	"net/http"
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
