package deletion

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/config"
	"stock/function"
	"stock/money"
	"stock/responses"
	"strconv"
	"time"

	"github.com/jackc/pgx"
)

const (
	sqlDeleteUser                                      = `update users set is_it_deleted = 'True' where user_id = $1`
	sqlDeleteCustomer                                  = `update customers set is_it_deleted = 'True' where customer_id = $1`
	sqlDeleteWorker                                    = `update workers set is_it_deleted = 'True' where worker_id = $1`
	sqlDeleteCategorie                                 = `update categories set is_it_deleted = 'True' where categorie_id = $1`
	sqlDeleteStore                                     = `update stores set is_it_deleted = 'True' where store_id = $1`
	sqlUpdateParentStoresAccount                       = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1, jemi_hasap_usd = jemi_hasap_usd - $2 where store_id = $3`
	sqlSelectMoneyFromDeletedStore                     = `select shahsy_hasap_tmt, shahsy_hasap_usd from stores where store_id = $1`
	sqlInsertMessagetoDeleteUser                       = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly useri sanawdan pozdy ')`
	sqlInsertMessagetoDeleteCustomer                   = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly musderini sanawdan pozdy ')`
	sqlInsertMessagetoDeleteWorker                     = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly isgari sanawdan pozdy ')`
	sqlInsertMessagetoDeleteCategorie                  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly kategoriyany pozdy ')`
	sqlInsertMessagetoDeleteStore                      = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly dukany pozdy ')`
	sqlTakeTheTimeOfTransaction                        = `select create_ts from money_transfers where id = $1`
	sqlSelectTheTransaction                            = `select total_payment_amount, currency, store_id, customer from money_transfers where id = $1`
	sqlReturnTheMoneyTMT                               = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1 , shahsy_hasap_tmt = shahsy_hasap_tmt - $2 where store_id = $3`
	sqlReturnTheMoneyUSD                               = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1 , shahsy_hasap_usd = shahsy_hasap_usd - $2 where store_id = $3`
	sqlGiveBackMoneyToCustomerTMT                      = `update customers set girdeyjisi_tmt = girdeyjisi_tmt - $1 where name = $2`
	sqlGiveBackMoneyToCustomerUSD                      = `update customers set girdeyjisi_usd = girdeyjisi_usd - $1 where name = $2`
	sqlReturningMoneyFromparentsTMT                    = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1  where store_id = $2`
	sqlReturningMoneyFromparentsUSD                    = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1  where store_id = $2`
	sqlUpdateTotalIncomeTMT                            = `update income_outcome set total_income_tmt = total_income_tmt - $1 where id = 1`
	sqlUpdateTotalIncomeUSD                            = `update income_outcome set total_income_usd = total_income_usd - $1 where id = 1`
	sqlInsertMessageToDeleteIncomeTransfer             = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyna bolan ' || $4 || $5 || ' pul girisini yzyna gaytardy ')`
	sqlInsertMessageToDeleteIncomeTransferNOTCONVERTED = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyna bolan ' || $4 || $5 || ' pul girisini sanawdan pozdy ')`
	sqlDeleteIncomeTransfer                            = `delete from money_transfers where id = $1`
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
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
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
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
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
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
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
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
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

func DeleteStore(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	deleter, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	if function.HasItGotChild(IntId) == true {
		err = responses.ErrConflict
		responses.SendResponse(w, err, nil, nil)
		return
	}

	DeleterId := function.SelectUserID(deleter)
	Storename := function.SelectStore(IntId)
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var MoneyInTmt int
	var MoneyInUsd int
	err = conn.QueryRow(context.Background(), sqlSelectMoneyFromDeletedStore, IntId).Scan(&MoneyInTmt, &MoneyInUsd)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}

	rows, err := conn.Exec(context.Background(), sqlDeleteStore, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	b := money.ParentStore(IntId)

	for j := 1; j < len(b)-1; j++ {
		rows, err := conn.Exec(context.Background(), sqlUpdateParentStoresAccount, MoneyInTmt, MoneyInUsd, b[j])
		if rows == nil {
			fmt.Println(rows, err)
		}

	}

	rows1, err1 := conn.Exec(context.Background(), sqlInsertMessagetoDeleteStore, DeleterId, " Dukany pozmak ", Storename)
	if rows == nil {
		fmt.Println(rows1, err1)
	}
	responses.SendResponse(w, err, nil, nil)
}

func DeletionOfIncomeTransfer(w http.ResponseWriter, r *http.Request) {
	Id := r.FormValue("id")
	IntId, _ := strconv.Atoi(Id)

	token := function.ExtractToken(r)
	deleter, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	DeleterId := function.SelectUserID(deleter)
	conf := config.ReadJsonFile()
	conn, err := pgx.Connect(context.Background(), os.Getenv(conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var date time.Time
	err = conn.QueryRow(context.Background(), sqlTakeTheTimeOfTransaction, IntId).Scan(&date)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	if function.IsItAvaiableForDeletingTransfer(date) == true {
		var Amount int
		var Currency string
		var Customer string
		var Storeid int
		err = conn.QueryRow(context.Background(), sqlSelectTheTransaction, IntId).Scan(&Amount, &Currency, &Storeid, &Customer)
		if err != nil {
			fmt.Println("error tmt")
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(12)
		}
		Storename := function.SelectStore(Storeid)
		StringFormOfAmount := strconv.Itoa(Amount)
		b := money.ParentStore(Storeid)
		ok := false

		if Currency == "TMT" {
			rows, err := conn.Exec(context.Background(), sqlReturnTheMoneyTMT, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
			}
			rows1, err1 := conn.Exec(context.Background(), sqlGiveBackMoneyToCustomerTMT, Amount, Customer)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), sqlReturningMoneyFromparentsTMT, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), sqlUpdateTotalIncomeTMT, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
			}
			rows3, err3 := conn.Exec(context.Background(), sqlDeleteIncomeTransfer, IntId)
			if rows3 == nil {

				fmt.Println("PozulmADY")
				fmt.Println(rows3, err3)
				ok = true
			}
			if ok == false {
				rows, err := conn.Exec(context.Background(), sqlInsertMessageToDeleteIncomeTransfer, DeleterId, "Pul girisini gaytarmak", Storename, StringFormOfAmount, Currency)
				if rows == nil {
					fmt.Println(rows, err)

				}
			}
		}
		if Currency == "USD" {
			rows, err := conn.Exec(context.Background(), sqlReturnTheMoneyUSD, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
			}
			rows1, err1 := conn.Exec(context.Background(), sqlGiveBackMoneyToCustomerUSD, Amount, Customer)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), sqlReturningMoneyFromparentsUSD, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), sqlUpdateTotalIncomeUSD, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
			}
			rows3, err3 := conn.Exec(context.Background(), sqlDeleteIncomeTransfer, IntId)
			if rows3 == nil {

				fmt.Println("PozulmADY")
				fmt.Println(rows3, err3)
				ok = true
			}
			if ok == false {
				rows, err := conn.Exec(context.Background(), sqlInsertMessageToDeleteIncomeTransfer, DeleterId, "Pul girisini gaytarmak", Storename, StringFormOfAmount, Currency)
				if rows == nil {
					fmt.Println("Gosulmady")
					fmt.Println(rows, err)

				}
			}
		}
		responses.SendResponse(w, err, nil, nil)
	}

	if function.IsItAvaiableForDeletingTransfer(date) == false {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
	}
}

func DeletionOfOutcomeTransfer(w http.ResponseWriter, r *http.Request) {

}
