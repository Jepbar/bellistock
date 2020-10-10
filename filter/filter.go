package filter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

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
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var filter responses.Filterworkers
	filter.Name = Name
	filter.Wezipesi = Wezipesi
	filter.DependingStore = DependingStore

	sqlFilterWorkers, sqlNumberOfWorkersInFilter, err := function.GenerateSqlFilterWorkers(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var NumberOfWorkers int
	err2 := conn.QueryRow(context.Background(), sqlNumberOfWorkersInFilter).Scan(&NumberOfWorkers)
	if err1 != nil {
		fmt.Println(err2)
	}

	rows, err := conn.Query(context.Background(), sqlFilterWorkers, intLimit, intOffset)
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

	allWorkers := &responses.AllWorkers{}

	allWorkers.NumberOfWorkers = NumberOfWorkers
	allWorkers.List = List

	item := allWorkers

	responses.SendResponse(w, err, item, nil)
}

func FilterMoneyTransfers(w http.ResponseWriter, r *http.Request) {
	Store := r.FormValue("store")
	TypeOfaccount := r.FormValue("type_of_account")
	Keyword := r.FormValue("keyword")
	Categorie := r.FormValue("categorie")
	Begin := r.FormValue("begin")
	End := r.FormValue("end")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var filter responses.FilterMoneyTransfers
	filter.Store = Store
	filter.TypeOfaccount = TypeOfaccount
	filter.Keyword = Keyword
	filter.Categorie = Categorie
	filter.Begin = Begin
	filter.End = End

	sqlFilterMoneyTransfers, sqlNumberOfTransfers, err := function.GenerateSqlFilterMoneyTransfers(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var NumberOfTransfers int
	err2 := conn.QueryRow(context.Background(), sqlNumberOfTransfers).Scan(&NumberOfTransfers)
	if err2 != nil {
		fmt.Println(err2)
		fmt.Println("HAHAH")
	}

	rows, err := conn.Query(context.Background(), sqlFilterMoneyTransfers, intLimit, intOffset)
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
	alltransfers := &responses.AllTransfers{}
	alltransfers.NumberOfTransfers = NumberOfTransfers
	alltransfers.List = List

	item := alltransfers
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
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var filter responses.FilterIncomes
	filter.Store = Store
	filter.TypeOfaccount = TypeOfaccount
	filter.Keyword = Keyword
	filter.Categorie = Categorie
	filter.Customer = Customer
	filter.TypeOfIncomePayment = TypeOfIncomePayment
	filter.Begin = Begin
	filter.End = End

	sqlFilterIncome, sqlNumberOfIncomes, err := function.GenerateSqlFilterIncomes(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())
	var NumberOfIncomes int
	err2 := conn.QueryRow(context.Background(), sqlNumberOfIncomes).Scan(&NumberOfIncomes)
	if err2 != nil {
		fmt.Println(err2)
	}

	rows, err := conn.Query(context.Background(), sqlFilterIncome, intLimit, intOffset)
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
	total := &responses.TotalIncome{}

	total.TotalIncomeTmt = 0
	total.TotalIncomeUsd = 0
	total.NumberOfIncomes = NumberOfIncomes

	total.List = List

	item := total

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
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var filter responses.FilterOutcomes
	filter.Store = Store
	filter.TypeOfaccount = TypeOfaccount
	filter.Keyword = Keyword
	filter.Categorie = Categorie
	filter.MoneyGoneTo = MoneyGoneTo
	filter.Begin = Begin
	filter.End = End

	sqlFilterOutcome, sqlNumberOfOutcomes, err := function.GenerateSqlFilterOutcomes(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var NumberOfOutcomes int
	err2 := conn.QueryRow(context.Background(), sqlNumberOfOutcomes).Scan(&NumberOfOutcomes)
	if err2 != nil {
		fmt.Println(err2)
	}

	rows, err := conn.Query(context.Background(), sqlFilterOutcome, intLimit, intOffset)
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

	total := &responses.TotalOutcome{}

	total.TotalOutcomeTmt = 0
	total.TotalOutcomeUsd = 0
	total.NumberOfOutcomes = NumberOfOutcomes

	total.List = List

	item := total

	responses.SendResponse(w, err, item, nil)

}

func FilterBetweenStores(w http.ResponseWriter, r *http.Request) {
	FromStoreName := r.FormValue("from_store_name")
	ToStoreName := r.FormValue("to_store_name")
	TypeOfAccount := r.FormValue("type_of_account")
	Begin := r.FormValue("begin")
	End := r.FormValue("end")
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	var filter responses.FilterBetweenStores

	filter.FromStoreName = FromStoreName
	filter.ToStoreName = ToStoreName
	filter.TypeOfAccount = TypeOfAccount
	filter.Begin = Begin
	filter.End = End

	sqlFilterBetweenStores, sqlNumberOfTransfers, err := function.GenerateSqlFilterTransferBetweenStores(filter)

	token := function.ExtractToken(r)
	_, err1 := function.VerifyAccessToken(token)
	if err1 != nil {
		err1 = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var NumberOfTransfers int
	err2 := conn.QueryRow(context.Background(), sqlNumberOfTransfers).Scan(&NumberOfTransfers)
	if err2 != nil {
		fmt.Println(err2)
	}

	rows, err := conn.Query(context.Background(), sqlFilterBetweenStores, intLimit, intOffset)

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

	allTransfersBetweenStores := &responses.AlltransfersBetweenStores{}

	allTransfersBetweenStores.NumberOfTransfers = NumberOfTransfers
	allTransfersBetweenStores.List = List

	item := allTransfersBetweenStores

	responses.SendResponse(w, err, item, nil)

}
