package dto

type CreateWebhookRequestDto struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Url         string   `json:"url"`
	Events      []string `json:"events"`
}
