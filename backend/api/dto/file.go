package dto

type FileQueryDto struct {
	FolderID string `json:"folder_id"`
	Search   string `json:"search"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
}

type CreateFileRequestDto struct {
	Name        string `json:"name"`
	FolderID    string `json:"folder_id,omitempty"`
	Description string `json:"description,omitempty"`
}
