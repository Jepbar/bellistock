package main

import (
	"net/http"
	"stock/authentication"
	"stock/creations"
	"stock/givingresponse"
	"stock/money"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/useradd", creations.CreateUser)
	r.HandleFunc("/createstore", creations.CreateStore)
	r.HandleFunc("/createcustomer", creations.CreateCustomer)
	r.HandleFunc("/transfer", money.StoreHasap)
	r.HandleFunc("/login", authentication.Login)
	r.HandleFunc("/betweenstores", money.BetweenStores)
	r.HandleFunc("/createcategorie", creations.CreateCategorie)
	r.HandleFunc("/getusers", givingresponse.GetUsers)
	r.HandleFunc("/getstores", givingresponse.GetStores)

	http.Handle("/", r)
	http.ListenAndServe("192.168.1.48:8080", nil)
}
