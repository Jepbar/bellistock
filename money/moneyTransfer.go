package money

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/function"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx"
)

const (
	sqlInsert2  = `insert into money_transfers(store_id, type_of_transfer, type_of_account, currency, categorie_id, customer_id,project, type_of_income_payment, total_payment_amount,user_id) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`
	sqlUpdate   = `update stores set jemi_hasap_tmt = jemi_hasap_tmt + $1 , shahsy_hasap_tmt = shahsy_hasap_tmt + $2 where store_id = $3 returning name`
	sqlUpdate2  = `update stores set jemi_hasap_usd = jemi_hasap_usd + $1 , shahsy_hasap_usd = shahsy_hasap_usd + $2 where store_id = $3 returning name`
	sqlUpdate3  = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1 , shahsy_hasap_tmt = shahsy_hasap_tmt - $2 where shahsy_hasap_tmt > $2 and store_id = $3 returning name`
	sqlUpdate4  = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1 , shahsy_hasap_usd = shahsy_hasap_usd - $2 where shahsy_hasap_usd > $2 and store_id = $3 returning name`
	sqlUpdate5  = `update stores set jemi_hasap_tmt = jemi_hasap_tmt + $1  where store_id = $2 returning name`
	sqlUpdate6  = `update stores set jemi_hasap_usd = jemi_hasap_usd + $1  where store_id = $2 returning name`
	sqlUpdate7  = `update stores set jemi_hasap_tmt = jemi_hasap_tmt - $1  where store_id = $2 returning name`
	sqlUpdate8  = `update stores set jemi_hasap_usd = jemi_hasap_usd - $1  where store_id = $2 returning name`
	sqlparent   = `select parent_store_id from stores where store_id = $1`
	sqlInsert3  = `insert into transfers_between_stores(user_id, from_store_name, to_store_name, total_payment_amount, currency, note) values($1 ,$2 ,$3 ,$4, $5, $6) returning id`
	sqlselectid = `select store_id from stores where name = $1`
	sqlInsert4  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || $4 || $5 || $6 || $7) returning id`
	sqlInsert5  = `insert into last_modifications(user_id, action, message) values($1, $2, $3 || ' dukanyndan ' || $4 || ' dukanyna '|| $5 || $6 || ' gecirdi') returning id`
	sqlSelectID = `select user_id from users where username = $1`
	sqlUpdate9  = `update customers set girdeyjisi_tmt = girdeyjisi_tmt + $1 where customer_id = $2 returning name`
	sqlUpdate10 = `update customers set girdeyjisi_usd = girdeyjisi_usd + $1 where customer_id = $2 returning name`
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
	categorieid := r.FormValue("categorie_id")
	customerid := r.FormValue("customer_id")
	project := r.FormValue("project")
	typeOfIncomePayment := r.FormValue("type_of_income_payment")
	totalPaymentAmount := r.FormValue("total_payment_amount")

	intTotalPaymentAmount, _ := strconv.Atoi(totalPaymentAmount)
	intStoreid, _ := strconv.Atoi(storeid)
	intCustomerID, _ := strconv.Atoi(customerid)

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
		os.Exit(10)
	}
	defer conn.Close(context.Background())

	ID := function.SelectUserID(adder)

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
			err = conn.QueryRow(context.Background(), sqlUpdate9, intTotalPaymentAmount, intCustomerID).Scan(&n)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(13)
			}

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
			err = conn.QueryRow(context.Background(), sqlUpdate10, intTotalPaymentAmount, intCustomerID).Scan(&n)
			if err != nil {
				fmt.Println("error tmt")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(15)
			}
		}

		var m int
		err = conn.QueryRow(context.Background(), sqlInsert4, ID, "Pul girisi", storeid, " dukanyna ", totalPaymentAmount, currency, " girizdi").Scan(&m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(100)
		}

	}
	if typeOfTransfer == "cykdy" {
		if currency == "TMT" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate3, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
			if err != nil {
				fmt.Println("There is no enough moeney!")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(3)
			}
		}
		if currency == "USD" {
			var name string
			err = conn.QueryRow(context.Background(), sqlUpdate4, intTotalPaymentAmount, intTotalPaymentAmount, intStoreid).Scan(&name)
			if err != nil {
				fmt.Println("error usd")
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(4)
			}
		}
		var m int
		err = conn.QueryRow(context.Background(), sqlInsert4, ID, "Pul cykysy", storeid, " dukanyndan ", totalPaymentAmount, currency, " cykardy").Scan(&m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(100)
		}
	}
	b := parentStore(intStoreid)

	for j := 1; j < len(b)-1; j++ {
		if typeOfTransfer == "girdi" {
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

		if typeOfTransfer == "cykdy" {
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
	id := 0
	err = conn.QueryRow(context.Background(), sqlInsert2, storeid, typeOfTransfer, typeOfAccount, currency, categorieid,
		customerid, project, typeOfIncomePayment, totalPaymentAmount, ID).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(100)
	}

}

func BetweenStores(w http.ResponseWriter, r *http.Request) {
	//userid := r.FormValue("user_id")
	currency := r.FormValue("currency")
	fromstorename := r.FormValue("from_store_name")
	toStorename := r.FormValue("to_store_name")
	totalPaymentAmount := r.FormValue("total_payment_amount")
	note := r.FormValue("note")

	fromstoreid := IDOfStore(fromstorename)
	tostoreid := IDOfStore(toStorename)

	intTotalPaymentAmount, _ := strconv.Atoi(totalPaymentAmount)

	function.VerifyToken(r)
	function.TokenValid(r)

	tokenString := function.ExtractToken(r)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})

	for key, val := range claims {
		if key == "username" {

			conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
				os.Exit(10)
			}
			defer conn.Close(context.Background())

			k := val
			var ID int
			err = conn.QueryRow(context.Background(), sqlSelectID, k).Scan(&ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(1)
			}

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
			err = conn.QueryRow(context.Background(), sqlInsert3, ID, fromstorename, toStorename, totalPaymentAmount, currency, note).Scan(&id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				os.Exit(11)
			}
		}
		fmt.Println(token, err)
	}
}