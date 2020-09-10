package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type GeneralResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status"`
	Data       interface{} `json:"data"`
}

type ResponseErrorCodeAndMessage struct {
	ErrorCode    string  `json:"error_code"`
	ErrorMessage *string `json:"error_msg,omitempty"`
}

const (
	ErrorCodeOK                  int = 200
	ErrorCodeTfaRequired         int = 250
	ErrorCodeBadRequest          int = 400
	ErrorCodeUnauthorized        int = 401
	ErrorCodeForbidden           int = 403
	ErrorCodeNotFound            int = 404
	ErrorCodeConflict            int = 409
	ErrorCodeExpired             int = 408
	ErrorCodeFileSizeTooLarge    int = 413
	ErrorCodeTooManyRequests     int = 429
	ErrorCodeInternalServerError int = 500
)

var ErrOK = errors.New("OK")
var ErrTfaRequired = errors.New("Tfa required")
var ErrBadRequest = errors.New("Bad request")
var ErrUnauthorized = errors.New("Unauthorized")
var ErrForbidden = errors.New("Forbidden")
var ErrNotFound = errors.New("Not found")
var ErrConflict = errors.New("Conflict")
var ErrExpired = errors.New("Expired")
var ErrFileSizeTooLarge = errors.New("File size too large")
var ErrTooManyRequests = errors.New("OTP retry limit exceeded")
var ErrInternalServerError = errors.New("Internal server error")

const (
	ErrorMessageOK                  = "ok"
	ErrorMessageTfaRequired         = "tfa_required"
	ErrorMessageBadRequest          = "bad_request"
	ErrorMessageUnauthorized        = "unauthorized"
	ErrorMessageForbidden           = "forbidden"
	ErrorMessageNotFound            = "not_found"
	ErrorMessageConflict            = "conflict"
	ErrorMessageExpired             = "expired"
	ErrorMessageFileSizeTooLarge    = "file_size_too_large"
	ErrorMessageTooManyRequests     = "otp_retry_limit_exceeded"
	ErrorMessageInternalServerError = "internal_server_error"
)

func SendResponse(w http.ResponseWriter, err error, data interface{}, clog *log.Entry) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var statusCode int = 200
	switch err {
	case ErrTfaRequired:
		statusCode = 250
	case ErrBadRequest:
		statusCode = 400
	case ErrUnauthorized:
		statusCode = 401
	case ErrForbidden:
		statusCode = 403
	case ErrNotFound:
		statusCode = 404
	case ErrConflict:
		statusCode = 409
	case ErrExpired:
		statusCode = 408
	case ErrFileSizeTooLarge:
		statusCode = 413
	case ErrTooManyRequests:
		statusCode = 429
	case ErrInternalServerError:
		statusCode = 500
	}

	w.WriteHeader(statusCode)

	var resp GeneralResponse
	if statusCode == ErrorCodeOK {
		resp.Success = true
		resp.StatusCode = statusCode
		resp.Data = data
	} else {
		resp.Success = false
		resp.StatusCode = statusCode
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		clog.WithError(err).Error(fmt.Sprint(" data: ", resp))
	}
}
