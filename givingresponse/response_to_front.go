package givingresponse

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/function"
	"stock/responses"

	"github.com/jackc/pgx"
)

const (
	sqlSelect                      = `select user_id, username, role, email, sowalga_tmt, sowalga_usd from users`
	sqlSelectStores                = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores`
	sqlSelectChildStore            = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores where parent_store_id = $1`
	sqlSelectLastActions           = `select u.username, l.action, l.message from last_modifications l inner join users u on l.user_id = u.user_id where is_it_seen = $1`
	sqlUpdateActions               = `update last_modifications set is_it_seen = $1`
	sqlSelectCustomer              = `select customer_id, name, girdeyjisi_tmt, girdeyjisi_usd from customers`
	sqlSelectTransferBetweenStores = `select id, user_id, from_store_name, to_store_name, total_payment_amount, currency, type_of_account from transfers_between_stores`
)

func GetUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Header)
	fmt.Println(r.Method)

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

	rows, err := conn.Query(context.Background(), sqlSelectLastActions, 0)

	defer rows.Close()

	ListofActions := make([]*responses.LastActions, 0)
	for rows.Next() {
		action := &responses.LastActions{}
		err = rows.Scan(&action.User, &action.Action, &action.Message)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1111)
		}

		ListofActions = append(ListofActions, action)
	}

	rows1, err := conn.Exec(context.Background(), sqlUpdateActions, 1)
	if rows1 == nil {
		fmt.Println(rows1)
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
		transfer := &responses.BetweenStores{}
		err = rows.Scan(&transfer.Id, &transfer.UserID, &transfer.FromStoreName, &transfer.ToStoreName, &transfer.TotalPaymentAmount, &transfer.Currency, &transfer.TypeOfAccount)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		List = append(List, transfer)

	}

	item := List

	responses.SendResponse(w, err, item, nil)

}
