package helpers

import "github.com/go-playground/validator/v10"

var customMessages = map[string]string{
	"required": "tidak boleh kosong",
	"email":    "harus berupa alamat email yang valid",
	"min":      "terlalu pendek",
	"max":      "terlalu panjang",
	"numeric":  "harus berupa angka",
	"gte":      "harus lebih besar atau sama dengan %s",
	"lte":      "harus lebih kecil atau sama dengan %s",
	"integer":  "harus berupa bilangan bulat",
}

func FormatValidationError(err error) map[string]string {
	validationErrors := make(map[string]string)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			if customMessage, exists := customMessages[tag]; exists {
				validationErrors[field] = customMessage
			} else {
				validationErrors[field] = e.Error()
			}

			if tag == "gte" || tag == "lte" {
				validationErrors[field] = customMessages[tag] + " " + param
			}
		}
	}

	return validationErrors
}
