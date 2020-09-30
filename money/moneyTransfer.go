package money

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
	sqlInsert2               = `insert into money_transfers(store_id, type_of_transfer, type_of_account, currency, categorie, customer ,project, type_of_income_payment, total_payment_amount,user_id, date) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`
	sqlUpdate                = `update stores set jemi_hasap_tmt = jemi_hasap_tmt + $1 , shahsy_hasap_tmt = shahsy_hasap_tmt + $2 where store_id = $3 returning name`
	sqlUpdate2               = `update stores set jemi_hasap_usd = jemi_hasap_usd + $1 , shahsy_hasap_usd = shahsy_hasap_usd + $2 where store_id = $3 returning name`
	sqlUpdate3               = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1 , shahsy_hasap_tmt = shahsy_hasap_tmt - $2 where shahsy_hasap_tmt >= $2 and store_id = $3 returning name`
	sqlUpdate4               = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1 , shahsy_hasap_usd = shahsy_hasap_usd - $2 where shahsy_hasap_usd >= $2 and store_id = $3 returning name`
	sqlUpdate5               = `update stores set jemi_hasap_tmt = jemi_hasap_tmt + $1  where store_id = $2 returning name`
	sqlUpdate6               = `update stores set jemi_hasap_usd = jemi_hasap_usd + $1  where store_id = $2 returning name`
	sqlUpdate7               = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1  where store_id = $2 returning name`
	sqlUpdate8               = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1  where store_id = $2 returning name`
	sqlparent                = `select parent_store_id from stores where store_id = $1`
	sqlInsert3               = `insert into transfers_between_stores(user_id, from_store_name, to_store_name, total_payment_amount, currency,type_of_account, note, date) values($1 ,$2 ,$3 ,$4, $5, $6, $7, $8) returning id`
	sqlselectid              = `select store_id from stores where name = $1`
	sqlInsert4               = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || $4 || $5 || $6 || $7) returning id`
	sqlInsert5               = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyndan ' || $4 || ' dukanyna '|| $5 || $6 || ' gecirdi ') returning id`
	sqlSelectID              = `select user_id from users where username = $1`
	sqlUpdate9               = `update customers set girdeyjisi_tmt = girdeyjisi_tmt + $1 where name = $2 returning name`
	sqlUpdate10              = `update customers set girdeyjisi_usd = girdeyjisi_usd + $1 where name = $2 returning name`
	sqlselectSowalga         = `select sowalga_tmt, sowalga_usd from users where user_id = $1`
	sqlUpdateSowalga         = `update users set sowalga_tmt = sowalga_tmt - $1 where user_id = $2`
	sqlUpdateSowalga2        = `update users set sowalga_usd = sowalga_usd - $1 where user_id = $2`
	sqlUpdateSowalga3        = `update users set sowalga_tmt = sowalga_tmt + $1 where user_id = $2`
	sqlUpdateSowalga4        = `update users set sowalga_usd = sowalga_usd + $1 where user_id = $2`
	sqlInsertMessage         = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' usere '|| $4 || $5 || ' sowalga berdi')`
	sqlUpdateTotalIncomeTmt  = `update income_outcome set total_income_tmt = total_income_tmt + $1 where id = 1`
	sqlUpdateTotalIncomeUsd  = `update income_outcome set total_income_usd = total_income_usd + $1 where id = 1`
	sqlUpdateTotalOutcomeTmt = `update income_outcome set total_outcome_tmt = total_outcome_tmt + $1 where id = 1`
	sqlUpdateTotalOutcomeUsd = `update income_outcome set total_outcome_usd = total_outcome_usd + $1 where id = 1`
)

func IDOfStore(x string) int {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	k := x
	var storeid int
	err = conn.QueryRow(context.Background(), sqlselectid, k).Scan(&storeid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return storeid
}

func parentStore(x int) []int {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var a []int
	var parentStoreID int
	a = append(a, x)
	for k := x; k > 0; {
		err := conn.QueryRow(context.Background(), sqlparent, k).Scan(&parentStoreID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}
		a = append(a, parentStoreID)
		k = parentStoreID
	}
	return a

}

func StoreHasap(w http.ResponseWriter, r *http.Request) {
	storeid := r.FormValue("store_id")
	typeOfTransfer := r.FormValue("type_of_transfer")
	typeOfAccount := r.FormValue("type_of_account")
	currency := r.FormValue("currency")
	categorie := r.FormValue("categorie")
	customer := r.FormValue("customer")
	project := r.FormValue("project")
	typeOfIncomePayment := r.FormValue("type_of_income_payment")
	totalPaymentAmount := r.FormValue("total_payment_amount")
	date := r.FormValue("date")

	intTotalPaymentAmount, _ := strconv.Atoi(totalPaymentAmount)
	intStoreid, _ := strconv.Atoi(storeid)

	time1 := function.ChangeStringToDate(date)

	token := function.ExtractToken(r)
	adder, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	var sowalgaTmt int
	var sowalgaUsd int

	err = conn.QueryRow(context.Background(), sqlselectSowalga, ID).Scan(&sowalgaTmt, &sowalgaUsd)
	if err != nil {
		fmt.Println("ERROR")
	}
	nok := false
	lok := false
	if typeOfTransfer == "girdi" {
		if currency == "TMT" {
			var nm string
			err = conn.QueryRow(context.Background(), sqlUpdate, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&nm)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(12)
			}

			var n string
			err = conn.QueryRow(context.Background(), sqlUpdate9, intTotalPaymentAmount, customer).Scan(&n)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(13)
			}

			rows, err := conn.Exec(context.Background(), sqlUpdateTotalIncomeTmt, intTotalPaymentAmount)
			if rows == nil {
				fmt.Println(rows, err)
			}
			nok = true

		}
		if currency == "USD" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate2, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
			if err != nil {
				fmt.Println("error usd")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(2)
			}
			var n string
			err = conn.QueryRow(context.Background(), sqlUpdate10, intTotalPaymentAmount, customer).Scan(&n)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(15)
			}
			rows, err := conn.Exec(context.Background(), sqlUpdateTotalIncomeUsd, intTotalPaymentAmount)
			if rows == nil {
				fmt.Println(rows, err)
			}
			nok = true

		}
		if nok == true {
			var m int
			err = conn.QueryRow(context.Background(), sqlInsert4, ID, "Pul girisi", storeid, " dukanyna ", totalPaymentAmount, currency, " girizdi").Scan(&m)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(100)
			}
		}

	}
	ok := false
	if typeOfTransfer == "cykdy" {
		if currency == "TMT" {
			if sowalgaTmt >= intTotalPaymentAmount {
				var name string
				err = conn.QueryRow(context.Background(), sqlUpdate3, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
				if err != nil {
					fmt.Println("There is no enough moeney!")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(3)
				}
				rows, err := conn.Exec(context.Background(), sqlUpdateSowalga, intTotalPaymentAmount, ID)
				if rows == nil {
					fmt.Println(rows, err)
				}

				rows1, err1 := conn.Exec(context.Background(), sqlUpdateTotalOutcomeTmt, intTotalPaymentAmount)
				if rows == nil {
					fmt.Println(rows1, err1)
				}
				ok = true

			}
		}
		if currency == "USD" {
			if sowalgaUsd >= intTotalPaymentAmount {
				var name string
				err = conn.QueryRow(context.Background(), sqlUpdate4, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
				if err != nil {
					fmt.Println("error usd")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(4)
				}
				rows, err := conn.Exec(context.Background(), sqlUpdateSowalga2, intTotalPaymentAmount, ID)
				if rows == nil {
					fmt.Println(rows, err)
				}

				rows1, err1 := conn.Exec(context.Background(), sqlUpdateTotalOutcomeUsd, intTotalPaymentAmount)
				if rows == nil {
					fmt.Println(rows1, err1)
				}
				ok = true
			}
		}
		if ok == true {
			var m int
			err = conn.QueryRow(context.Background(), sqlInsert4, ID, "Pul cykysy", storeid, " dukanyndan ", totalPaymentAmount, currency, " cykardy").Scan(&m)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(100)
			}
		}
	}

	b := parentStore(intStoreid)

	for j := 1; j < len(b)-1; j++ {
		if typeOfTransfer == "girdi" && nok == true {
			if currency == "TMT" {
				var name string
				err = conn.QueryRow(context.Background(), sqlUpdate5, intTotalPaymentAmount, b[j]).Scan(&name)
				if err != nil {
					fmt.Println("error tmt")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(5)
				}
			}
			if currency == "USD" {
				var name string
				err = conn.QueryRow(context.Background(), sqlUpdate6, intTotalPaymentAmount, b[j]).Scan(&name)
				if err != nil {
					fmt.Println("error usd")
					fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
					os.Exit(8)
				}
			}
		}

		if typeOfTransfer == "cykdy" && ok == true {
			if ok == true {
				if currency == "TMT" {
					var name string
					err = conn.QueryRow(context.Background(), sqlUpdate7, intTotalPaymentAmount, b[j]).Scan(&name)
					if err != nil {
						fmt.Println("error tmt")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						os.Exit(6)
					}
				}
				if currency == "USD" {
					var name string
					err = conn.QueryRow(context.Background(), sqlUpdate8, intTotalPaymentAmount, b[j]).Scan(&name)
					if err != nil {
						fmt.Println("error usd")
						fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
						os.Exit(7)
					}
				}
			}
		}
	}
	if ok == true || nok == true {
		id := 0
		err = conn.QueryRow(context.Background(), sqlInsert2, intStoreid, typeOfTransfer, typeOfAccount, currency, categorie,
			customer, project, typeOfIncomePayment, intTotalPaymentAmount, ID, time1).Scan(&id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(100)
		}
		lok = true
	}
	if lok == true {

		responses.SendResponse(w, err, nil, nil)
	} else {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

}

func BetweenStores(w http.ResponseWriter, r *http.Request) {
	currency := r.FormValue("currency")
	fromstorename := r.FormValue("from_store_name")
	toStorename := r.FormValue("to_store_name")
	totalPaymentAmount := r.FormValue("total_payment_amount")
	typeOfAccount := r.FormValue("type_of_account")
	date := r.FormValue("date")
	note := r.FormValue("note")

	fromstoreid := IDOfStore(fromstorename)
	tostoreid := IDOfStore(toStorename)

	intTotalPaymentAmount, _ := strconv.Atoi(totalPaymentAmount)

	time1 := function.ChangeStringToDate(date)

	token := function.ExtractToken(r)
	adder, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	if currency == "TMT" {
		var name string
		err = conn.QueryRow(context.Background(), sqlUpdate3, intTotalPaymentAmount, intTotalPaymentAmount, fromstoreid).Scan(&name)
		if err != nil {
			fmt.Println("error tmt")
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(12)
		}
		var currency string
		err = conn.QueryRow(context.Background(), sqlUpdate, intTotalPaymentAmount, intTotalPaymentAmount, tostoreid).Scan(&currency)
		if err != nil {
			fmt.Println("error tmt")
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(13)
		}

	}
	if currency == "USD" {
		var asd string
		err = conn.QueryRow(context.Background(), sqlUpdate4, intTotalPaymentAmount, intTotalPaymentAmount, fromstoreid).Scan(&asd)
		if err != nil {
			fmt.Println("error tmt")
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(15)
		}

		var as string
		err = conn.QueryRow(context.Background(), sqlUpdate2, intTotalPaymentAmount, intTotalPaymentAmount, tostoreid).Scan(&as)
		if err != nil {
			fmt.Println("error tmt")
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(15)
		}
	}

	c := parentStore(fromstoreid)
	d := parentStore(tostoreid)

	for j := 1; j < len(c)-1; j++ {
		if currency == "TMT" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate7, intTotalPaymentAmount, c[j]).Scan(&name)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(16)
			}
		}
		if currency == "USD" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate8, intTotalPaymentAmount, c[j]).Scan(&name)
			if err != nil {
				fmt.Println("error usd")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(17)
			}
		}
	}

	for j := 1; j < len(d)-1; j++ {
		if currency == "TMT" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate5, intTotalPaymentAmount, d[j]).Scan(&name)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(18)
			}
		}
		if currency == "USD" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate6, intTotalPaymentAmount, d[j]).Scan(&name)
			if err != nil {
				fmt.Println("error usd")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(19)
			}
		}
	}
	var m int
	err = conn.QueryRow(context.Background(), sqlInsert5, ID, "Dukandan dukana pul gecirim", fromstorename, toStorename, totalPaymentAmount, currency).Scan(&m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(11)
	}

	id := 0
	err = conn.QueryRow(context.Background(), sqlInsert3, ID, fromstorename, toStorename, totalPaymentAmount, currency, typeOfAccount, note, time1).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(11)
	}

	responses.SendResponse(w, err, nil, nil)
}

func GiveMoneyToUser(w http.ResponseWriter, r *http.Request) {
	amount := r.FormValue("amount")
	currency := r.FormValue("currency")
	userid := r.FormValue("user_id")

	intAmount, _ := strconv.Atoi(amount)
	intUserid, _ := strconv.Atoi(userid)

	token := function.ExtractToken(r)
	adder, err := function.VerifyAccessToken(token)
	if err != nil {
		err = responses.ErrForbidden
		responses.SendResponse(w, err, nil, nil)
		return
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

	if currency == "TMT" {

		rows, err := conn.Exec(context.Background(), sqlUpdateSowalga3, intAmount, intUserid)
		if rows == nil {
			fmt.Println(rows, err)
		}
	}

	if currency == "USD" {
		rows, err := conn.Exec(context.Background(), sqlUpdateSowalga4, intAmount, intUserid)
		if rows == nil {
			fmt.Println(rows, err)
		}
	}

	rows, err := conn.Exec(context.Background(), sqlInsertMessage, ID, "Usere sowalga bermek", userid, amount, currency)
	if rows == nil {
		fmt.Println(rows, err)
	}
	responses.SendResponse(w, err, nil, nil)
}
