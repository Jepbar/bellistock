package filter

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"stock/config"
	"stock/function"
	"stock/responses"
	"time"

	"github.com/jackc/pgx"
)

func FilterWorkers(w http.ResponseWriter, r *http.Request) {
	Name := r.FormValue("name")
	Wezipesi := r.FormValue("wezipesi")
	DependingStore := r.FormValue("depending_store")

	var filter responses.Filterworkers
	filter.Name = Name
	filter.Wezipesi = Wezipesi
	filter.DependingStore = DependingStore

	sqlFilterWorkers, err := function.GenerateSqlFilterWorkers(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlFilterWorkers)
	defer rows.Close()

	List := make([]*responses.Workers, 0)
	for rows.Next() {
		worker := &responses.Workers{}
		err = rows.Scan(&worker.Workerid, &worker.Fullname, &worker.Wezipesi, &worker.Salary, &worker.DegisliDukany)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		List = append(List, worker)

	}

	item := List

	responses.SendResponse(w, err, item, nil)
}

func FilterMoneyTransfers(w http.ResponseWriter, r *http.Request) {
	Store := r.FormValue("store")
	TypeOfaccount := r.FormValue("type_of_account")
	Keyword := r.FormValue("keyword")
	Categorie := r.FormValue("categorie")
	Begin := r.FormValue("begin")
	End := r.FormValue("end")

	var filter responses.FilterMoneyTransfers
	filter.Store = Store
	filter.TypeOfaccount = TypeOfaccount
	filter.Keyword = Keyword
	filter.Categorie = Categorie
	filter.Begin = Begin
	filter.End = End

	sqlFilterMoneyTransfers, err := function.GenerateSqlFilterMoneyTransfers(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlFilterMoneyTransfers)
	defer rows.Close()

	List := make([]*responses.MoneyTransfer, 0)
	for rows.Next() {
		var Userid int
		var date time.Time
		moneyTransfer := &responses.MoneyTransfer{}
		err = rows.Scan(&moneyTransfer.Id, &moneyTransfer.Store, &moneyTransfer.TypeOfTransfer, &Userid, &moneyTransfer.TypeOfAccount, &moneyTransfer.TotalPaymentAmount, &moneyTransfer.Currency, &date, &moneyTransfer.Categorie)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}
		dateOfTransfer := date.Format("2006-01-02")
		Username := function.SelectUsername(Userid)

		moneyTransfer.DoneBy = Username
		moneyTransfer.Date = dateOfTransfer

		List = append(List, moneyTransfer)

	}

	item := List

	responses.SendResponse(w, err, item, nil)

}

func FilterIncomes(w http.ResponseWriter, r *http.Request) {
	Store := r.FormValue("store")
	TypeOfaccount := r.FormValue("type_of_account")
	Keyword := r.FormValue("keyword")
	Categorie := r.FormValue("categorie")
	Customer := r.FormValue("customer")
	TypeOfIncomePayment := r.FormValue("type_of_income_payment")
	Begin := r.FormValue("begin")
	End := r.FormValue("end")

	var filter responses.FilterIncomes
	filter.Store = Store
	filter.TypeOfaccount = TypeOfaccount
	filter.Keyword = Keyword
	filter.Categorie = Categorie
	filter.Customer = Customer
	filter.TypeOfIncomePayment = TypeOfIncomePayment
	filter.Begin = Begin
	filter.End = End

	sqlFilterIncome, err := function.GenerateSqlFilterIncomes(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlFilterIncome)
	defer rows.Close()

	List := make([]*responses.Incomes, 0)
	for rows.Next() {
		var date time.Time
		income := &responses.Incomes{}
		err = rows.Scan(&income.Id, &income.Store, &income.Customer, &income.Project, &income.TypeOfAccount, &income.TotalPaymentAmount, &income.Currency, &date, &income.Categorie)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}
		dateOfTransfer := date.Format("2006-01-02")

		income.Date = dateOfTransfer

		List = append(List, income)

	}

	item := List

	responses.SendResponse(w, err, item, nil)

}

func FilterOutcomes(w http.ResponseWriter, r *http.Request) {
	Store := r.FormValue("store")
	TypeOfaccount := r.FormValue("type_of_account")
	Keyword := r.FormValue("keyword")
	Categorie := r.FormValue("categorie")
	MoneyGoneTo := r.FormValue("money_gone_to")
	Begin := r.FormValue("begin")
	End := r.FormValue("end")

	var filter responses.FilterOutcomes
	filter.Store = Store
	filter.TypeOfaccount = TypeOfaccount
	filter.Keyword = Keyword
	filter.Categorie = Categorie
	filter.MoneyGoneTo = MoneyGoneTo
	filter.Begin = Begin
	filter.End = End

	sqlFilterOutcome, err := function.GenerateSqlFilterOutcomes(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlFilterOutcome)
	defer rows.Close()

	List := make([]*responses.Outcomes, 0)
	for rows.Next() {
		var date time.Time

		outcome := &responses.Outcomes{}
		err = rows.Scan(&outcome.Id, &outcome.Store, &outcome.MoneyGoneTo, &outcome.TotalPaymentAmount, &outcome.Currency, &outcome.TypeOfAccount, &date, &outcome.Categorie)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}
		dateOfTransfer := date.Format("2006-01-02")

		outcome.Date = dateOfTransfer

		List = append(List, outcome)

	}

	item := List

	responses.SendResponse(w, err, item, nil)

}

func FilterBetweenStores(w http.ResponseWriter, r *http.Request) {
	FromStoreName := r.FormValue("from_store_name")
	ToStoreName := r.FormValue("to_store_name")
	TypeOfAccount := r.FormValue("type_of_account")
	Begin := r.FormValue("begin")
	End := r.FormValue("end")

	var filter responses.FilterBetweenStores

	filter.FromStoreName = FromStoreName
	filter.ToStoreName = ToStoreName
	filter.TypeOfAccount = TypeOfAccount
	filter.Begin = Begin
	filter.End = End

	sqlFilterBetweenStores, err := function.GenerateSqlFilterTransferBetweenStores(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlFilterBetweenStores)

	defer rows.Close()

	List := make([]*responses.BetweenStores, 0)
	for rows.Next() {
		var date time.Time
		transfer := &responses.BetweenStores{}
		err = rows.Scan(&transfer.Id, &transfer.UserID, &transfer.FromStoreName, &transfer.ToStoreName, &transfer.TotalPaymentAmount, &transfer.Currency, &transfer.TypeOfAccount, &date)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}
		dateOfTransfer := date.Format("2006-01-02")

		transfer.Date = dateOfTransfer

		List = append(List, transfer)

	}

	item := List

	responses.SendResponse(w, err, item, nil)

}
