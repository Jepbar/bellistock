package update

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/config"
	"stock/deletion"
	"stock/function"
	"stock/money"
	"stock/responses"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

const (
	sqlUpdateUser                           = `update users set username = $1, email = $2, role = $3 where user_id = $4`
	sqlUpdateWorker                         = `update workers set fullname = $1, degisli_dukany = $2, wezipesi = $3, salary = $4, phone = $5, home_phone = $6, home_addres = $7, email = $8, note = $9 where worker_id = $10 `
	sqlUpdatePasswordOfUser                 = `update users set password = $1 where user_id = $2`
	sqlUpdateCustomer                       = `update customers set name = $1 , note = $2 where customer_id = $3`
	sqlUpdateCategorie                      = `update categories set name =$1 where categorie_id = $2`
	sqlSelectPasswordOfUser                 = `select password from users where user_id =$1`
	sqlUpdateIncomeTransfer                 = `update money_transfers set store_id = $1, type_of_account = $2, currency = $3, categorie =$4, customer = $5 ,project = $6, type_of_income_payment =$7, total_payment_amount =$8 , date = $9, user_id = $10, keyword =$11, update_ts = $12 where id = $13`
	sqlUpdateOutcomeTransfer                = `update money_transfers set store_id = $1, type_of_account = $2, currency = $3, categorie =$4, money_gone_to =$5 , type_of_income_payment =$6, total_payment_amount =$7 , date = $8, user_id = $9, keyword =$10, update_ts = $11 where id =$12`
	sqlInsertMessageUpdatingUserData        = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli userin maglumaty tazelendi ')`
	sqlInsertMessageUpdatingUsersPassword   = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli userin passwordy tazelendi ')`
	sqlInsertMessageUpdatingWorkerData      = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli ishgarin maglumaty tazelendi ')`
	sqlInsertMessageUpdatingCustomerData    = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli musderinin maglumaty tazelendi ')`
	sqlInsertMessageUpdatingCategorieData   = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' Idli kategoryanyn ady uytgedildi ')`
	sqlInsertMessageUpdatingIncomeTransfer  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyna bolan '|| $4 || $5 || ' pul girisi tazelendi ')`
	sqlInsertMessageUpdatingOutcometransfer = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyna bolan '|| $4 || $5 || ' pul cykysy tazelendi ')`
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
	var x int
	var y string
	if function.SelectRoleOfUser(Editorid) == "Admin" {
		if len(Id) > 0 {
			x = Intid
			y = Id
		} else {
			x = Editorid
			y = strconv.Itoa(Editorid)
		}
	} else {
		x = Editorid
		y = strconv.Itoa(Editorid)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var hashedpass string
	err = conn.QueryRow(context.Background(), sqlSelectPasswordOfUser, x).Scan(&hashedpass)
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
		rows, err := conn.Exec(context.Background(), sqlUpdatePasswordOfUser, NewHashedPassword, x)
		if rows == nil {
			fmt.Println(rows, err)
		}
		rows1, err1 := conn.Exec(context.Background(), sqlInsertMessageUpdatingUsersPassword, Editorid, "Password tazelemek", y)
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

func UpdateIncomeTransferData(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	Id := r.FormValue("id")
	storeid := r.FormValue("store_id")
	typeOfAccount := r.FormValue("type_of_account")
	currency := r.FormValue("currency")
	categorie := r.FormValue("categorie")
	customer := r.FormValue("customer")
	project := r.FormValue("project")
	typeOfIncomePayment := r.FormValue("type_of_income_payment")
	totalPaymentAmount := r.FormValue("total_payment_amount")
	keyword := r.FormValue("keyword")
	Date := r.FormValue("date")

	CurrentTime := time.Now()
	IntId, _ := strconv.Atoi(Id)
	intTotalPaymentAmount, _ := strconv.Atoi(totalPaymentAmount)
	intStoreid, _ := strconv.Atoi(storeid)

	time1 := function.ChangeStringToDate(Date)

	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var date time.Time
	err = conn.QueryRow(context.Background(), deletion.SqlTakeTheTimeOfTransaction, IntId).Scan(&date)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}

	var Amount int
	var Currency string
	var Customer string
	var Storeid int
	var Userid int
	err = conn.QueryRow(context.Background(), deletion.SqlSelectTheTransaction, IntId).Scan(&Amount, &Currency, &Storeid, &Customer, &Userid)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	RoleOfEditor := function.SelectRoleOfUser(Editorid)
	if function.IsItAvaiableForDeletingTransfer(date) == true && Editorid == Userid || RoleOfEditor == "Admin" {
		Storename := function.SelectStore(Storeid)
		StringFormOfAmount := strconv.Itoa(Amount)
		b := money.ParentStore(Storeid)
		ok := false
		IsEverythingOk := true
		if Currency == "TMT" {
			rows, err := conn.Exec(context.Background(), deletion.SqlReturnTheMoneyTMT, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
				IsEverythingOk = false
			}
			rows1, err1 := conn.Exec(context.Background(), deletion.SqlGiveBackMoneyToCustomerTMT, Amount, Customer)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
				IsEverythingOk = false
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), deletion.SqlReturningMoneyFromparentsTMT, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
						IsEverythingOk = false
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), deletion.SqlUpdateTotalIncomeTMT, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
				IsEverythingOk = false
			}
		}
		if Currency == "USD" {
			rows, err := conn.Exec(context.Background(), deletion.SqlReturnTheMoneyUSD, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
				IsEverythingOk = false
			}
			rows1, err1 := conn.Exec(context.Background(), deletion.SqlGiveBackMoneyToCustomerUSD, Amount, Customer)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
				IsEverythingOk = false
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), deletion.SqlReturningMoneyFromparentsUSD, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
						IsEverythingOk = false
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), deletion.SqlUpdateTotalIncomeUSD, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
				IsEverythingOk = false
			}
		}
		IsItTransfered := false
		CanIinsertmessage := false
		if IsEverythingOk == true {
			if currency == "TMT" {
				var nm string
				err = conn.QueryRow(context.Background(), money.SqlUpdate, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&nm)
				if err != nil {
					fmt.Println("error tmt")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(12)
				}

				var n string
				err = conn.QueryRow(context.Background(), money.SqlUpdate9, intTotalPaymentAmount, customer).Scan(&n)
				if err != nil {
					fmt.Println("error tmt")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(13)
				}

				rows, err := conn.Exec(context.Background(), money.SqlUpdateTotalIncomeTmt, intTotalPaymentAmount)
				if rows == nil {
					fmt.Println(rows, err)
				}
				IsItTransfered = true

			}
			if currency == "USD" {
				var name string
				err = conn.QueryRow(context.Background(), money.SqlUpdate2, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
				if err != nil {
					fmt.Println("error usd")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(2)
				}
				var n string
				err = conn.QueryRow(context.Background(), money.SqlUpdate10, intTotalPaymentAmount, customer).Scan(&n)
				if err != nil {
					fmt.Println("error tmt")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(15)
				}
				rows, err := conn.Exec(context.Background(), money.SqlUpdateTotalIncomeUsd, intTotalPaymentAmount)
				if rows == nil {
					fmt.Println(rows, err)
				}
				IsItTransfered = true
			}
			c := money.ParentStore(intStoreid)
			if IsItTransfered == true {
				for j := 1; j < len(b)-1; j++ {
					if currency == "TMT" {
						var name string
						err = conn.QueryRow(context.Background(), money.SqlUpdate5, intTotalPaymentAmount, c[j]).Scan(&name)
						if err != nil {
							fmt.Println("error tmt")
							fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
							os.Exit(5)
						}
					}
					if currency == "USD" {
						var name string
						err = conn.QueryRow(context.Background(), money.SqlUpdate6, intTotalPaymentAmount, c[j]).Scan(&name)
						if err != nil {
							fmt.Println("error usd")
							fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
							os.Exit(8)
						}
					}
				}
				CanIinsertmessage = true
			}
			if CanIinsertmessage == true {
				k := true
				rows, err := conn.Exec(context.Background(), sqlUpdateIncomeTransfer, intStoreid, typeOfAccount, currency, categorie, customer, project, typeOfIncomePayment, intTotalPaymentAmount, time1, Editorid, keyword, CurrentTime, IntId)
				if rows == nil {
					fmt.Println(rows, err)
					k = false
				}
				if k == true {
					rows, err := conn.Exec(context.Background(), sqlInsertMessageUpdatingIncomeTransfer, Editorid, "Pul girisini tazelemek", Storename, StringFormOfAmount, Currency)
					if rows == nil {
						fmt.Println(rows, err)
						fmt.Println("hereketlere gosmadym")
					}
				}
			}
		}

		responses.SendResponse(w, err, nil, nil)
	} else {
		if function.IsItAvaiableForDeletingTransfer(date) == true {
			err = responses.ErrForbidden
			responses.SendResponse(w, err, nil, nil)
		} else {
			err = responses.ErrUnauthorized
			responses.SendResponse(w, err, nil, nil)
		}
	}
}

