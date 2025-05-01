package backend

type FileMetadata struct {
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
}