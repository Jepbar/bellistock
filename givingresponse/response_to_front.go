package givingresponse

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/function"
	"stock/responses"
	"strconv"
	"time"

	"github.com/jackc/pgx"
)

const (
	sqlSelect                      = `select user_id, username, role, email, sowalga_tmt, sowalga_usd from users where is_it_deleted = 'False'`
	sqlSelectStores                = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores`
	sqlSelectChildStore            = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores where parent_store_id = $1`
	sqlSelectLastActions           = `select l.id, u.username, l.action, l.message, l.create_ts, l.is_it_seen from last_modifications l inner join users u on l.user_id = u.user_id order by id desc limit $1 offset $2`
	sqlUpdateActions               = `update last_modifications set is_it_seen = $1 where id = $2`
	sqlSelectCustomer              = `select customer_id, name, girdeyjisi_tmt, girdeyjisi_usd from customers where is_it_deleted = 'False'`
	sqlSelectTransferBetweenStores = `select id, user_id, from_store_name, to_store_name, total_payment_amount, currency, type_of_account, date from transfers_between_stores`
	sqlSelectWorkers               = `select worker_id , fullname, wezipesi, salary, degisli_dukany from workers`
	sqlSelectMoneyTransfers        = `select m.id, s.name, m.type_of_transfer, m.user_id, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id`
	sqlSelectIncomes               = `select m.id, s.name, m.customer, m.project, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id where m.type_of_transfer = 'girdi'`
	sqlSelectOutcomes              = `select m.id, s.name, m.money_gone_to, m.total_payment_amount, m.currency, m.type_of_account, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id where m.type_of_transfer = 'cykdy'`
	sqlSelectCategories            = `select categorie_id, name, parent_categorie from categories`
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelect)

	defer rows.Close()

	List := make([]*responses.Users, 0)
	for rows.Next() {
		user := &responses.Users{}
		err = rows.Scan(&user.Userid, &user.Username, &user.Role, &user.Email, &user.SowalgaTmt, &user.SowalgaUsd)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		List = append(List, user)

	}

	item := List

	responses.SendResponse(w, err, item, nil)

}

func GetStores(w http.ResponseWriter, r *http.Request) {

	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectStores)

	defer rows.Close()

	ListofStores := make([]*responses.Stores, 0)
	for rows.Next() {
		store := &responses.Stores{}
		err = rows.Scan(&store.Storeid, &store.Name, &store.OverallTmt, &store.OverallUsd, &store.OwnTmt, &store.OwnUsd)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		rows1, err1 := conn.Query(context.Background(), sqlSelectChildStore, store.Storeid)
		if err1 != nil {
			fmt.Println("ERRORRRRR")
		}

		ListofChilds := make([]*responses.Stores1, 0)

		for rows1.Next() {
			child := &responses.Stores1{}
			err = rows1.Scan(&child.Storeid, &child.Name, &child.OverallTmt, &child.OverallUsd, &child.OwnTmt, &child.OwnUsd)
			if err != nil {
				fmt.Println("ERROR")
				os.Exit(1101)
			}
			ListofChilds = append(ListofChilds, child)
		}

		store.Childs = ListofChilds
		ListofStores = append(ListofStores, store)
	}

	item := ListofStores

	responses.SendResponse(w, err, item, nil)
}

func GetLastActions(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectLastActions, intLimit, intOffset)

	defer rows.Close()

	ListofActions := make([]*responses.LastActions, 0)
	for rows.Next() {
		var date time.Time
		var status int
		var id int
		action := &responses.LastActions{}
		err = rows.Scan(&id, &action.User, &action.Action, &action.Message, &date, &status)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1111)
		}

		dateOfTransfer := date.Format("2006-01-02")
		action.Date = dateOfTransfer

		if status == 0 {
			action.LastStatus = "It has not seen yet"
		}
		if status == 1 {
			action.LastStatus = "It has already seen"
		}

		conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		rows1, err := conn.Exec(context.Background(), sqlUpdateActions, 1, id)
		if err == nil {
			fmt.Println(rows1)
		}

		action.Id = id

		ListofActions = append(ListofActions, action)

	}

	item := ListofActions

	responses.SendResponse(w, err, item, nil)
}

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectCustomer)

	defer rows.Close()

	List := make([]*responses.Customers, 0)
	for rows.Next() {
		customer := &responses.Customers{}
		err = rows.Scan(&customer.Customerid, &customer.Name, &customer.GirdeyjisiTmt, &customer.GirdeyjisiUsd)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		List = append(List, customer)

	}
	item := List

	responses.SendResponse(w, err, item, nil)

}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectCategories)

	defer rows.Close()
	List := make([]*responses.Categories, 0)
	for rows.Next() {
		categorie := &responses.Categories{}
		err = rows.Scan(&categorie.Categorieid, &categorie.Name, &categorie.ParentCategorie)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		List = append(List, categorie)

	}
	item := List

	responses.SendResponse(w, err, item, nil)
}

func GetTransferBetweenStores(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectTransferBetweenStores)

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

func GetWorkers(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectWorkers)

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

func GetMoneyTransfers(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectMoneyTransfers)
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

func GetIncomes(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectIncomes)
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

func GetOutcomes(w http.ResponseWriter, r *http.Request) {
	error2 := function.TokenValid(r)
	if error2 != nil {
		fmt.Println("Time is over!")
		os.Exit(112)
	}

	_, error1 := function.VerifyToken(r)
	if error1 != nil {
		fmt.Println("It is not mine!")
		os.Exit(111)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectOutcomes)
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
