package utils

import (
	"errors"
	"fmt"
	"image"
	"mime/multipart"
	"path/filepath"

	//"image/gif"
	"bytes"
	"image/jpeg"
	"net/http"
	"os"
	"time"

	//"encoding/base64"

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

// CreateThumbnailFromFile Reads and inputFilePath Image from the Disk and Creates a Thumbnail
// in the outputFilePath
func CreateThumbnailFromFile(inputFilePath string, outputfilePath string) (err error) {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	img, err := jpeg.Decode(inputFile)
	if err != nil {
		return err
	}

	imgCropped, err := CropImage(&img, 150, 150)
	if err != nil {
		return err
	}

	//imgFileName := fmt.Sprintf("%v.%s", , format)
	thumbnailFile, err := SaveImage(outputfilePath, &imgCropped, 75)
	if err != nil {
		if thumbnailFile != nil {
			os.Remove(thumbnailFile.Name())
			thumbnailFile.Close()
		}

		return err
	}

	return nil
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

// ProcessImage saves an image to the disk and creates a thumbnail
func ProcessImage(img *image.Image, format string, dirs Directories) (fileName string, thumbnailFileName string, err error) {
	fileName = fmt.Sprintf("%v.jpeg", time.Now().UnixNano())
	//imgFileName := fmt.Sprintf("%v.%s", time.Now().UnixNano(), "png") // TODO put user id into filename to be save for duplicates
	imgFile, err := SaveImage(filepath.Join(dirs.Image, fileName), img, 75)
	if err != nil {
		if imgFile != nil {
			os.Remove(imgFile.Name())
			imgFile.Close()
		}

		return "", "", err
	}
	defer imgFile.Close()

	imgCropped, err := CropImage(img, 150, 150)
	if err != nil {
		os.Remove(imgFile.Name())
		imgFile.Close()

		return "", "", err
	}

	//imgFileName := fmt.Sprintf("%v.%s", , format)
	thumbnailFile, err := SaveImage(filepath.Join(dirs.Thumbnail, fileName), &imgCropped, 100)
	if err != nil {
		if thumbnailFile != nil {
			os.Remove(thumbnailFile.Name())
			thumbnailFile.Close()
		}

		os.Remove(imgFile.Name())
		imgFile.Close()

		return "", "", err
	}

	return imgFile.Name(), thumbnailFile.Name(), nil
}

// ProcessImageFromURL ... TODO
func ProcessImageFromURL(response *http.Response, format string, dirs Directories) (fileName string, thumbnailFileName string, err error) {
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return "", "", err
	}

	return ProcessImage(&img, format, dirs)

}

// ProcessImageFromMultipart ... TODO
func ProcessImageFromMultipart(file *multipart.File, format string, dirs Directories) (fileName string, thumbnailFileName string, err error) {
	img, _, err := image.Decode(*file)
	if err != nil {
		return "", "", err
	}
	filePath, thumbnailFilePath, err := ProcessImage(&img, format, dirs)
	return filepath.Base(filePath), filepath.Base(thumbnailFilePath), err
}

// SaveImage the image to the disk
func SaveImage(name string, image *image.Image, quality uint) (file *os.File, err error) {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	imagick.Initialize()
	// Schedule cleanup
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()

	//fileName := filepath.Base(name)
	filePath := name

	file, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer

	err = jpeg.Encode(&buff, *image, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}

	err = mw.ReadImageBlob(buff.Bytes())
	if err != nil {
		return nil, err
	}
	mw.SetFormat("jpeg")

	mw.SetCompression(imagick.COMPRESSION_BZIP)
	mw.BlurImage(5, 0.05)

	err = mw.SetImageCompressionQuality(quality)
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
	bestCrop, err := analyzer.FindBestCrop(*imgIn, width, height)

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
