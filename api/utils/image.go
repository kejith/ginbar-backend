package utils

import (
	"bufio"
	"fmt"
	"image"
	"io"
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

// --------------------------------------
// STRUCTs
// --------------------------------------

// ImageProcessResult yields the return data for the image processing function
type ImageProcessResult struct {
	Filename          string
	ThumbnailFilename string
	UploadedFilename  string
	PerceptionHash    *goimagehash.ExtImageHash
}

// --------------------------------------
// Process
// --------------------------------------

// ProcessImageFromURL ... TODO
func ProcessImageFromURL(
	url string,
	format string,
	dirs Directories,
) (result *ImageProcessResult, err error) {
	fileName := fmt.Sprintf("%v.webp", time.Now().UnixNano())

	tmpFilePath, err := DownloadImage(url, fileName, dirs.Tmp)
	if err != nil {
		return &ImageProcessResult{}, fmt.Errorf("Image download failed: %w", err)
	}

	return ProcessImage(tmpFilePath, dirs)

}

// ProcessImageFromMultipart ... TODO
func ProcessImageFromMultipart(
	file *os.File,
	format string,
	dirs Directories,
) (result *ImageProcessResult, err error) {
	fileName := GenerateFilename("webp")
	filePath := filepath.Join(dirs.Upload, fileName)

	processResult, err := ProcessImage(filePath, dirs)
	if err != nil {
		return &ImageProcessResult{}, fmt.Errorf("Image Processing failed: %w", err)
	}

	processResult.UploadedFilename = filePath
	return processResult, err
}

// ProcessImage ...
func ProcessImage(inputFilePath string, dirs Directories) (result *ImageProcessResult, err error) {
	fileName := filepath.Base(inputFilePath)

	outputFilePath := filepath.Join(dirs.Image, fileName)
	err = ConvertImageToWebp(inputFilePath, outputFilePath, 75)
	if err != nil {
		return nil, fmt.Errorf("Convert Image to WEBP failed: %w", err)
	}

	img, err := LoadImageFile(inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("Loading Image(%s) from disk failed: %w", inputFilePath, err)
	}

	hash, err := goimagehash.ExtPerceptionHash(*img, 16, 16)
	if err != nil {
		return nil, fmt.Errorf("Perception Hash generation failed: %w", err)
	}

	outputThumbnailFilePath := filepath.Join(dirs.Thumbnail, fileName)
	if err = CreateThumbnailFromImage(img, outputThumbnailFilePath, dirs); err != nil {
		return nil, fmt.Errorf("Thumbnail creation from Image failed: %w", err)
	}

	result = &ImageProcessResult{
		Filename:          filepath.Base(outputFilePath),
		ThumbnailFilename: filepath.Base(outputThumbnailFilePath),
		PerceptionHash:    hash,
	}

	return result, err
}

// LoadImageFile opens an image from the disk and decodes it into image.Image
func LoadImageFile(inputFilePath string) (*image.Image, error) {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return nil, fmt.Errorf("Image File couldn't be opened: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("Image decoding failed: %w", err)
	}

	return &img, err
}

// SaveImageJPEG saves an image.Image to the disk as a JPEG
func SaveImageJPEG(img *image.Image, directory string, name string) (filePath string, err error) {
	filePath = filepath.Join(directory, name)
	file, err := os.Create(filePath)
	if err != nil {
		return filePath, fmt.Errorf("Image File creation failed: %w", err)
	}
	defer file.Close()

	err = jpeg.Encode(file, *img, &jpeg.Options{Quality: 100})
	if err != nil {
		if errRemove := os.Remove(filePath); errRemove != nil {
			err = fmt.Errorf("Couldnt remove created File: %w", err)
		}
		return filePath, fmt.Errorf("Encoding image failed: %w", err)
	}

	return filePath, nil
}

// CreateThumbnailFromFile Reads and inputFilePath Image from the Disk and Creates a Thumbnail
// in the outputFilePath
func CreateThumbnailFromFile(inputFilePath string, dstFilePath string, dirs Directories) (err error) {
	img, err := LoadImageFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("Loading image from disk failed: %w", err)
	}

	if err = CreateThumbnailFromImage(img, dstFilePath, dirs); err != nil {
		return fmt.Errorf("Thumbnail creation from Image failed: %w", err)
	}

	return nil
}

