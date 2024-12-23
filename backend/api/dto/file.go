package dto

type FileQueryDto struct {
	FolderID string `json:"folder_id"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
}
