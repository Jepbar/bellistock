package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/responses"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

const (
	layoutISO          = "2006-01-02"
	sqlSelect          = `select user_id from users where username = $1`
	sqlSelectUsername  = `select username from users where user_id = $1`
	sqlSelectcategorie = `select name from categories where categorie_id = $1`
	sqlSelectCustomer  = `select name from customers where customer_id = $1`
)

func Hash(x string) string {
	password := []byte(x)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)

}

func Ascii(x string) bool {
	str := x
	runes := []rune(str)
	k := 0
	var result []int
	for i := 0; i < len(runes); i++ {
		result = append(result, int(runes[i]))
	}
	for j := 0; j < len(result); j++ {
		if int(123) > int(result[j]) && int(result[j]) > int(96) {
			k = k + 1
			break
		}
	}
	for j := 0; j < len(result); j++ {
		if int(91) > int(result[j]) && int(result[j]) > int(64) {
			k = k + 1
			break
		}
	}
	for j := 0; j < len(result); j++ {
		if int(58) > int(result[j]) && int(result[j]) > int(47) {
			k = k + 1
			break
		}
	}
	if k == 3 {
		return true
	}
	return false

}

func CreateToken(username string) (string, string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["username"] = username
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return "", "", err
	}
	return token, rt, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err

	}
	return nil
}

func TokenData(r *http.Request) string {

	tokenString := ExtractToken(r)
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	if token == nil {
		fmt.Println("Hello")
	}

	username := fmt.Sprintf("%v", claims["username"])

	return username
}

func SelectUserID(x string) int {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1000)
	}
	defer conn.Close(context.Background())

	var ID int
	err = conn.QueryRow(context.Background(), sqlSelect, x).Scan(&ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	return ID
}

func SelectUsername(x int) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1000)
	}
	defer conn.Close(context.Background())

	var Username string
	err = conn.QueryRow(context.Background(), sqlSelectUsername, x).Scan(&Username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	return Username
}

func SelectCategorie(x int) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1000)
	}
	defer conn.Close(context.Background())

	var Categorie string
	err = conn.QueryRow(context.Background(), sqlSelectcategorie, x).Scan(&Categorie)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	return Categorie
}

func SelectCustomer(x int) string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1000)
	}
	defer conn.Close(context.Background())

	var Customer string
	err = conn.QueryRow(context.Background(), sqlSelectCustomer, x).Scan(&Customer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(12)
	}
	return Customer
}

func ChangeStringToDate(x string) time.Time {

	date := x
	t, _ := time.Parse(layoutISO, date)
	return t
}

func GenerateSqlFilterWorkers(filter responses.Filterworkers) (sql string, err error) {
	sql = `select worker_id , fullname, wezipesi, salary, degisli_dukany from workers`
	ok := false
	if len(filter.Name) > 0 {
		sql += ` where fullname ` + `ILIKE` + `'%` + filter.Name + `%'`
		ok = true
	}
	if len(filter.Wezipesi) > 0 {
		if ok == false {
			sql += ` where wezipesi ` + `ILIKE` + `'%` + filter.Wezipesi + `%'`
			ok = true
		} else {
			sql += ` and wezipesi ` + `ILIKE` + `'%` + filter.Wezipesi + `%'`
		}
	}
	if len(filter.DependingStore) > 0 {
		if ok == false {
			sql += ` where degisli_dukany ` + `ILIKE` + `'%` + filter.DependingStore + `%'`
			ok = true
		} else {
			sql += ` and degisli_dukany ` + `ILIKE` + `'%` + filter.DependingStore + `%'`
		}
	}

	return

}

func GenerateSqlFilterMoneyTransfers(filter responses.FilterMoneyTransfers) (sql string, err error) {
	sql = `select m.id, s.name, m.type_of_transfer, m.user_id, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id`
	ok := false

	if len(filter.Store) > 0 {
		sql += ` where s.name ` + `ILIKE` + `'%` + filter.Store + `%'`
		ok = true
	}
	if len(filter.TypeOfaccount) > 0 {
		if ok == false {
			sql += ` where m.type_of_account ` + `ILIKE` + `'%` + filter.TypeOfaccount + `%'`
			ok = true
		} else {
			sql += ` and type_of_account ` + `ILIKE` + `'%` + filter.TypeOfaccount + `%'`
		}
	}
	if len(filter.Categorie) > 0 {
		if ok == false {
			sql += ` where m.categorie ` + `ILIKE` + `'%` + filter.Categorie + `%'`
			ok = true
		} else {
			sql += ` and m.categorie ` + `ILIKE` + `'%` + filter.Categorie + `%'`
		}
	}
	if len(filter.Keyword) > 0 {
		if ok == false {
			sql += ` where m.keyword ` + `ILIKE` + `'%` + filter.Keyword + `%'`
			ok = true
		} else {
			sql += ` and m.keyword ` + `ILIKE` + `'%` + filter.Keyword + `%'`
		}
	}

	if len(filter.Begin) > 0 {
		if ok == false {
			sql += ` where m.date ` + ` >=` + `'` + filter.Begin + `'`
			ok = true
		} else {
			sql += ` and m.date ` + ` >=` + `'` + filter.Begin + `'`
		}
	}
	if len(filter.End) > 0 {
		if ok == false {
			sql += ` where m.date ` + ` <=` + `'` + filter.End + `'`
			ok = true
		} else {
			sql += ` and m.date ` + ` <=` + `'` + filter.End + `'`
		}
	}
	return

}

