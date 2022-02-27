package main

import (
	"avito_task/app/db"
	functions "avito_task/app/functions"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getbalance", GetBalanceHandler)
	r.HandleFunc("/maketransaction", MakeTransactionHandler)
	r.HandleFunc("/gettransactions", GetTransactionsHandler)

	s := &http.Server{
		Addr:           ":1111",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	database := db.Connect()

	transactions := &functions.Transactions{}

	user := &functions.User{}

	functions.Read(user, w, r)

	functions.ScanRows(transactions, db.GetTransactions(database, user.Id, 10))

	functions.Write(transactions, w, r)
}

func MakeTransactionHandler(w http.ResponseWriter, r *http.Request) {
	database := db.Connect()

	transaction := &functions.Transaction{}

	functions.Read(transaction, w, r)

	transaction.Status = transaction.Status+"/pending"

	sender := &functions.User{Id: transaction.Sender}

	functions.ScanRows(sender, db.Select(database,  "SELECT id, balance FROM "+
	"wallets where id='"+sender.Id+"';"))

	if sender.Balance < transaction.Sum {
		fmt.Println("Transaction failed, sender haven`t enough amount of money")
		transaction.Status = transaction.Status+"/canceled, sender haven`t enough amount of money"
	} else {
		transaction.Status = transaction.Status+
		db.AddSum(database, transaction.Receiver, transaction.Sum)
		
		transaction.Status = transaction.Status+
		db.RemoveSum(database, transaction.Sender, transaction.Sum)
	}

	db.RecordTransaction(database, transaction.Sender, transaction.Receiver, transaction.Sum, 
	transaction.Status)
	
	functions.Write(transaction, w, r)
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	database := db.Connect()

	user := &functions.User{}

	functions.Read(user, w, r)

	functions.ScanRows(user, db.Select(database, "SELECT id, balance FROM "+
	"wallets where id='"+user.Id+"';"))

	functions.Write(user, w, r)
}
