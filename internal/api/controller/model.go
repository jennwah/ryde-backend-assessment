package controller

import "fmt"

// ErrorResponseModel represents the response payload for error response.
// swagger:model ErrorResponseModel
type ErrorResponseModel struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func InternalServerErrorResp() ErrorResponseModel {
	return ErrorResponseModel{
		ErrorCode:    "500",
		ErrorMessage: "Unexpected error encountered. Please try again later",
	}
}

func BadRequestErrorResp(err error) ErrorResponseModel {
	return ErrorResponseModel{
		ErrorCode:    "400",
		ErrorMessage: fmt.Sprintf("Bad request body: %s", err.Error()),
	}
}
