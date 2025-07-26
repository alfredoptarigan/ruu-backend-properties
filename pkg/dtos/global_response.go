package dtos

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponseDTO struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Code    int         `json:"code,omitempty"` // opsional: kode error internal
	Errors  interface{} `json:"errors,omitempty"`
}
