package givedataforediting

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/config"
	"stock/function"
	"stock/responses"
	"strconv"
	"time"

	"github.com/jackc/pgx"
)

const (
	sqlGiveDataAboutUser            = `select username, email, role from users where user_id = $1`
	sqlGiveDataAboutWorker          = `select fullname, degisli_dukany, wezipesi, salary, phone, home_phone, home_addres, email, note from workers where worker_id = $1`
	sqlGiveDataAboutCustomer        = `select name, note from customers where customer_id = $1`
	sqlGiveDataAboutCategorie       = `select name from categories where categorie_id = $1`
	sqlGiveDataAboutIncomeTransfer  = `select s.name, m.customer, m.project, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie, m.type_of_income_payment, m.keyword, m.note from money_transfers m inner join stores s on s.store_id = m.store_id where m.id = $1`
	sqlGiveDataAboutOutcomeTransfer = `select s.name, m.money_gone_to, m.total_payment_amount, m.currency, m.type_of_account, m.date, m.categorie, m.keyword, m.note from money_transfers m inner join stores s on s.store_id = m.store_id where m.id = $1`
)

func GiveDataAboutUserForEditing(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	user := &responses.UserDataForEditing{}
	err = conn.QueryRow(context.Background(), sqlGiveDataAboutUser, IntId).Scan(&user.Username, &user.Email, &user.Role)
	if err != nil {
		fmt.Println("ERROR")
	}

	item := user

	responses.SendResponse(w, err, item, nil)

}

func GiveDataAboutWorkerForEditing(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	worker := &responses.WorkerDataForEditing{}

	err = conn.QueryRow(context.Background(), sqlGiveDataAboutWorker, IntId).Scan(&worker.Fullname, &worker.DegisliDukany, &worker.Wezipesi, &worker.Salary, &worker.Phone, &worker.HomePhone, &worker.HomeAddres, &worker.Email, &worker.Note)
	if err != nil {
		fmt.Println("ERROR")
	}

	item := worker

	responses.SendResponse(w, err, item, nil)
}

func GiveDataAboutCustomerForediting(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	customer := &responses.CustomerDataForEditing{}
	ok := false
	err1 := conn.QueryRow(context.Background(), sqlGiveDataAboutCustomer, IntId).Scan(&customer.Name, &customer.Note)
	if err1 != nil {
		fmt.Println(err1)
		ok = true
	}
	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	} else {
		item := customer

		responses.SendResponse(w, err, item, nil)
	}
}

func GiveDataAboutCategorieForEditing(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	ok := false
	categorie := &responses.CategorieDataForEditing{}
	err1 := conn.QueryRow(context.Background(), sqlGiveDataAboutCategorie, IntId).Scan(&categorie.Name)
	if err1 != nil {
		fmt.Println(err1)
		ok = true
	}

	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	} else {
		item := categorie

		responses.SendResponse(w, err, item, nil)
	}

}

func GiveDataAboutIncomeTransferForEditing(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	ok := false

	income := &responses.IncomeDataForEditing{}
	var date time.Time
	err1 := conn.QueryRow(context.Background(), sqlGiveDataAboutIncomeTransfer, IntId).Scan(&income.Store, &income.Customer, &income.Project, &income.TypeOfAccount, &income.TotalPaymentAmount, &income.Currency, &date, &income.Categorie, &income.TypeOfIncomePayment, &income.Keyword, &income.Note)
	if err1 != nil {
		fmt.Println(err1)
		ok = true
	}
	dateOfTransfer := date.Format("2006-01-02")

	income.Date = dateOfTransfer
	income.TypeOfTransfer = "girdi"
	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	} else {
		item := income

		responses.SendResponse(w, err, item, nil)
	}

}

func GiveDataAboutOutcomeTransfer(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	ok := false

	outcome := &responses.OutcomeDataForEditing{}

	var date time.Time
	err1 := conn.QueryRow(context.Background(), sqlGiveDataAboutOutcomeTransfer, IntId).Scan(&outcome.Store, &outcome.MoneyGoneTo, &outcome.TotalPaymentAmount, &outcome.Currency, &outcome.TypeOfAccount, &date, &outcome.Categorie, &outcome.Keyword, &outcome.Note)
	if err1 != nil {
		fmt.Println(err1)
		ok = true
	}
	dateOfTransfer := date.Format("2006-01-02")

	outcome.Date = dateOfTransfer
	outcome.TypeOfTransfer = "cykdy"
	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	} else {
		item := outcome

		responses.SendResponse(w, err, item, nil)
	}

}
