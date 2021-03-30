package validation

//ErrorResponse is the response error model
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}