// CreateThumbnailFromImage takes and image.Image an saves a thumbnail of it
// to the disk in the outputFilePath
func CreateThumbnailFromImage(img *image.Image, dstFilePath string, dirs Directories) (err error) {

	imgCropped, err := CropImage(img, 150, 150)
	if err != nil {
		return fmt.Errorf("Cropping image failed: %w", err)
	}

	// create a temporary image file that we can convert
	tmpFilePath, err := SaveImageJPEG(imgCropped, filepath.Join(dirs.Tmp, "thumbnails"), GenerateFilename("jpeg"))
	if err != nil {
		return fmt.Errorf("Saving Image to disk as a JPEG failed: %w", err)
	}

	err = ConvertImageToWebp(tmpFilePath, dstFilePath, 75)
	if err != nil {
		return fmt.Errorf("Can't Convert Image to Webp: %w", err)
	}

	return nil
}

// --------------------------------------
// HELPER
// --------------------------------------

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

// DownloadImage ...
func DownloadImage(
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
		"\"%s\" -q %v -preset picture -m 6 -mt -o \"%s\"",
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
func CropImage(imgIn *image.Image, w int, h int) (i *image.Image, err error) {
	width, height := GetCropDimensions(imgIn, w, h)
	resizer := nfnt.NewDefaultResizer()
	analyzer := smartcrop.NewAnalyzer(resizer)
	bestCrop, err := analyzer.FindBestCrop(*imgIn, width, height)

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	simg, ok := (*imgIn).(subImager)
	if !ok {
		return nil, fmt.Errorf("Image does not support cropping: %w", err)
	}

	img := simg.SubImage(bestCrop)

	if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
		img = resizer.Resize(img, uint(width), uint(height))
	}

	return &img, nil
}

// GenerateFilename generates a new Unique Filename
// with a given Format Extension
func GenerateFilename(fileFormat string) string {
	fileName := fmt.Sprintf("%v.%s", time.Now().UnixNano(), fileFormat)
	return fileName
}

// // ProcessImage saves an image to the disk and creates a thumbnail
// func ProcessImage(img *image.Image, format string, dirs Directories) (result ImageProcessResult, err error) {
// 	fileName := fmt.Sprintf("%v.webp", time.Now().UnixNano())

// 	hash, err := goimagehash.ExtPerceptionHash(*img, 16, 16)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	imgFile, err := SaveImage(filepath.Join(dirs.Image, fileName), img, 65)
// 	if err != nil {
// 		if imgFile != nil {
// 			os.Remove(imgFile.Name())
// 			imgFile.Close()
// 		}

// 		return ImageProcessResult{}, err
// 	}
// 	defer imgFile.Close()

// 	imgCropped, err := CropImage(img, 150, 150)
// 	if err != nil {
// 		os.Remove(imgFile.Name())
// 		imgFile.Close()

// 		return ImageProcessResult{}, err
// 	}

// 	//imgFileName := fmt.Sprintf("%v.%s", , format)
// 	thumbnailFile, err := SaveImage(filepath.Join(dirs.Thumbnail, fileName), &imgCropped, 65)
// 	if err != nil {
// 		if thumbnailFile != nil {
// 			os.Remove(thumbnailFile.Name())
// 			thumbnailFile.Close()
// 		}

// 		os.Remove(imgFile.Name())
// 		imgFile.Close()

// 		return ImageProcessResult{}, err
// 	}

// 	return ImageProcessResult{Filename: imgFile.Name(), ThumbnailFilename: thumbnailFile.Name(), PerceptionHash: hash}, nil
// }
