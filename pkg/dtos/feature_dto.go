package dtos

type FeatureRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type FeatureResponse struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
