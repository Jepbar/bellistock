package delete

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/function"
	"stock/responses"
	"strconv"

	"github.com/jackc/pgx"
)

const (
	sqlDeleteUser                     = `update users set is_it_deleted = 'True' where user_id = $1`
	sqlDeleteCustomer                 = `update customers set is_it_deleted = 'True' where customer_id = $1`
	sqlDeleteWorker                   = `update workers set is_it_deleted = 'True' where worker_id = $1`
	sqlDeleteCategorie                = `update categories set is_it_deleted = 'True' where categorie_id = $1`
	sqlInsertMessagetoDeleteUser      = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly useri sanawdan pozdy ')`
	sqlInsertMessagetoDeleteCustomer  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly musderini sanawdan pozdy ')`
	sqlInsertMessagetoDeleteWorker    = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly isgari sanawdan pozdy ')`
	sqlInsertMessagetoDeleteCategorie = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly kategoriyany pozdy ')`
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")

	IntId, _ := strconv.Atoi(Id)

	DeletedUser := function.SelectUsername(IntId)

	token := function.ExtractToken(r)
	deleter, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	DeleterId := function.SelectUserID(deleter)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), sqlDeleteUser, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), sqlInsertMessagetoDeleteUser, DeleterId, "User pozmak", DeletedUser)
	if rows1 == nil {
		fmt.Println(rows1, err1)
	}
	responses.SendResponse(w, err, nil, nil)

}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")

	IntId, _ := strconv.Atoi(Id)

	DeletedCustomer := function.SelectCustomer(IntId)

	token := function.ExtractToken(r)
	deleter, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	DeleterId := function.SelectUserID(deleter)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), sqlDeleteCustomer, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), sqlInsertMessagetoDeleteCustomer, DeleterId, "Musderi pozmak", DeletedCustomer)
	if rows1 == nil {
		fmt.Println(rows1, err1)
	}

	responses.SendResponse(w, err, nil, nil)

}

func DeleteWorker(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")

	IntId, _ := strconv.Atoi(Id)

	DeletedWorker := function.SelectWorker(IntId)

	token := function.ExtractToken(r)
	deleter, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	DeleterId := function.SelectUserID(deleter)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), sqlDeleteWorker, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), sqlInsertMessagetoDeleteWorker, DeleterId, "Isgari pozmak", DeletedWorker)
	if rows1 == nil {
		fmt.Println(rows1, err1)
	}

	responses.SendResponse(w, err, nil, nil)

}

func DeleteCategorie(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")

	IntId, _ := strconv.Atoi(Id)

	Deletedcategorie := function.SelectCategorie(IntId)
	token := function.ExtractToken(r)
	deleter, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	DeleterId := function.SelectUserID(deleter)

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), sqlDeleteCategorie, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), sqlInsertMessagetoDeleteCategorie, DeleterId, "Kategoriyany pozmak", Deletedcategorie)
	if rows1 == nil {
		fmt.Println(rows1, err1)
	}

	responses.SendResponse(w, err, nil, nil)
}
