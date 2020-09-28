package delete

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
	sqlDeleteUser                = `update users set is_it_deleted = 'True' where user_id = $1`
	sqlInsertMessagetoDeleteUser = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly useri sanawdan pozdy ')`
)

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")

	IntId, _ := strconv.Atoi(Id)

	DeletedUser := function.SelectUsername(IntId)

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

	deleter := function.TokenData(r)

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

}
