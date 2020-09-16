package creations

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/function"
	"strconv"

	"github.com/jackc/pgx"
)

const (
	sqlInsertUser      = `insert into users(username, email, password, role) values($1, $2, $3, $4) returning user_id`
	sqlInsertStore     = `insert into stores(name, parent_store_id) values($1, $2) returning store_id`
	sqlInsertMessage1  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyny doretdi') returning id`
	sqlInsertCustomer  = `insert into customers(name, note) values($1, $2) returning customer_id `
	sqlInsertMessage2  = `insert into last_modifications(user_id, action, message) values($1, $2, $3  || ' atly musderini sanawa gosdy') returning id`
	sqlInsertMessage3  = `insert into last_modifications(user_id, action,message) values($1, $2, $3 || ' atly useri sanawa gosdy') returning id`
	sqlInsertCategorie = `insert into categories(name) values($1) returning name`
	sqlInsertMessage4  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' kategoriyasyny doretdi') returning id`
	sqlInsertWorker    = `insert into workers(fullname, degisli_dukany, wezipesi, salary, phone, home_addres, home_phone, email, file_name, note) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	sqlInsertMessage5  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly ishgari sanawa gosdy ')`
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")
	email := r.FormValue("email")

	x := function.Hash(password)

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

	adder := function.TokenData(r)

	if len(x) > 7 && function.Ascii(x) == true {
		conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		userid := 0
		err = conn.QueryRow(context.Background(), sqlInsertUser, username, email, x, role).Scan(&userid)

		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")

		ID := function.SelectUserID(adder)
		var m int
		err = conn.QueryRow(context.Background(), sqlInsertMessage3, ID, "User gosmak", username).Scan(&m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(100)
		}
	} else {
		w.WriteHeader(400)
	}

}

func CreateStore(w http.ResponseWriter, r *http.Request) {
	StoreName := r.FormValue("name")
	ParentStoreid := r.FormValue("parent_store_id")

	intParentStoreID, _ := strconv.Atoi(ParentStoreid)

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

	adder := function.TokenData(r)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1000)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	var n int
	err = conn.QueryRow(context.Background(), sqlInsertStore, StoreName, intParentStoreID).Scan(&n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(101)
	}

	var m int
	err = conn.QueryRow(context.Background(), sqlInsertMessage1, ID, "Dukan doretmek", StoreName).Scan(&m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(100)
	}

}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	CustomerName := r.FormValue("name")
	note := r.FormValue("note")

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

	adder := function.TokenData(r)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	var n int
	err = conn.QueryRow(context.Background(), sqlInsertCustomer, CustomerName, note).Scan(&n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(103)
	}

	var m int
	err = conn.QueryRow(context.Background(), sqlInsertMessage2, ID, "Musderi gosmak", CustomerName).Scan(&m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(104)
	}

}

func CreateCategorie(w http.ResponseWriter, r *http.Request) {
	CategorieName := r.FormValue("name")

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

	adder := function.TokenData(r)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	var n string
	err = conn.QueryRow(context.Background(), sqlInsertCategorie, CategorieName).Scan(&n)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(103)
	}
	var m int
	err = conn.QueryRow(context.Background(), sqlInsertMessage4, ID, "Kategoriya doretmek", CategorieName).Scan(&m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(104)
	}

}

func CreateWorker(w http.ResponseWriter, r *http.Request) {
	fullname := r.FormValue("fullname")
	degisliDukany := r.FormValue("degisli_dukany")
	wezipesi := r.FormValue("wezipesi")
	phone := r.FormValue("phone")
	salary := r.FormValue("salary")
	homeAddres := r.FormValue("home_addres")
	homePhone := r.FormValue("home_phone")
	fileName := r.FormValue("file_name")
	email := r.FormValue("email")
	note := r.FormValue("note")

	intSalary, _ := strconv.Atoi(salary)

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

	adder := function.TokenData(r)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1000)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	rows, err := conn.Exec(context.Background(), sqlInsertWorker, fullname, degisliDukany, wezipesi, intSalary, phone, homeAddres, homePhone, email, fileName, note)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(101)
	}
	if rows == nil {
		fmt.Println(rows)
	}

	rows1, err := conn.Exec(context.Background(), sqlInsertMessage5, ID, "Ishgar gosmak", fullname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(104)
	}
	if rows1 == nil {
		fmt.Println(rows1)
	}
}
