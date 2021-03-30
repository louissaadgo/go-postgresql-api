package validation

import "github.com/go-playground/validator"

//NewAuthor is the data model for creating an author
type NewAuthor struct {
	Email    string `json:"email" validate:"required,email,min=7,max=200"`
	Password string `json:"password" validate:"required,min=8,max=256"`
	Name     string `json:"name" validate:"required,min=1,max=25"`
	LastName string `json:"lastName" validate:"required,min=1,max=25"`
}

func (author NewAuthor) Validate() []*ErrorResponse {
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
