package main

import (
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unrolled/render"
	"go-debts/infrastructure"
	"go-debts/interfaces"
	"net/http"
	"os"
)

func main() {
	r := render.New(render.Options{Layout: "layout"})

	dbHandler := infrastructure.NewSqliteHandler("/var/tmp/go-debts.sqlite")

	accountController := accountController{r: r, handler: dbHandler}

	router := mux.NewRouter()
	router.HandleFunc("/accounts", accountController.index)
	router.HandleFunc("/accounts/{id:[0-9]+}", accountController.show)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Handle("/", router)
	http.ListenAndServe(":"+port, nil)
}

type AccountsListingInput interface {
	fetchAccountsForUser(userId int) accountsViewModel
}

type AccountsListingUseCase struct {
	handler        interfaces.DbHandler
	userGateway    userGateway
	accountGateway accountGateway
}

func (usecase AccountsListingUseCase) fetchAccountsForUser(userId int) accountsViewModel {
	debitor := usecase.userGateway.fetchDebitorByUserId(userId)
	return accountsViewModel{UserName: debitor.name,
		Accounts: usecase.accountGateway.fetchAccountsByDebitorId(debitor.id)}
}

type AccountDetailUseCase struct {
	accountGateway accountGateway
	paymentGateway paymentGateway
}

func (usecase AccountDetailUseCase) fetchAccountDetails(accountID int) accountViewModel {
	account := usecase.accountGateway.fetchAccountByID(accountID)
	payments := usecase.paymentGateway.fetchPaymentsByAccountID(accountID)
	return accountViewModel{Account: account, Payments: payments}
}

type accountGateway interface {
	fetchAccountsByDebitorId(debitorId int) []account
	fetchAccountByID(accountID int) account
}

type userGateway interface {
	fetchDebitorByUserId(userId int) debitor
}

type paymentGateway interface {
	fetchPaymentsByAccountID(accountID int) []payment
}
