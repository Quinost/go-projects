package thumbnail

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	"thumbnail/internal/services/worker"
)

const size = 100

type ThumbnailJob struct {
	FileName  string `json:"file_name"`
	ImageData []byte `json:"-"`
}

type ThumbnailJobResult struct {
	FileName string `json:"file_name"`
	Success  bool   `json:"success"`
}

type ThumbnailService struct {
	Pool       *worker.Pool[ThumbnailJob]
	ResultChan chan ThumbnailJobResult
}

func NewThumbnailService(pool *worker.Pool[ThumbnailJob]) *ThumbnailService {
	return &ThumbnailService{
		Pool: pool,
		ResultChan: make(chan ThumbnailJobResult, 100),
	}
}

func (th *ThumbnailService) Process(data ThumbnailJob) error {
	log.Println("Processing", data.FileName)
	img, _, err := image.Decode(bytes.NewReader(data.ImageData))
	if err != nil {
		return fmt.Errorf("Error while decoding (%s): %v \n", data.FileName, err)
	}

	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	var newWidth, newHeight int
	if originalWidth > originalHeight {
		newWidth = size
		newHeight = (originalHeight * size) / originalWidth
	} else {
		newHeight = size
		newWidth = (originalWidth * size) / originalHeight
	}

	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, bounds, draw.Over, nil)

	nameWithoutExt := strings.TrimSuffix(data.FileName, filepath.Ext(data.FileName))
	fileName := fmt.Sprintf("%s_thumb_%dpx%s", nameWithoutExt, size, ".jpeg")
	envPath, _ := os.Getwd()
	dirPath := filepath.Join(envPath, "/images/thumb")
	ensureDir(dirPath)
	filePath := filepath.Join(dirPath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
	}
	defer file.Close()

	err = jpeg.Encode(file, thumbnail, &jpeg.Options{Quality: 85})

	if err != nil {
		return fmt.Errorf("Error while decoding (%s): %v \n", data.FileName, err)
	}

	log.Println("Processed", data.FileName)

	th.ResultChan <- ThumbnailJobResult{
		FileName: data.FileName,
		Success:  true,
	}

	return nil
}

func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}
