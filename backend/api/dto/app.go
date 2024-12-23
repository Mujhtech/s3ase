package dto

type CreateAppRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Region      string `json:"region,omitempty"`
	Slug        string `json:"slug,omitempty"`
}
