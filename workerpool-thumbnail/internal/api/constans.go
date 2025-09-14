package api

const maxSize = 20 << 20

type Response struct {
	Status          string            `JSON:"status"`
	Message         string            `JSON:"message"`
	SkippedResponse []SkippedResponse `JSON:"skipped,omitempty"`
	FileCount       int               `JSON:"file_count,omitempty"`
}

type SkippedResponse struct {
	FileName string `JSON:"file_name"`
	Reason   string `JSON:"reason"`
}

var validExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}
