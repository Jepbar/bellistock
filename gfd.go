package main

import (
	"net/http"
	"stock/authentication"
	"stock/creations"
	"stock/delete"
	"stock/filter"
	"stock/givingresponse"
	"stock/money"

	"github.com/gorilla/mux"
)

const (
	ConnectToDatabase = "postgres://jepbar:bjepbar2609@localhost:5432/jepbar"
)

func main() {
	r := mux.NewRouter()

	/*----Login---*/

	r.HandleFunc("/api/login", authentication.Login)

	/*---Creations---*/

	r.HandleFunc("/api/useradd", creations.CreateUser)
	r.HandleFunc("/api/createstore", creations.CreateStore)
	r.HandleFunc("/api/createcustomer", creations.CreateCustomer)
	r.HandleFunc("/api/createcategorie", creations.CreateCategorie)
	r.HandleFunc("/api/createworker", creations.CreateWorker)

	/*---MoneyTransfers---*/

	r.HandleFunc("/api/transfer", money.StoreHasap)
	r.HandleFunc("/api/betweenstores", money.BetweenStores)
	r.HandleFunc("/api/givemoney", money.GiveMoneyToUser)

	/*---ResponsToFront---*/

	r.HandleFunc("/api/getusers", givingresponse.GetUsers)
	r.HandleFunc("/api/getworkers", givingresponse.GetWorkers)
	r.HandleFunc("/api/getstores", givingresponse.GetStores)
	r.HandleFunc("/api/zxcvb", givingresponse.GetAllStores)
	r.HandleFunc("/api/gettransferstores", givingresponse.GetTransferBetweenStores)
	r.HandleFunc("/api/getlastactions", givingresponse.GetLastActions)
	r.HandleFunc("/api/getcustomers", givingresponse.GetCustomers)
	r.HandleFunc("/api/getcategories", givingresponse.GetCategories)
	r.HandleFunc("/api/getmoneytransfers", givingresponse.GetMoneyTransfers)
	r.HandleFunc("/api/getincomes", givingresponse.GetIncomes)
	r.HandleFunc("/api/getoutcomes", givingresponse.GetOutcomes)

	/*---Filters---*/

	r.HandleFunc("/api/filtermoneytransfers", filter.FilterMoneyTransfers)
	r.HandleFunc("/api/filterworkers", filter.FilterWorkers)
	r.HandleFunc("/api/filterincomes", filter.FilterIncomes)
	r.HandleFunc("/api/filteroutcomes", filter.FilterOutcomes)
	r.HandleFunc("/api/filtertransferbetweenstores", filter.FilterBetweenStores)

	/*---Delete---*/

	r.HandleFunc("/api/deleteuser", delete.DeleteUser)
	r.HandleFunc("/api/deletecustomer", delete.DeleteCustomer)
	r.HandleFunc("/api/deleteworker", delete.DeleteWorker)
	r.HandleFunc("/api/deletecategorie", delete.DeleteCategorie)
	r.HandleFunc("/api/deletestore", delete.DeleteStore)

	http.Handle("/", r)
	http.ListenAndServe("localhost:8000", nil)
}
