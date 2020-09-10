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
	sqlSelect           = `select user_id, username, role, email, sowalgasy from users`
	sqlSelectStores     = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores`
	sqlSelectChildStore = `select store_id, name, jemi_hasap_tmt, jemi_hasap_usd, shahsy_hasap_tmt, shahsy_hasap_usd from stores where parent_store_id = $1`
)

type Users struct {
	Userid    int    `json:"userid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Sowalgasy int    `json:"sowalgasy"`
}

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

	List := make([]*Users, 0)
	for rows.Next() {
		user := &Users{}
		err = rows.Scan(&user.Userid, &user.Username, &user.Role, &user.Email, &user.Sowalgasy)
		if err != nil {
			fmt.Println("ERROR")
			os.Exit(1101)
		}

		List = append(List, user)

	}

	item := List

	responses.SendResponse(w, err, item, nil)
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
	Storeid    int     `json:"store_id"`
	Name       string  `json:"store_name"`
	OverallTmt int     `json:"overall_tmt"`
	OverallUsd int     `json:"overall_usd"`
	OwnTmt     int     `json:"own_tmt"`
	OwnUsd     int     `json:"own_usd"`
	Child      Stores1 `json:"child"`
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

	ListofStores := make([]*Stores, 0)
	for rows.Next() {
		store := &Stores{}
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
		defer rows1.Close()

		for rows1.Next() {
			err = rows1.Scan(&store.Child.Storeid, &store.Child.Name, &store.Child.OverallTmt, &store.Child.OverallUsd, &store.Child.OwnTmt, &store.Child.OwnUsd)
			if err != nil {
				fmt.Println("ERROR")
				os.Exit(1101)
			}
		}

		ListofStores = append(ListofStores, store)

	}

	item := ListofStores

	responses.SendResponse(w, err, item, nil)
}
