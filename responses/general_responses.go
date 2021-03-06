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

type UserDataForEditing struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type Customers struct {
	Customerid    int    `json:"customerid"`
	Name          string `json:"name"`
	GirdeyjisiTmt int    `json:"girdeyjisi_tmt"`
	GirdeyjisiUsd int    `json:"girdeyjisi_usd"`
}

type CustomerDataForEditing struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

type Categories struct {
	Categorieid     int    `json:"categorieid"`
	Name            string `json:"name"`
	ParentCategorie string `json:"parent_categorie"`
}

type CategorieDataForEditing struct {
	Name string `json:"name"`
}

type Stores struct {
	Storeid    int    `json:"store_id"`
	Name       string `json:"store_name"`
	OverallTmt int    `json:"overall_tmt"`
	OverallUsd int    `json:"overall_usd"`
	OwnTmt     int    `json:"own_tmt"`
	OwnUsd     int    `json:"own_usd"`
}

type Stores1 struct {
	Storeid    int       `json:"store_id"`
	Name       string    `json:"store_name"`
	OverallTmt int       `json:"overall_tmt"`
	OverallUsd int       `json:"overall_usd"`
	OwnTmt     int       `json:"own_tmt"`
	OwnUsd     int       `json:"own_usd"`
	Childs     []*Stores `json:"childs"`
}

type AllActions struct {
	NumberOfActions int            `json:"number_of_actions"`
	List            []*LastActions `json:"list"`
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
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
}

type AllWorkers struct {
	NumberOfWorkers int        `json:"number_of_workers"`
	List            []*Workers `json:"list"`
}

type Workers struct {
	Workerid      int    `json:"worker_id"`
	Fullname      string `json:"fullname"`
	Salary        int    `json:"salary"`
	Wezipesi      string `json:"wezipesi"`
	DegisliDukany string `json:"degisli_dukany"`
}

type WorkerDataForEditing struct {
	Fullname      string `json:"fullname"`
	Salary        int    `json:"salary"`
	Wezipesi      string `json:"wezipesi"`
	DegisliDukany string `json:"degisli_dukany"`
	HomePhone     string `json:"home_phone"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	HomeAddres    string `json:"home_addres"`
	Note          string `json:"note"`
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

type AlltransfersBetweenStores struct {
	NumberOfTransfers int              `json:"number_of_transfers"`
	List              []*BetweenStores `json:"list"`
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

type AllTransfers struct {
	NumberOfTransfers int              `json:"number_of_transfers"`
	List              []*MoneyTransfer `json:"list"`
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

type IncomeDataForEditing struct {
	Store               string `json:"store"`
	TypeOfAccount       string `json:"type_of_account"`
	TotalPaymentAmount  int    `json:"total_payment_amount"`
	TypeOfTransfer      string `json:"type_of_transfer"`
	Currency            string `json:"currency"`
	Categorie           string `json:"categorie"`
	Date                string `json:"date"`
	Customer            string `json:"customer"`
	Project             string `json:"project"`
	Keyword             string `json:"keyword"`
	TypeOfIncomePayment string `json:"type_of_income_payment"`
	Note                string `json:"note"`
}

type TotalIncome struct {
	TotalIncomeTmt  int        `json:"total_income_tmt"`
	TotalIncomeUsd  int        `json:"total_income_usd"`
	NumberOfIncomes int        `json:"number_of_incomes"`
	List            []*Incomes `json:"list"`
}

type TotalOutcome struct {
	TotalOutcomeTmt  int         `json:"total_outcome_tmt"`
	TotalOutcomeUsd  int         `json:"total_outcome_usd"`
	NumberOfOutcomes int         `json:"number_of_outcomes"`
	List             []*Outcomes `json:"list"`
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

type OutcomeDataForEditing struct {
	Store              string `json:"store"`
	MoneyGoneTo        string `json:"money_gone_to"`
	TypeOfAccount      string `json:"type_of_account"`
	TypeOfTransfer     string `json:"type_of_transfer"`
	TotalPaymentAmount int    `json:"total_payment_amount"`
	Currency           string `json:"currency"`
	Categorie          string `json:"categorie"`
	Date               string `json:"date"`
	Keyword            string `json:"keyword"`
	Note               string `json:"note"`
}

type Filterworkers struct {
	Name           string `json:"name"`
	DependingStore string `json:"depending_store"`
	Wezipesi       string `json:"wezipesi"`
}

type FilterMoneyTransfers struct {
	Store         string `json:"store"`
	TypeOfaccount string `json:"type_of_account"`
	Keyword       string `json:"keyword"`
	Categorie     string `json:"categorie"`
	Begin         string `json:"begin"`
	End           string `json:"end"`
}

type FilterIncomes struct {
	Store               string `json:"store"`
	TypeOfaccount       string `json:"type_of_account"`
	Keyword             string `json:"keyword"`
	Categorie           string `json:"categorie"`
	Customer            string `json:"customer"`
	TypeOfIncomePayment string `json:"type_of_income_payment"`
	Begin               string `json:"begin"`
	End                 string `json:"end"`
}

type FilterOutcomes struct {
	Store         string `json:"store"`
	TypeOfaccount string `json:"type_of_account"`
	MoneyGoneTo   string `json:"money_gone_to"`
	Keyword       string `json:"keyword"`
	Categorie     string `json:"categorie"`
	Begin         string `json:"begin"`
	End           string `json:"end"`
}

type FilterBetweenStores struct {
	FromStoreName string `json:"from_store_name"`
	ToStoreName   string `json:"to_store_name"`
	TypeOfAccount string `json:"type_of_account"`
	Begin         string `json:"begin"`
	End           string `json:"end"`
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
