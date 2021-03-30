package validation

import "github.com/go-playground/validator"

//OldAuthor is the data model for signing in an author
type OldAuthor struct {
	Email    string `json:"email" validate:"required,email,min=7,max=200"`
	Password string `json:"password" validate:"required,min=8,max=256"`
}

func (author OldAuthor) Validate() []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(author)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
