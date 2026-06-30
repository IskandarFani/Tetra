package dto

type CreateFileInput struct {
	UserID       uint
	OriginalName string
	StorageName  string
	MimeType     string
	Size         int64
}
