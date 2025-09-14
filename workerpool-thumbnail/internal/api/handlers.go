package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"thumbnail/internal/services/thumbnail"
	"thumbnail/internal/services/worker"
	"time"
)

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxSize)
	if err := r.ParseMultipartForm(maxSize); err != nil {
		http.Error(w, "File is too big", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "File not selected", http.StatusBadRequest)
		return
	}

	var skippedResponse []SkippedResponse

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}

		if reason, ok := ValidateFile(fileHeader); !ok {
			skippedResponse = append(skippedResponse, *reason)
			file.Close()
			continue
		}

		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			continue
		}

		job := thumbnail.ThumbnailJob{
			FileName:  fileHeader.Filename,
			ImageData: data,
		}

		s.thumbnailService.Pool.Submit(worker.NewJob(job, s.thumbnailService))
	}

	response := &Response{
		Status:    "Accepted",
		Message:   "Accepted to processing",
		FileCount: len(files),
	}

	json, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(json)
}

func ValidateFile(fileHeader *multipart.FileHeader) (*SkippedResponse, bool) {
	extensions := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !validExts[extensions] {
		return &SkippedResponse{
			FileName: fileHeader.Filename,
			Reason:   "Invalid file extension",
		}, false
	}

	return nil, true
}

func (s *Server) defaultHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Works! %s", time.Now().Format(time.RFC850))))
}
