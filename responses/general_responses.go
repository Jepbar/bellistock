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

type Users struct {
	Userid     int    `json:"userid"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	SowalgaTmt int    `json:"sowalga_tmt"`
	SowalgaUsd int    `json:"sowalga_usd"`
}

type Customers struct {
	Customerid    int    `json:"customerid"`
	Name          string `json:"name"`
	GirdeyjisiTmt int    `json:"girdeyjisi_tmt"`
	GirdeyjisiUsd int    `json:"girdeyjisi_usd"`
}

type Stores1 struct {
	Storeid    int    `json:"store_id"`
	Name       string `json:"store_name"`
	OverallTmt int    `json:"overall_tmt"`
	OverallUsd int    `json:"overall_usd"`
	OwnTmt     int    `json:"own_tmt"`
	OwnUsd     int    `json:"own_usd"`
}

type Stores struct {
	Storeid    int        `json:"store_id"`
	Name       string     `json:"store_name"`
	OverallTmt int        `json:"overall_tmt"`
	OverallUsd int        `json:"overall_usd"`
	OwnTmt     int        `json:"own_tmt"`
	OwnUsd     int        `json:"own_usd"`
	Childs     []*Stores1 `json:"childs"`
}

type LastActions struct {
	Id         int    `json:"id"`
	User       string `json:"done_by"`
	Action     string `json:"action"`
	Message    string `json:"message"`
	Date       string `json:"date"`
	LastStatus string `json:"last_status"`
}

type UserLogin struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Workers struct {
	Workerid      int    `json:"worker_id"`
	Fullname      string `json:"fullname"`
	Salary        int    `json:"salary"`
	Wezipesi      string `json:"wezipes"`
	DegisliDukany string `json:"degisli_dukany"`
}

type BetweenStores struct {
	Id                 int    `json:"id"`
	UserID             int    `json:"user_id"`
	FromStoreName      string `json:"from_store_name"`
	ToStoreName        string `json:"to_store_name"`
	TypeOfAccount      string `json:"type_of_account"`
	Currency           string `json:"currency"`
	TotalPaymentAmount int    `json:"total_payment_amount"`
	Date               string `json:"date"`
}

type MoneyTransfer struct {
	Id                 int    `json:"id"`
	Store              string `json:"store"`
	TypeOfTransfer     string `json:"type_of_transfer"`
	TypeOfAccount      string `json:"type_of_account"`
	DoneBy             string `json:"done_by"`
	TotalPaymentAmount int    `json:"total_payment_amount"`
	Currency           string `json:"currency"`
	Categorie          string `json:"categorie"`
	Date               string `json:"date"`
}

type Incomes struct {
	Id                 int    `json:"id"`
	Store              string `json:"store"`
	TypeOfAccount      string `json:"type_of_account"`
	TotalPaymentAmount int    `json:"total_payment_amount"`
	Currency           string `json:"currency"`
	Categorie          string `json:"categorie"`
	Date               string `json:"date"`
	Customer           string `json:"customer"`
	Project            string `json:"project"`
}

type Outcomes struct {
	Id                 int    `json:"id"`
	Store              string `json:"store"`
	MoneyGoneTo        string `json:"money_gone_to"`
	TypeOfAccount      string `json:"type_of_account"`
	TotalPaymentAmount int    `json:"total_payment_amount"`
	Currency           string `json:"currency"`
	Categorie          string `json:"categorie"`
	Date               string `json:"date"`
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
