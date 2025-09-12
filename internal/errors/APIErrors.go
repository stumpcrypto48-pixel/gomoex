package errors

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppErrorCode string

const (
	ErrCodeValidation AppErrorCode = "validation_failed"
	ErrCodeDecode     AppErrorCode = "decode_error"
	ErrCodeInternal   AppErrorCode = "internal_error"
	ErrURLParse       AppErrorCode = "url_parse_error"
	OkEndRequest      AppErrorCode = "ok_end_request"
)

type ApiOk struct {
	Code       AppErrorCode `json:"code"`
	Message    string       `json:"message"`
	Details    any          `json:"details,omitempty"`
	HTTPStatus int          `json:"-"`
}

func (ok *ApiOk) Error() string {
	return string(ok.Code)
}

type APIError struct {
	Code       AppErrorCode `json:"code"`
	Message    string       `json:"message,omitempty"`
	Details    any          `json:"details,omitempty"`
	HTTPStatus int          `json:"-"`
}

func (e *APIError) Error() string {
	return string(e.Code)
}

func CreateErrCodeValidation(details string) *APIError {
	return &APIError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		Details:    details,
		HTTPStatus: 400,
	}
}

func CreateErrCodeDecode(details string) *APIError {
	return &APIError{
		Code:       ErrCodeDecode,
		Message:    "Decode failed",
		Details:    details,
		HTTPStatus: 400,
	}
}

func CreateErrCodeInternal(details string) *APIError {
	return &APIError{
		Code:       ErrCodeInternal,
		Message:    "Internal error",
		Details:    details,
		HTTPStatus: 500,
	}
}

func WriteAPIError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *APIError:
		c.JSON(e.HTTPStatus, e)
	case *json.SyntaxError, *json.UnmarshalTypeError:
		c.JSON(http.StatusBadRequest, CreateErrCodeDecode(e.Error()))
	case validator.ValidationErrors:
		c.JSON(http.StatusBadRequest, CreateErrCodeValidation(e.Error()))
	// case *ApiOk:
	// 	c.JSON(http.StatusOK, e)
	default:
		c.JSON(http.StatusInternalServerError, CreateErrCodeInternal(e.Error()))
	}
}
