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
	SqlDeleteUser                           = `update users set is_it_deleted = 'True' where user_id = $1`
	SqlDeleteCustomer                       = `update customers set is_it_deleted = 'True' where customer_id = $1`
	SqlDeleteWorker                         = `update workers set is_it_deleted = 'True' where worker_id = $1`
	SqlDeleteCategorie                      = `update categories set is_it_deleted = 'True' where categorie_id = $1`
	SqlDeleteStore                          = `update stores set is_it_deleted = 'True' where store_id = $1`
	SqlUpdateParentStoresAccount            = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1, jemi_hasap_usd = jemi_hasap_usd - $2 where store_id = $3`
	SqlSelectMoneyFromDeletedStore          = `select shahsy_hasap_tmt, shahsy_hasap_usd from stores where store_id = $1`
	SqlInsertMessagetoDeleteUser            = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly useri sanawdan pozdy ')`
	SqlInsertMessagetoDeleteCustomer        = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly musderini sanawdan pozdy ')`
	SqlInsertMessagetoDeleteWorker          = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly isgari sanawdan pozdy ')`
	SqlInsertMessagetoDeleteCategorie       = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly kategoriyany pozdy ')`
	SqlInsertMessagetoDeleteStore           = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' atly dukany pozdy ')`
	SqlTakeTheTimeOfTransaction             = `select create_ts from money_transfers where id = $1`
	SqlSelectTheTransaction                 = `select total_payment_amount, currency, store_id, customer, user_id from money_transfers where id = $1`
	SqlReturnTheMoneyTMT                    = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1 , shahsy_hasap_tmt = shahsy_hasap_tmt - $2 where store_id = $3`
	SqlReturnTheMoneyUSD                    = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1 , shahsy_hasap_usd = shahsy_hasap_usd - $2 where store_id = $3`
	SqlGiveBackTheMoneyTMT                  = `update stores set jemi_hasap_tmt = jemi_hasap_tmt + $1 , shahsy_hasap_tmt = shahsy_hasap_tmt + $2 where store_id = $3`
	SqlGiveBackTheMoneyUSD                  = `update stores set jemi_hasap_usd = jemi_hasap_usd + $1 , shahsy_hasap_usd = shahsy_hasap_usd + $2 where store_id = $3`
	SqlGiveBackMoneyToCustomerTMT           = `update customers set girdeyjisi_tmt = girdeyjisi_tmt - $1 where name = $2`
	SqlGiveBackMoneyToCustomerUSD           = `update customers set girdeyjisi_usd = girdeyjisi_usd - $1 where name = $2`
	SqlGiveBackMoneyToUserTMT               = `update users set sowalga_tmt = sowalga_tmt + $1 where user_id = $2`
	SqlGiveBackMoneyToUserUSD               = `update users set sowalga_usd = sowalga_usd + $1 where user_id = $2`
	SqlReturningMoneyFromparentsTMT         = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1  where store_id = $2`
	SqlReturningMoneyFromparentsUSD         = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1  where store_id = $2`
	SqlGivingBackMoneyToparentsTMT          = `update stores set jemi_hasap_tmt = jemi_hasap_tmt + $1  where store_id = $2`
	SqlGivingBackMoneyToparentsUSD          = `update stores set jemi_hasap_usd = jemi_hasap_usd + $1  where store_id = $2`
	SqlUpdateTotalIncomeTMT                 = `update income_outcome set total_income_tmt = total_income_tmt - $1 where id = 1`
	SqlUpdateTotalIncomeUSD                 = `update income_outcome set total_income_usd = total_income_usd - $1 where id = 1`
	SqlUpdateTotalOutcomeTMT                = `update income_outcome set total_outcome_tmt = total_outcome_tmt - $1 where id = 1`
	SqlUpdateTotalOutcomeUSD                = `update income_outcome set total_outcome_usd = total_outcome_usd - $1 where id = 1`
	SqlInsertMessageToDeleteIncomeTransfer  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyna bolan ' || $4 || $5 || ' pul girisini yzyna gaytardy ')`
	SqlInsertMessageToDeleteOutcomeTransfer = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyndan cykan ' || $4 || $5 || ' pul cykysyny yzyna aldy ')`
	SqlDeleteTransfer                       = `delete from money_transfers where id = $1`
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), SqlDeleteUser, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), SqlInsertMessagetoDeleteUser, DeleterId, "User pozmak", DeletedUser)
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), SqlDeleteCustomer, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), SqlInsertMessagetoDeleteCustomer, DeleterId, "Musderi pozmak", DeletedCustomer)
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), SqlDeleteWorker, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), SqlInsertMessagetoDeleteWorker, DeleterId, "Isgari pozmak", DeletedWorker)
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Exec(context.Background(), SqlDeleteCategorie, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	rows1, err1 := conn.Exec(context.Background(), SqlInsertMessagetoDeleteCategorie, DeleterId, "Kategoriyany pozmak", Deletedcategorie)
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var MoneyInTmt int
	var MoneyInUsd int
	err = conn.QueryRow(context.Background(), SqlSelectMoneyFromDeletedStore, IntId).Scan(&MoneyInTmt, &MoneyInUsd)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}

	rows, err := conn.Exec(context.Background(), SqlDeleteStore, IntId)
	if rows == nil {
		fmt.Println(rows, err)
	}

	b := money.ParentStore(IntId)

	for j := 1; j < len(b)-1; j++ {
		rows, err := conn.Exec(context.Background(), SqlUpdateParentStoresAccount, MoneyInTmt, MoneyInUsd, b[j])
		if rows == nil {
			fmt.Println(rows, err)
		}

	}

	rows1, err1 := conn.Exec(context.Background(), SqlInsertMessagetoDeleteStore, DeleterId, " Dukany pozmak ", Storename)
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var date time.Time
	err = conn.QueryRow(context.Background(), SqlTakeTheTimeOfTransaction, IntId).Scan(&date)
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

	err = conn.QueryRow(context.Background(), SqlSelectTheTransaction, IntId).Scan(&Amount, &Currency, &Storeid, &Customer, &Userid)
	if err != nil {
		fmt.Println("error tmt")
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}

	RoleOfDeleter := function.SelectRoleOfUser(DeleterId)
	if function.IsItAvaiableForDeletingTransfer(date) == true && DeleterId == Userid || RoleOfDeleter == "Admin" {

		Storename := function.SelectStore(Storeid)
		StringFormOfAmount := strconv.Itoa(Amount)
		b := money.ParentStore(Storeid)
		ok := false

		if Currency == "TMT" {
			rows, err := conn.Exec(context.Background(), SqlReturnTheMoneyTMT, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
			}
			rows1, err1 := conn.Exec(context.Background(), SqlGiveBackMoneyToCustomerTMT, Amount, Customer)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), SqlReturningMoneyFromparentsTMT, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), SqlUpdateTotalIncomeTMT, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
			}
			rows3, err3 := conn.Exec(context.Background(), SqlDeleteTransfer, IntId)
			if rows3 == nil {

				fmt.Println("PozulmADY")
				fmt.Println(rows3, err3)
				ok = true
			}
			if ok == false {
				rows, err := conn.Exec(context.Background(), SqlInsertMessageToDeleteIncomeTransfer, DeleterId, "Pul girisini gaytarmak", Storename, StringFormOfAmount, Currency)
				if rows == nil {
					fmt.Println(rows, err)

				}
			}
		}
		if Currency == "USD" {
			rows, err := conn.Exec(context.Background(), SqlReturnTheMoneyUSD, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
			}
			rows1, err1 := conn.Exec(context.Background(), SqlGiveBackMoneyToCustomerUSD, Amount, Customer)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), SqlReturningMoneyFromparentsUSD, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), SqlUpdateTotalIncomeUSD, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
			}
			rows3, err3 := conn.Exec(context.Background(), SqlDeleteTransfer, IntId)
			if rows3 == nil {

				fmt.Println("PozulmADY")
				fmt.Println(rows3, err3)
				ok = true
			}
			if ok == false {
				rows, err := conn.Exec(context.Background(), SqlInsertMessageToDeleteIncomeTransfer, DeleterId, "Pul girirsini gaytarmak", Storename, StringFormOfAmount, Currency)
				if rows == nil {
					fmt.Println("Gosulmady")
					fmt.Println(rows, err)

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

func DeletionOfOutcomeTransfer(w http.ResponseWriter, r *http.Request) {
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

	conn, err := pgx.Connect(context.Background(), os.Getenv(config.Conf.DbConnect))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	var date time.Time
	err = conn.QueryRow(context.Background(), SqlTakeTheTimeOfTransaction, IntId).Scan(&date)
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
		var Userid int
		err = conn.QueryRow(context.Background(), SqlSelectTheTransaction, IntId).Scan(&Amount, &Currency, &Storeid, &Customer, &Userid)
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
			rows, err := conn.Exec(context.Background(), SqlGiveBackTheMoneyTMT, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
			}
			rows1, err1 := conn.Exec(context.Background(), SqlGiveBackMoneyToUserTMT, Amount, Userid)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), SqlGivingBackMoneyToparentsTMT, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), SqlUpdateTotalOutcomeTMT, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
			}
			rows3, err3 := conn.Exec(context.Background(), SqlDeleteTransfer, IntId)
			if rows3 == nil {
				fmt.Println("PozulmADY")
				fmt.Println(rows3, err3)
				ok = true
			}
			if ok == false {
				rows, err := conn.Exec(context.Background(), SqlInsertMessageToDeleteOutcomeTransfer, DeleterId, "Pul cykysyny yzyna almak", Storename, StringFormOfAmount, Currency)
				if rows == nil {
					fmt.Println(rows, err)

				}
			}
		}
		if Currency == "USD" {
			rows, err := conn.Exec(context.Background(), SqlGiveBackTheMoneyUSD, Amount, Amount, Storeid)
			if rows == nil {
				fmt.Println(rows, err)
				ok = true
			}
			rows1, err1 := conn.Exec(context.Background(), SqlGiveBackMoneyToUserUSD, Amount, Userid)
			if rows1 == nil {
				fmt.Println(rows1, err1)
				ok = true
			}

			for j := 1; j < len(b)-1; j++ {
				if ok == false {
					rows, err := conn.Exec(context.Background(), SqlGivingBackMoneyToparentsUSD, Amount, b[j])
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						fmt.Println(rows)
						ok = true
					}
				}
			}
			rows2, err2 := conn.Exec(context.Background(), SqlUpdateTotalOutcomeUSD, Amount)
			if rows2 == nil {
				fmt.Println(rows2, err2)
				ok = true
			}
			rows3, err3 := conn.Exec(context.Background(), SqlDeleteTransfer, IntId)
			if rows3 == nil {
				fmt.Println("PozulmADY")
				fmt.Println(rows3, err3)
				ok = true
			}
			if ok == false {
				rows, err := conn.Exec(context.Background(), SqlInsertMessageToDeleteOutcomeTransfer, DeleterId, "Pul cykysyny yzyna almak", Storename, StringFormOfAmount, Currency)
				if rows == nil {
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
