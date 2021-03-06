package givingresponse

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
	sqlSelectUsers                 = `select user_id, username, role, email, sowalga_tmt, sowalga_usd from users where is_it_deleted = 'False' order by username limit $1 offset $2`
	sqlSelectStore                 = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores where store_id = $1 and is_it_deleted = 'False'`
	sqlSelectAllStore              = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores where is_it_deleted = 'False'`
	sqlSelectChildStores           = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores where parent_store_id = $1 and is_it_deleted = 'False'`
	sqlSelectLastActions           = `select l.id, u.username, l.action, l.message, l.create_ts, l.is_it_seen from last_modifications l inner join users u on l.user_id = u.user_id order by id desc limit $1 offset $2`
	sqlUpdateActions               = `update last_modifications set is_it_seen = $1 where id = $2`
	sqlSelectCustomer              = `select customer_id, name, girdeyjisi_tmt, girdeyjisi_usd from customers where is_it_deleted = 'False'`
	sqlSelectTransferBetweenStores = `select id, user_id, from_store_name, to_store_name, total_payment_amount, currency, type_of_account, date from transfers_between_stores order by create_ts desc limit $1 offset $2`
	sqlSelectWorkers               = `select worker_id , fullname, wezipesi, salary, degisli_dukany from workers where is_it_deleted = 'False' order by fullname limit $1 offset $2`
	sqlSelectMoneyTransfers        = `select m.id, s.name, m.type_of_transfer, m.user_id, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id order by m.create_ts desc limit $1 offset $2`
	sqlSelectIncomes               = `select m.id, s.name, m.customer, m.project, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id where m.type_of_transfer = 'girdi' order by create_ts desc limit $1 offset $2`
	sqlSelectOutcomes              = `select m.id, s.name, m.money_gone_to, m.total_payment_amount, m.currency, m.type_of_account, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id where m.type_of_transfer = 'cykdy' order by create_ts desc limit $1 offset $2`
	sqlSelectCategories            = `select categorie_id, name, parent_categorie from categories where is_it_deleted = 'False'`
	sqlSelectTotalIncome           = `select total_income_tmt , total_income_usd from income_outcome where id = 1`
	sqlSelectTotalOutcome          = `select total_outcome_tmt , total_outcome_usd from income_outcome where id = 1`
	sqlCountWorkers                = `select count(*) from workers where is_it_deleted = 'False'`
	sqlCountActions                = `select count(*) from last_modifications`
	sqlCountMoneyTransfer          = `select count(*) from  money_transfers`
	sqlCountIncomes                = `select count(*) from money_transfers where type_of_transfer = 'girdi'`
	sqlCountOutcomes               = `select count(*) from money_transfers where type_of_transfer = 'cykdy'`
	sqlCountTransfersBetweenStores = `select count(*) from transfers_between_stores`
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectUsers, intLimit, intOffset)

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

	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	Id := r.FormValue("id")
	k := 0
	if len(Id) > 0 {
		intId, _ := strconv.Atoi(Id)
		k = intId
	} else {
		k = 1
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), sqlSelectStore, k)

	defer rows.Close()

	ListofStores := make([]*responses.Stores1, 0)
	for rows.Next() {
		store := &responses.Stores1{}
		err = rows.Scan(&store.Storeid, &store.Name, &store.OverallTmt, &store.OverallUsd, &store.OwnTmt, &store.OwnUsd)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		rows1, err1 := conn.Query(context.Background(), sqlSelectChildStores, k)
		if err1 != nil {
			fmt.Println("ERRORRRRR")
		}

		ListofChilds := make([]*responses.Stores, 0)

		for rows1.Next() {
			child := &responses.Stores{}
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

func GetAllStores(w http.ResponseWriter, r *http.Request) {
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

	rows, err := conn.Query(context.Background(), sqlSelectAllStore)

	defer rows.Close()

	ListofAllStores := make([]*responses.Stores, 0)

	for rows.Next() {
		store := &responses.Stores{}
		err = rows.Scan(&store.Storeid, &store.Name, &store.OverallTmt, &store.OverallUsd, &store.OwnTmt, &store.OwnUsd)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}
		ListofAllStores = append(ListofAllStores, store)
	}

	item := ListofAllStores

	responses.SendResponse(w, err, item, nil)
}

func GetLastActions(w http.ResponseWriter, r *http.Request) {
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

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
	var NumberOfActions int
	err1 := conn.QueryRow(context.Background(), sqlCountActions).Scan(&NumberOfActions)
	if err1 != nil {
		fmt.Println(err1)
	}

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
			action.LastStatus = "Not seen"
		}
		if status == 1 {
			action.LastStatus = "Seen"
		}

		conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
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

	allactions := &responses.AllActions{}
	allactions.NumberOfActions = NumberOfActions
	allactions.List = ListofActions

	item := allactions

	responses.SendResponse(w, err, item, nil)
}

func GetCustomers(w http.ResponseWriter, r *http.Request) {
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
	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var NumberOfTransfers int
	err1 := conn.QueryRow(context.Background(), sqlCountTransfersBetweenStores).Scan(&NumberOfTransfers)
	if err1 != nil {
		fmt.Println(err1)
	}

	rows, err := conn.Query(context.Background(), sqlSelectTransferBetweenStores, intLimit, intOffset)

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

func GetWorkers(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var NumberOfWorkers int
	err1 := conn.QueryRow(context.Background(), sqlCountWorkers).Scan(&NumberOfWorkers)
	if err1 != nil {
		fmt.Println(err1)
		fmt.Println(NumberOfWorkers)
	}

	rows, err := conn.Query(context.Background(), sqlSelectWorkers, intLimit, intOffset)

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

func GetMoneyTransfers(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")

	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	var NumberOfTransfers int
	err1 := conn.QueryRow(context.Background(), sqlCountMoneyTransfer).Scan(&NumberOfTransfers)
	if err1 != nil {
		fmt.Println(err1)
	}

	rows, err := conn.Query(context.Background(), sqlSelectMoneyTransfers, intLimit, intOffset)
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

func GetIncomes(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var tmt int
	var usd int
	err1 := conn.QueryRow(context.Background(), sqlSelectTotalIncome).Scan(&tmt, &usd)
	if err1 != nil {
		fmt.Println("erroro")
	}

	var NumberOfIncomes int
	err2 := conn.QueryRow(context.Background(), sqlCountIncomes).Scan(&NumberOfIncomes)
	if err1 != nil {
		fmt.Println(err2)
	}

	rows, err := conn.Query(context.Background(), sqlSelectIncomes, intLimit, intOffset)
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

	total.TotalIncomeTmt = tmt
	total.TotalIncomeUsd = usd
	total.NumberOfIncomes = NumberOfIncomes

	total.List = List

	item := total

	responses.SendResponse(w, err, item, nil)
}

func GetOutcomes(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	_, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	limit := r.FormValue("limit")
	offset := r.FormValue("offset")
	intLimit, _ := strconv.Atoi(limit)
	intOffset, _ := strconv.Atoi(offset)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var tmt int
	var usd int
	err1 := conn.QueryRow(context.Background(), sqlSelectTotalOutcome).Scan(&tmt, &usd)
	if err1 != nil {
		fmt.Println("erroro")
	}

	var NumberOutcomes int
	err2 := conn.QueryRow(context.Background(), sqlCountOutcomes).Scan(&NumberOutcomes)
	if err2 != nil {
		fmt.Println(err2)
	}
	rows, err := conn.Query(context.Background(), sqlSelectOutcomes, intLimit, intOffset)
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

	total.TotalOutcomeTmt = tmt
	total.TotalOutcomeUsd = usd
	total.NumberOfOutcomes = NumberOutcomes

	total.List = List

	item := total

	responses.SendResponse(w, err, item, nil)
}
