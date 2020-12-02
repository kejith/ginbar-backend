package utils

import (
	"errors"
	"fmt"
	"image"

	//"image/gif"
	"bytes"
	"image/jpeg"
	"net/http"
	"os"
	"time"

	//"encoding/base64"
	"path/filepath"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"

	"gopkg.in/gographics/imagick.v2/imagick"
)

// GetCropDimensions returns the Dimensions needed for Cropping
func GetCropDimensions(img *image.Image, width, height int) (int, int) {
	// if we don't have width or height set use the smaller image dimension
	// as both width and height
	if width == 0 && height == 0 {
		bounds := (*img).Bounds()
		x := bounds.Dx()
		y := bounds.Dy()
		if x < y {
			width = x
			height = x
		} else {
			width = y
			height = y
		}
	}
	return width, height
}

// DownloadImage downloads the image and decodes it
func DownloadImage(url string) (img image.Image, format string, err error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	img, format, err = image.Decode(response.Body)
	if err != nil {
		return
	}

	return
}

// ProcessUploadedImage ... TODO
func ProcessUploadedImage(url string) (fileName string, err error) {
	// create Filepaths
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	imageDir := filepath.Join(cwd, "public", "images")
	thumbnailDir := filepath.Join(imageDir, "thumbnails")

	// load image
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}

	img, format, err := image.Decode(response.Body)
	if err != nil {
		return "", err
	}

	fileName = fmt.Sprintf("%v", time.Now().UnixNano())
	//imgFileName := fmt.Sprintf("%v.%s", time.Now().UnixNano(), "png") // TODO put user id into filename to be save for duplicates
	imgFile, err := SaveImage(imageDir, fileName, format, &img)
	if err != nil {
		if imgFile != nil {
			os.Remove(imgFile.Name())
			imgFile.Close()
		}

		return "", err
	}
	defer imgFile.Close()

	imgCropped, err := CropImage(&img, 250, 250)
	if err != nil {
		os.Remove(imgFile.Name())
		imgFile.Close()

		return "", err
	}

	//imgFileName := fmt.Sprintf("%v.%s", , format)
	thumbnailFile, err := SaveImage(thumbnailDir, fileName, format, &imgCropped)
	if err != nil {
		if thumbnailFile != nil {
			os.Remove(thumbnailFile.Name())
			thumbnailFile.Close()
		}

		os.Remove(imgFile.Name())
		imgFile.Close()

		return "", err
	}

	return imgFile.Name(), nil
}

// SaveImage the image to the disk
func SaveImage(cwd, name, format string, image *image.Image) (file *os.File, err error) {
	imagick.Initialize()
	// Schedule cleanup
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()

	fileName := fmt.Sprintf("%s.%s", name, "png")
	filePath := filepath.Join(cwd, fileName)

	file, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	err = jpeg.Encode(&buff, *image, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Imagick")
	err = mw.ReadImageBlob(buff.Bytes())
	if err != nil {
		return nil, err
	}
	mw.SetFormat("jpeg")

	mw.SetCompression(imagick.COMPRESSION_BZIP)
	mw.BlurImage(5, 0.05)

	err = mw.SetImageCompressionQuality(85)
	if err != nil {
		return nil, err
	}
	err = mw.WriteImageFile(file)
	if err != nil {
		return nil, err
	}
	return
}

// CropImage crops the given Image with a smart cropper that calculates
// the best position to crop the image. i.e. a face or a distinct object
func CropImage(imgIn *image.Image, w int, h int) (img image.Image, err error) {
	width, height := GetCropDimensions(imgIn, w, h)
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	bestCrop, err := analyzer.FindBestCrop(*imgIn, 240, 240)

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := (*imgIn).(subImager)
	if !ok {
		err = errors.New("image does not support cropping")
		return
	}

	img = simg.SubImage(bestCrop)

	if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
		img = resizer.Resize(img, uint(width), uint(height))
	}

	return
}
