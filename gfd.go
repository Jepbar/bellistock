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
	r.HandleFunc("/gettransferstores", givingresponse.GetTransferBetweenStores)
	r.HandleFunc("/getlastactions", givingresponse.GetLastActions)
	r.HandleFunc("/getcustomers", givingresponse.GetCustomers)
	r.HandleFunc("/createworker", creations.CreateWorker)
	r.HandleFunc("/givemoney", money.GiveMoneyToUser)

	http.Handle("/", r)
	http.ListenAndServe("192.168.1.140:8000", nil)
}
