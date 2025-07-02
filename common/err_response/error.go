package err_response

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"todo_project/constant"
)

type CustomError struct {
	error
	Status   int      `json:"status"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	ErrorMsg string   `json:"error,omitempty"`
	Details  []string `json:"details,omitempty"`
}

func NewCustomError(status int, code string, message string, errs ...error) huma.StatusError {
	details := make([]string, len(errs))
	for i, err := range errs {
		details[i] = err.Error()
	}
	errMsg := message
	if len(details) > 0 {
		errMsg = fmt.Sprintf("%s", details[0])
	}
	return &CustomError{
		Status:   status,
		Code:     code,
		Message:  message,
		ErrorMsg: errMsg,
		Details:  details,
	}
}

func (e *CustomError) Error() string {
	if e.error != nil {
		return e.error.Error()
	}
	return e.Message
}

func (e *CustomError) GetStatus() int {
	return e.Status
}

func HandleError(err error) huma.StatusError {
	errorCode := "ERR_SYSTEM_ERROR"
	errorCodeMessage := http.StatusText(http.StatusInternalServerError)
	errorMessage := err.Error()
	status := http.StatusInternalServerError

	if code, ok := constant.MAP_ERROR_CODE[constant.ERROR_CODE(err.Error())]; ok {
		status = http.StatusOK
		errorCode = err.Error()
		errorCodeMessage = code
		errorMessage = fmt.Sprintf("%s: %s", err.Error(), code)
	}

	return &CustomError{
		error:    err,
		Status:   status,
		Message:  errorCodeMessage,
		Code:     errorCode,
		ErrorMsg: errorMessage,
	}
}

func NewHumaError() {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		details := make([]string, len(errs))
		for i, err := range errs {
			details[i] = err.Error()
		}
		code := string(constant.ERR_REQUEST_INVALID)
		if message == string(constant.ERR_UNAUTHORIZED) {
			code = string(constant.ERR_UNAUTHORIZED)
			message = "User not authorized"
		}
		return &CustomError{
			Status:  http.StatusOK,
			Code:    code,
			Message: message,
			Details: details,
		}
	}
}

func ErrBadRequest(message string, locs ...string) *CustomError {
	details := make([]string, len(locs))
	for i, loc := range locs {
		details[i] = loc
	}
	return &CustomError{
		Status:   http.StatusBadRequest,
		Message:  message,
		Code:     string(constant.ERR_REQUEST_INVALID),
		ErrorMsg: fmt.Sprintf("%s: %s", constant.ERR_REQUEST_INVALID, message),
		Details:  details,
	}
}

func ErrUnauthorized(err error, message string, details ...string) *CustomError {
	return &CustomError{
		error:    err,
		Status:   http.StatusUnauthorized,
		Message:  message,
		Code:     string(constant.ERR_UNAUTHORIZED),
		ErrorMsg: fmt.Sprintf("%s: %s", constant.ERR_UNAUTHORIZED, message),
		Details:  details,
	}
}

func ErrNotFound(message string, notFoundCode string, details ...string) *CustomError {
	return &CustomError{
		Status:   http.StatusNotFound,
		Message:  message,
		Code:     notFoundCode,
		ErrorMsg: message,
		Details:  details,
	}
}

func ErrInternalServerError(err error, message string, internalServerErrorCode string, errs ...error) *CustomError {
	var details []string
	if len(errs) > 0 {
		for _, e := range errs {
			details = append(details, e.Error())
		}
	}
	if err != nil {
		details = append(details, err.Error())
	}
	return &CustomError{
		error:    err,
		Status:   http.StatusInternalServerError,
		Message:  message,
		Code:     internalServerErrorCode,
		ErrorMsg: message,
		Details:  details,
	}
}
