package update

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/config"
	"stock/function"
	"stock/responses"
	"strconv"

	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

const (
	sqlUpdateUser                         = `update users set username = $1, email = $2, role = $3 where user_id = $4`
	sqlUpdateWorker                       = `update workers set fullname = $1, degisli_dukany = $2, wezipesi = $3, salary = $4, phone = $5, home_phone = $6, home_addres = $7, email = $8, note = $9 where worker_id = $10 `
	sqlUpdatePasswordOfUser               = `update users set password = $1 where user_id = $2`
	sqlUpdateCustomer                     = `update customers set name = $1 , note = $2 where customer_id = $3`
	sqlUpdateCategorie                    = `update categories set name =$1 where categorie_id = $2`
	sqlSelectPasswordOfUser               = `select password from users where user_id =$1`
	sqlInsertMessageUpdatingUserData      = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli userin maglumaty tazelendi ')`
	sqlInsertMessageUpdatingUsersPassword = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli userin passwordy tazelendi ')`
	sqlInsertMessageUpdatingWorkerData    = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli ishgarin maglumaty tazelendi ')`
	sqlInsertMessageUpdatingCustomerData  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli musderinin maglumaty tazelendi ')`
	sqlInsertMessageUpdatingCategorieData = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli kategoryanyn ady uytgedildi ')`
)

func UpdateUserData(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	username := r.FormValue("username")
	email := r.FormValue("email")
	role := r.FormValue("role")

	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	Intid, _ := strconv.Atoi(Id)
	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), sqlUpdateUser, username, email, role, Intid)
	if rows == nil {
		fmt.Println(rows, err)
	}
	rows1, err1 := conn.Exec(context.Background(), sqlInsertMessageUpdatingUserData, Editorid, "Userin Maglumatyny tazelemek", Id)
	if rows1 == nil {
		fmt.Println(rows1, err1)
	}
	responses.SendResponse(w, err, nil, nil)

}

func UpdatePasswordOfUser(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	OldPassword := r.FormValue("old_password")
	NewPassword := r.FormValue("new_password")

	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	Intid, _ := strconv.Atoi(Id)
	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var hashedpass string
	err = conn.QueryRow(context.Background(), sqlSelectPasswordOfUser, Intid).Scan(&hashedpass)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(2)
	}
	ok := false
	err = bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(OldPassword))
	if err != nil {
		fmt.Println("KKKKKKK")
		ok = true
	}
	NewHashedPassword := function.Hash(NewPassword)
	if ok == false && len(NewPassword) > 7 && function.Ascii(NewPassword) == true {
		rows, err := conn.Exec(context.Background(), sqlUpdatePasswordOfUser, NewHashedPassword, Intid)
		if rows == nil {
			fmt.Println(rows, err)
		}
		rows1, err1 := conn.Exec(context.Background(), sqlInsertMessageUpdatingUsersPassword, Editorid, "Password tazelemek", Id)
		if rows == nil {
			fmt.Println(rows1, err1)
		}
		responses.SendResponse(w, err, nil, nil)
	} else {

		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	}
}

func UpdateWorkerData(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	Id := r.FormValue("id")
	fullname := r.FormValue("fullname")
	degisliDukany := r.FormValue("degisli_dukany")
	wezipesi := r.FormValue("wezipesi")
	phone := r.FormValue("phone")
	salary := r.FormValue("salary")
	homeAddres := r.FormValue("home_addres")
	homePhone := r.FormValue("home_phone")
	email := r.FormValue("email")
	note := r.FormValue("note")

	Intid, _ := strconv.Atoi(Id)
	IntSalary, _ := strconv.Atoi(salary)
	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())
	ok := false
	rows, err := conn.Exec(context.Background(), sqlUpdateWorker, fullname, degisliDukany, wezipesi, IntSalary, phone, homePhone, homeAddres, email, note, Intid)
	if rows == nil {
		fmt.Println(rows, err)
		ok = true
	}
	if ok == false {
		rows, err := conn.Exec(context.Background(), sqlInsertMessageUpdatingWorkerData, Editorid, "Isgarin maglumatyny tazelemk", Id)
		if rows == nil {
			fmt.Println(rows, err)
		}

		responses.SendResponse(w, err, nil, nil)
	}
	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	}
}

func UpdateCustomerData(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	Id := r.FormValue("id")
	name := r.FormValue("name")
	note := r.FormValue("note")

	Intid, _ := strconv.Atoi(Id)
	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())
	ok := false
	rows, err := conn.Exec(context.Background(), sqlUpdateCustomer, name, note, Intid)
	if rows == nil {
		fmt.Println(rows, err)
		ok = true
	}
	if ok == false {
		rows, err := conn.Exec(context.Background(), sqlInsertMessageUpdatingCustomerData, Editorid, "Musderinin maglumatyny tazelemek", Id)
		if rows == nil {
			fmt.Println(rows, err)
		}

		responses.SendResponse(w, err, nil, nil)
	}
	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	}
}

func UpdateCategorieData(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	Id := r.FormValue("id")
	name := r.FormValue("name")

	Intid, _ := strconv.Atoi(Id)
	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())
	ok := false
	rows, err := conn.Exec(context.Background(), sqlUpdateCategorie, name, Intid)
	if rows == nil {
		fmt.Println(rows, err)
		ok = true
	}
	if ok == false {
		rows, err := conn.Exec(context.Background(), sqlInsertMessageUpdatingCategorieData, Editorid, "Kategoryanyn adyny tazelemek", Id)
		if rows == nil {
			fmt.Println(rows, err)
		}

		responses.SendResponse(w, err, nil, nil)
	}
	if ok == true {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	}

}
