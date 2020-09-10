package authentication

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"stock/function"
	"stock/responses"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx"
)

const (
	sqlselect = `select password from users where username = $1`
)

type UserLogin struct {
	AccessToken string `json:"access_token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	conn, err := pgx.Connect(context.Background(), os.Getenv("postgres://jepbar:bjepbar2609@localhost:5432/jepbar"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var hashedpass string
	err = conn.QueryRow(context.Background(), sqlselect, username).Scan(&hashedpass)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(2)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedpass), []byte(password))
	if err != nil {
		fmt.Println("error password")
		return
	}

	token, err := function.CreateToken(username)
	if err != nil {
		panic(err)
	}

	item := &UserLogin{
		AccessToken: token,
	}

	responses.SendResponse(w, err, item, nil)
}
