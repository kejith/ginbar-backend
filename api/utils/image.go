package utils

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	//"image/gif"

	"image/jpeg"
	"net/http"
	"os"
	"time"

	//"encoding/base64"

	"github.com/corona10/goimagehash"
	"github.com/harukasan/go-libwebp/webp"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
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
func CreateThumbnailFromFile(inputFilePath string, outputfilePath string, dirs Directories) (err error) {

	sourceFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	img, _, err := image.Decode(sourceFile)
	if err != nil {
		return err
	}

	imgCropped, err := CropImage(&img, 150, 150)
	if err != nil {
		return err
	}

	tmpThumbnailFilePath := filepath.Join(dirs.Tmp, "thumbnails", GenerateFilename("jpeg"))
	tmpFile, err := os.Create(tmpThumbnailFilePath)
	if err != nil {
		return err
	}

	err = jpeg.Encode(tmpFile, imgCropped, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}

	err = ConvertImageToWebp(
		tmpThumbnailFilePath,
		outputfilePath,
		75,
	)

	if err != nil {
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

// ImageProcessResult yields the return data for the image processing function
type ImageProcessResult struct {
	Filename          string
	ThumbnailFilename string
	UploadedFilename  string
	PerceptionHash    *goimagehash.ExtImageHash
}

// ProcessImage saves an image to the disk and creates a thumbnail
func ProcessImage(img *image.Image, format string, dirs Directories) (result ImageProcessResult, err error) {
	fileName := fmt.Sprintf("%v.webp", time.Now().UnixNano())

	hash, err := goimagehash.ExtPerceptionHash(*img, 16, 16)
	if err != nil {
		fmt.Println(err)
	}

	imgFile, err := SaveImage(filepath.Join(dirs.Image, fileName), img, 65)
	if err != nil {
		if imgFile != nil {
			os.Remove(imgFile.Name())
			imgFile.Close()
		}

		return ImageProcessResult{}, err
	}
	defer imgFile.Close()

	imgCropped, err := CropImage(img, 150, 150)
	if err != nil {
		os.Remove(imgFile.Name())
		imgFile.Close()

		return ImageProcessResult{}, err
	}

	//imgFileName := fmt.Sprintf("%v.%s", , format)
	thumbnailFile, err := SaveImage(filepath.Join(dirs.Thumbnail, fileName), &imgCropped, 65)
	if err != nil {
		if thumbnailFile != nil {
			os.Remove(thumbnailFile.Name())
			thumbnailFile.Close()
		}

		os.Remove(imgFile.Name())
		imgFile.Close()

		return ImageProcessResult{}, err
	}

	return ImageProcessResult{Filename: imgFile.Name(), ThumbnailFilename: thumbnailFile.Name(), PerceptionHash: hash}, nil
}

// ProcessImageNew ...
func ProcessImageNew(innputFilePath string, dirs Directories) (result *ImageProcessResult, err error) {
	fileName := filepath.Base(innputFilePath)

	outputFilePath := filepath.Join(dirs.Image, fileName)
	err = ConvertImageToWebp(innputFilePath, outputFilePath, 75)
	if err != nil {
		return nil, err
	}

	downloadedFile, err := os.Open(innputFilePath)
	if err != nil {
		return nil, err
	}
	defer downloadedFile.Close()

	img, _, err := image.Decode(downloadedFile)
	if err != nil {
		return nil, err
	}

	hash, err := goimagehash.ExtPerceptionHash(img, 16, 16)
	if err != nil {
		return nil, err
	}

	imgCropped, err := CropImage(&img, 150, 150)
	if err != nil {
		return nil, err
	}

	outputThumbnailTmpFilePath := filepath.Join(dirs.Tmp, "thumbnails", fileName)
	file, err := os.Create(outputThumbnailTmpFilePath)
	if err != nil {
		return nil, err
	}

	err = jpeg.Encode(file, imgCropped, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}

	outputThumbnailFilePath := filepath.Join(dirs.Thumbnail, fileName)
	err = ConvertImageToWebp(
		outputThumbnailTmpFilePath,
		outputThumbnailFilePath,
		100,
	)
	if err != nil {
		return nil, err
	}

	result = &ImageProcessResult{
		Filename:          fileName,
		ThumbnailFilename: filepath.Base(outputThumbnailFilePath),
		PerceptionHash:    hash,
	}

	return result, err
}

// ProcessImageFromURL ... TODO
func ProcessImageFromURL(
	url string,
	format string,
	dirs Directories,
) (result *ImageProcessResult, err error) {
	fileName := fmt.Sprintf("%v.webp", time.Now().UnixNano())

	tmpFilePath, err := DownloadImageNew(url, fileName, dirs.Tmp)
	if err != nil {
		return &ImageProcessResult{}, err
	}

	return ProcessImageNew(tmpFilePath, dirs)

}

// ProcessImageFromMultipart ... TODO
func ProcessImageFromMultipart(
	file *multipart.File,
	format string,
	dirs Directories,
) (result *ImageProcessResult, err error) {
	fileName := GenerateFilename("webp")
	filePath := filepath.Join(dirs.Upload, fileName)

	dst, err := os.Create(filePath)
	defer dst.Close()
	if err != nil {
		return &ImageProcessResult{}, err
	}

	_, err = io.Copy(dst, *file)
	if err != nil {
		return &ImageProcessResult{}, err
	}

	processResult, err := ProcessImageNew(filePath, dirs)
	processResult.UploadedFilename = filePath
	return processResult, err
}

// DownloadImageNew ...
func DownloadImageNew(
	URL string,
	fileName string,
	directory string,
) (filePath string, err error) {
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", err
	}

	fileName = path.Base(URL)

	filePath = filepath.Join(directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil

}

// ConvertImageToWebp ...
func ConvertImageToWebp(inputFilePath string, outputFilePath string, quality uint) error {
	commandArgs := fmt.Sprintf(
		"%s -q %v -preset picture -m 6 -mt -o -f 100 %s",
		inputFilePath,
		quality,
		outputFilePath)

	fmt.Println(commandArgs)
	cmd := exec.Command("cwebp", strings.Split(commandArgs, " ")...)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

// SaveImage the image to the disk
func SaveImage(filePath string, img *image.Image, quality uint) (file *os.File, err error) {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }

	file, err = os.Create(filePath)
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush()
		file.Close()
	}()

	config, err := webp.ConfigPreset(webp.PresetPhoto, float32(quality/100))
	if err != nil {
		return nil, err
	}

	err = webp.EncodeRGBA(writer, *img, config)
	if err != nil {
		return nil, err
	}

	return file, nil
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

// GenerateFilename generates a new Unique Filename
// with a given Format Extension
func GenerateFilename(fileFormat string) string {
	fileName := fmt.Sprintf("%v.%s", time.Now().UnixNano(), fileFormat)
	return fileName
}