func GenerateSqlFilterIncomes(filter responses.FilterIncomes) (sql string, err error) {
	sql = `select m.id, s.name, m.customer, m.project, m.type_of_account, m.total_payment_amount, m.currency, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id where m.type_of_transfer = 'girdi'`

	if len(filter.Store) > 0 {
		sql += ` and s.name ` + `ILIKE` + `'%` + filter.Store + `%'`
	}

	if len(filter.TypeOfaccount) > 0 {
		sql += ` and type_of_account ` + `ILIKE` + `'%` + filter.TypeOfaccount + `%'`
	}

	if len(filter.Categorie) > 0 {
		sql += ` and m.categorie ` + `ILIKE` + `'%` + filter.Categorie + `%'`
	}

	if len(filter.Keyword) > 0 {
		sql += ` and m.keyword ` + `ILIKE` + `'%` + filter.Keyword + `%'`
	}

	if len(filter.Customer) > 0 {
		sql += ` and m.customer ` + `ILIKE` + `'%` + filter.Customer + `%'`
	}

	if len(filter.TypeOfIncomePayment) > 0 {
		sql += ` and m.type_of_income_payment ` + `=` + `'` + filter.TypeOfIncomePayment + `'`
	}

	if len(filter.Begin) > 0 {
		sql += ` and m.date ` + ` >=` + `'` + filter.Begin + `'`
	}
	if len(filter.End) > 0 {
		sql += ` and m.date ` + ` <=` + `'` + filter.End + `'`
	}
	return

}

func GenerateSqlFilterOutcomes(filter responses.FilterOutcomes) (sql string, err error) {
	sql = `select m.id, s.name, m.money_gone_to, m.total_payment_amount, m.currency, m.type_of_account, m.date, m.categorie from money_transfers m inner join stores s on s.store_id = m.store_id where m.type_of_transfer = 'cykdy'`

	if len(filter.Store) > 0 {
		sql += ` and s.name ` + `ILIKE` + `'%` + filter.Store + `%'`
	}

	if len(filter.TypeOfaccount) > 0 {
		sql += ` and type_of_account ` + `ILIKE` + `'%` + filter.TypeOfaccount + `%'`
	}

	if len(filter.Categorie) > 0 {
		sql += ` and m.categorie ` + `ILIKE` + `'%` + filter.Categorie + `%'`
	}

	if len(filter.Keyword) > 0 {
		sql += ` and m.keyword ` + `ILIKE` + `'%` + filter.Keyword + `%'`
	}

	if len(filter.MoneyGoneTo) > 0 {
		sql += ` and m.money_gone_to ` + `ILIKE` + `'%` + filter.MoneyGoneTo + `%'`
	}

	if len(filter.Begin) > 0 {
		sql += ` and m.date ` + ` >=` + `'` + filter.Begin + `'`
	}

	if len(filter.End) > 0 {
		sql += ` and m.date ` + ` <=` + `'` + filter.End + `'`
	}
	return

}

func GenerateSqlFilterTransferBetweenStores(filter responses.FilterBetweenStores) (sql string, err error) {
	sql = `select id, user_id, from_store_name, to_store_name, total_payment_amount, currency, type_of_account, date from transfers_between_stores`
	ok := false

	if len(filter.FromStoreName) > 0 {
		sql += ` where from_store_name ` + `ILIKE` + `'%` + filter.FromStoreName + `%'`
		ok = true
	}
	if len(filter.ToStoreName) > 0 {
		if ok == false {
			sql += ` where to_store_name ` + `ILIKE` + `'%` + filter.ToStoreName + `%'`
			ok = true
		} else {
			sql += ` and to_store_name ` + `ILIKE` + `'%` + filter.ToStoreName + `%'`
		}
	}
	if len(filter.TypeOfAccount) > 0 {
		if ok == false {
			sql += ` where type_of_account ` + `ILIKE` + `'%` + filter.TypeOfAccount + `%'`
			ok = true
		} else {
			sql += ` and type_of_account ` + `ILIKE` + `'%` + filter.TypeOfAccount + `%'`
		}
	}
	if len(filter.Begin) > 0 {
		if ok == false {
			sql += ` where date ` + ` >=` + `'` + filter.Begin + `'`
			ok = true
		} else {
			sql += ` and date ` + ` >=` + `'` + filter.Begin + `'`
		}
	}
	if len(filter.End) > 0 {
		if ok == false {
			sql += ` where date ` + ` <=` + `'` + filter.End + `'`
			ok = true
		} else {
			sql += ` and date ` + ` <=` + `'` + filter.End + `'`
		}
	}
	return

}
