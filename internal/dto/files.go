package dto

import "time"

type CreateFileInput struct {
	UserID       uint
	FolderID     *uint
	OriginalName string
	StorageName  string
	MimeType     string
	Size         int64
}

type CreateFolderInput struct {
	UserID   uint   `gorm:"not null;index"`
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name"`
}

type UpdateFolderRequest struct {
	Name string `json:"name"`
}

type FolderResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateFolderRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type CreateFolderResponse struct {
	Status string         `json:"status"`
	Folder FolderResponse `json:"folder"`
}

type FileResponse struct {
	ID           uint      `json:"id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"created_at"`
}

type BreadcrumbItem struct {
	ID   *uint  `json:"id"`
	Name string `json:"name"`
}

type FolderContentResponse struct {
	Status        string           `json:"status"`
	CurrentFolder *FolderResponse  `json:"current_folder"`
	Breadcrumbs   []BreadcrumbItem `json:"breadcrumbs"`
	Folders       []FolderResponse `json:"folders"`
	Files         []FileResponse   `json:"files"`
}
