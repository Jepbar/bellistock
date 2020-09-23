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

	r.HandleFunc("/api/login", authentication.Login)

	r.HandleFunc("/api/useradd", creations.CreateUser)
	r.HandleFunc("/api/createstore", creations.CreateStore)
	r.HandleFunc("/api/createcustomer", creations.CreateCustomer)
	r.HandleFunc("/api/createcategorie", creations.CreateCategorie)
	r.HandleFunc("/api/createworker", creations.CreateWorker)

	r.HandleFunc("/api/transfer", money.StoreHasap)
	r.HandleFunc("/api/betweenstores", money.BetweenStores)
	r.HandleFunc("/api/givemoney", money.GiveMoneyToUser)

	r.HandleFunc("/api/getusers", givingresponse.GetUsers)
	r.HandleFunc("/api/getworkers", givingresponse.GetWorkers)
	r.HandleFunc("/api/getstores", givingresponse.GetStores)
	r.HandleFunc("/api/gettransferstores", givingresponse.GetTransferBetweenStores)
	r.HandleFunc("/api/getlastactions", givingresponse.GetLastActions)
	r.HandleFunc("/api/getcustomers", givingresponse.GetCustomers)
	r.HandleFunc("/api/getmoneytransfers", givingresponse.GetMoneyTransfers)
	r.HandleFunc("/api/getincomes", givingresponse.GetIncomes)
	r.HandleFunc("/api/getoutcomes", givingresponse.GetOutcomes)

	http.Handle("/", r)
	http.ListenAndServe("localhost:8000", nil)
}
