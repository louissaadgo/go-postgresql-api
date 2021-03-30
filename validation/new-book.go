package validation

import "github.com/go-playground/validator"

//NewBook is the data model for creating a new book
type NewBook struct {
	Title       string `json:"title" validate:"required,min=1,max=40"`
	AuthorID    int    `json:"authorID" validate:"required"`
	Description string `json:"description" validate:"required,min=40,max=1000"`
}

func (book NewBook) Validate() []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(book)
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