func UpdateOutcomeTransferData(w http.ResponseWriter, r *http.Request) {
	token := function.ExtractToken(r)
	editor, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}
	Id := r.FormValue("id")
	storeid := r.FormValue("store_id")
	typeOfAccount := r.FormValue("type_of_account")
	currency := r.FormValue("currency")
	categorie := r.FormValue("categorie")
	typeOfIncomePayment := r.FormValue("type_of_income_payment")
	totalPaymentAmount := r.FormValue("total_payment_amount")
	keyword := r.FormValue("keyword")
	MoneyGoneTo := r.FormValue("money_gone_to")
	Date := r.FormValue("date")

	CurrentTime := time.Now()
	IntId, _ := strconv.Atoi(Id)
	intTotalPaymentAmount, _ := strconv.Atoi(totalPaymentAmount)
	intStoreid, _ := strconv.Atoi(storeid)

	time1 := function.ChangeStringToDate(Date)

	Editorid := function.SelectUserID(editor)

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var Amount int
	var Currency string
	var Customer string
	var Storeid int
	var Userid int
	err = conn.QueryRow(context.Background(), deletion.SqlSelectTheTransaction, IntId).Scan(&Amount, &Currency, &Storeid, &Customer, &Userid)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	Storename := function.SelectStore(Storeid)
	StringFormOfAmount := strconv.Itoa(Amount)
	b := money.ParentStore(Storeid)

	var date time.Time
	err = conn.QueryRow(context.Background(), deletion.SqlTakeTheTimeOfTransaction, IntId).Scan(&date)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}

	RoleOfEditor := function.SelectRoleOfUser(Editorid)
	ok := false
	IsEverythingOk := true
	if function.IsItAvaiableForDeletingTransfer(date) == true && Editorid == Userid || RoleOfEditor == "Admin" {
		if Currency == "TMT" {
			rows, err := conn.Exec(context.Background(), deletion.SqlGiveBackTheMoneyTMT, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
				IsEverythingOk = false
			}
			rows1, err1 := conn.Exec(context.Background(), deletion.SqlGiveBackMoneyToUserTMT, Amount, Userid)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
				IsEverythingOk = false

			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), deletion.SqlGivingBackMoneyToparentsTMT, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
						IsEverythingOk = false

					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), deletion.SqlUpdateTotalOutcomeTMT, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
				IsEverythingOk = false

			}
		}

		if Currency == "USD" {
			rows, err := conn.Exec(context.Background(), deletion.SqlGiveBackTheMoneyUSD, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
				IsEverythingOk = false

			}
			rows1, err1 := conn.Exec(context.Background(), deletion.SqlGiveBackMoneyToUserUSD, Amount, Userid)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
				IsEverythingOk = false

			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), deletion.SqlGivingBackMoneyToparentsUSD, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
						IsEverythingOk = false

					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), deletion.SqlUpdateTotalOutcomeUSD, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
				IsEverythingOk = false

			}
		}
		var sowalgaTmt int
		var sowalgaUsd int
		err = conn.QueryRow(context.Background(), money.SqlselectSowalga, Editorid).Scan(&sowalgaTmt, &sowalgaUsd)
		if err != nil {
			fmt.Println("ERROR")
		}
		IsItTransfered := false
		CanIinsertmessage := false
		if IsEverythingOk == true {
			if currency == "TMT" {
				if sowalgaTmt >= intTotalPaymentAmount {
					var name string
					err = conn.QueryRow(context.Background(), money.SqlUpdate3, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
					if err != nil {
						fmt.Println("There is no enough moeney!")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						os.Exit(3)
					}
					rows, err := conn.Exec(context.Background(), money.SqlUpdateSowalga, intTotalPaymentAmount, Editorid)
					if rows == nil {
						fmt.Println(rows, err)
					}

					rows1, err1 := conn.Exec(context.Background(), money.SqlUpdateTotalOutcomeTmt, intTotalPaymentAmount)
					if rows1 == nil {
						fmt.Println(rows1, err1)
					}
					IsItTransfered = true

				}
			}
			if currency == "USD" {
				if sowalgaUsd >= intTotalPaymentAmount {
					var name string
					err = conn.QueryRow(context.Background(), money.SqlUpdate4, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
					if err != nil {
						fmt.Println("error usd")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						os.Exit(4)
					}
					rows, err := conn.Exec(context.Background(), money.SqlUpdateSowalga2, intTotalPaymentAmount, Editorid)
					if rows == nil {
						fmt.Println(rows, err)
					}

					rows1, err1 := conn.Exec(context.Background(), money.SqlUpdateTotalOutcomeUsd, intTotalPaymentAmount)
					if rows == nil {
						fmt.Println(rows1, err1)
					}
					IsItTransfered = true
				}
			}
			c := money.ParentStore(intStoreid)
			fmt.Println(c)
			if IsItTransfered == true {
				for j := 1; j < len(c)-1; j++ {
					if currency == "TMT" {
						var name string
						err = conn.QueryRow(context.Background(), money.SqlUpdate7, intTotalPaymentAmount, c[j]).Scan(&name)
						if err != nil {
							fmt.Println("error tmt")
							fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
							os.Exit(6)
						}
					}
					if currency == "USD" {
						var name string
						err = conn.QueryRow(context.Background(), money.SqlUpdate8, intTotalPaymentAmount, c[j]).Scan(&name)
						if err != nil {
							fmt.Println("error usd")
							fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
							os.Exit(7)
						}
					}
				}
				CanIinsertmessage = true

				if CanIinsertmessage == true {
					k := true
					rows, err := conn.Exec(context.Background(), sqlUpdateOutcomeTransfer, intStoreid, typeOfAccount, currency, categorie, MoneyGoneTo, typeOfIncomePayment, intTotalPaymentAmount, time1, Editorid, keyword, CurrentTime, IntId)
					if rows == nil {
						fmt.Println(rows, err)
						k = false
					}
					if k == true {
						rows, err := conn.Exec(context.Background(), sqlInsertMessageUpdatingOutcometransfer, Editorid, "Pul cykysyny tazelemek", Storename, StringFormOfAmount, Currency)
						if rows == nil {
							fmt.Println(rows, err)
							fmt.Println("hereketlere gosmadym")
						}
					}
				}
			}
		}
		fmt.Println(IsEverythingOk)
		responses.SendResponse(w, err, nil, nil)
	} else {
		if function.IsItAvaiableForDeletingTransfer(date) == true {
			err = responses.ErrForbidden
			responses.SendResponse(w, err, nil, nil)
		} else {
			err = responses.ErrUnauthorized
			responses.SendResponse(w, err, nil, nil)
		}
	}

}
