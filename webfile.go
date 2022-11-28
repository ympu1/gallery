package main

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/sunshineplan/imgconv"
	"image"
	"os"
	"path/filepath"
	"strings"
)

type webFile struct {
	Name string
	Type string
	Path string
}

func (file *webFile) fillFields(osFile os.DirEntry, path string, config config) {
	file.Name = osFile.Name()
	file.Type = getFileType(osFile, config)
	file.Path = strings.ReplaceAll(path+"/"+file.Name, config.RootPath, "")
}

func (file *webFile) getThumb(config config) []byte {
	switch file.Type {
	case "image":
		return file.getImageThumb(config)
	case "folder":
		return file.getFolderThumb(config)
	default:
		return nil
	}
}

func (file *webFile) getImageThumb(config config) []byte {
	src, err := imgconv.Open(config.RootPath + file.Path)
	if err != nil {
		fmt.Println(err, file.Path)
		return nil
	}

	srcWidth := src.Bounds().Dx()
	srcHeight := src.Bounds().Dy()

	thumbWidthRatio := srcWidth / config.ThumbSize
	thumbHeightRatio := srcHeight / config.ThumbSize

	thumbWidth := config.ThumbSize
	thumbHeight := config.ThumbSize

	if thumbWidthRatio > thumbHeightRatio {
		thumbWidth = 0
	} else {
		thumbHeight = 0
	}

	thumb := imgconv.Resize(src, &imgconv.ResizeOption{Width: thumbWidth, Height: thumbHeight})

	thumbPositionX := (thumb.Bounds().Dx() - config.ThumbSize) / 2 * -1
	thumbPositionY := (thumb.Bounds().Dy() - config.ThumbSize) / 2 * -1

	outputImage := gg.NewContext(config.ThumbSize, config.ThumbSize)
	outputImage.DrawImage(thumb, thumbPositionX, thumbPositionY)

	var buf bytes.Buffer
	outputImage.EncodeJPG(&buf, nil)
	//outputImage.EncodePNG(&buf)
	return buf.Bytes()
}

func (file *webFile) getFolderThumb(config config) []byte {
	folderBG, _ := gg.LoadImage("static/folder.png")
	folderBG = imgconv.Resize(folderBG, &imgconv.ResizeOption{Width: config.ThumbSize})

	outputImage := gg.NewContext(config.ThumbSize, config.ThumbSize)
	outputImage.DrawImage(folderBG, 0, 0)

	subFilesThumbs := file.getSubfilesThubms(config)

	subFilesThumbPadding := 10
	subFileThumbSize := (config.ThumbSize - subFilesThumbPadding*3) / 2

	x, y := subFilesThumbPadding, subFilesThumbPadding
	drawSubfileThumbToOutput(subFilesThumbs[0], x, y, subFileThumbSize, subFileThumbSize, outputImage)

	x = subFilesThumbPadding*2 + subFileThumbSize
	drawSubfileThumbToOutput(subFilesThumbs[1], x, y, subFileThumbSize, subFileThumbSize, outputImage)

	x = subFilesThumbPadding
	y = subFilesThumbPadding*2 + subFileThumbSize
	drawSubfileThumbToOutput(subFilesThumbs[2], x, y, subFileThumbSize, subFileThumbSize, outputImage)

	x = subFilesThumbPadding*2 + subFileThumbSize
	drawSubfileThumbToOutput(subFilesThumbs[3], x, y, subFileThumbSize, subFileThumbSize, outputImage)

	var buf bytes.Buffer
	outputImage.EncodeJPG(&buf, nil)

	return buf.Bytes()
}

func drawSubfileThumbToOutput(subfileThumb []byte, x int, y int, width int, height int, outputImage *gg.Context) {
	if subfileThumb != nil {
		img, _, _ := image.Decode(bytes.NewReader(subfileThumb))
		img = imgconv.Resize(img, &imgconv.ResizeOption{Width: width, Height: height})
		outputImage.DrawImage(img, x, y)
	}
}

func (file *webFile) getSubfilesThubms(config config) [4][]byte {
	var subFilesThumbs [4][]byte

	path := config.RootPath + file.Path
	osFiles, _ := os.ReadDir(path)
	i := 0
	for _, osFile := range osFiles {
		var subFile webFile
		subFile.fillFields(osFile, path, config)
		if subFile.Type == "image" {
			subFilesThumbs[i] = subFile.getImageThumb(config)
			i++
		}

		if i > 3 {
			break
		}
	}

	return subFilesThumbs
}

func getFileType(file os.DirEntry, config config) string {
	if file.IsDir() {
		return "folder"
	}

	ext := getExtensionFromFileName(file.Name())

	for _, val := range config.ImageExt {
		if ext == val {
			return "image"
		}
	}

	for _, val := range config.VideoExt {
		if ext == val {
			return "video"
		}
	}

	return "other"
}

func getExtensionFromFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	if len(ext) > 0 {
		return strings.ToLower(ext)[1:]
	} else {
		return ""
	}
}
