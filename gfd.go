package main

import (
	"net/http"
	"stock/authentication"
	"stock/config"
	"stock/creations"
	"stock/deletion"
	"stock/filter"
	"stock/givingresponse"
	"stock/money"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	conf := config.ReadJsonFile()

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
	r.HandleFunc("/api/getallstores", givingresponse.GetAllStores)
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

	r.HandleFunc("/api/deleteuser", deletion.DeleteUser)
	r.HandleFunc("/api/deletecustomer", deletion.DeleteCustomer)
	r.HandleFunc("/api/deleteworker", deletion.DeleteWorker)
	r.HandleFunc("/api/deletecategorie", deletion.DeleteCategorie)
	r.HandleFunc("/api/deletestore", deletion.DeleteStore)
	r.HandleFunc("/api/deleteincometransfer", deletion.DeletionOfIncomeTransfer)

	http.Handle("/", r)
	http.ListenAndServe(conf.ListenAndServe, nil)
}
