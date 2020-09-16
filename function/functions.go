package function

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx"
	"golang.org/x/crypto/bcrypt"
)

const (
	sqlSelect = `select user_id from users where username = $1`
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
