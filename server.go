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
	r.HandleFunc("/getbalance", GetBalanceHandler).Methods("POST")
	r.HandleFunc("/maketransaction", MakeTransactionHandler).Methods("POST")
	r.HandleFunc("/gettransactions", GetTransactionsHandler).Methods("POST")
	r.HandleFunc("/changebalance", ChangeBalanceHandler).Methods("POST")

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

	database, dbError := db.Connect()
	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t connect to database"))
		return
	}

	transactions := &functions.Transactions{}

	user := &functions.User{}

	readError := functions.ParseJson(user, r)
	if readError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400: server can`t parse your request"))
		return
	}

	rows, getError := db.GetTransactions(database, user.Id, 10)
	if getError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t get transactions info from database"))
		return
	}

	scanError := functions.ScanRows(transactions, rows)
	if scanError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t operate results from database"))
		return
	}

	JsonBody, writingJsonError := functions.WriteJson(transactions)
	if writingJsonError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t write json response for you"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonBody)
}

func MakeTransactionHandler(w http.ResponseWriter, r *http.Request) {

	database, dbError := db.Connect()
	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t connect to database"))
		return
	}

	transaction := &functions.Transaction{}

	readError := functions.ParseJson(transaction, r)
	if readError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400: server can`t parse your json"))
		return
	}

	transaction.Status = transaction.Status + "/pending"

	sender := &functions.User{Id: transaction.Sender}

	rows, dbError := db.Select(database, "SELECT id, balance FROM "+
		"wallets where id='"+sender.Id+"';")
	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t get balance from database"))
		return
	}

	scanError := functions.ScanRows(sender, rows)
	if scanError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t scan result from database"))
		return
	}

	if sender.Balance < transaction.Sum {
		fmt.Println("Transaction failed, sender haven`t enough amount of money")
		transaction.Status = transaction.Status + "/canceled, sender haven`t enough amount of money"
	} else {
		addSumStatus, addSumError := db.AddSum(database, transaction.Receiver, transaction.Sum)
		transaction.Status = transaction.Status + addSumStatus
		if addSumError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t operate transaction"))
			return
		}

		removeSumStatus, removeSumError := db.AddSum(database, transaction.Sender, transaction.Sum)
		transaction.Status = transaction.Status + removeSumStatus
		if removeSumError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t operate transaction"))
			return
		}
	}

	recordToDbError := db.RecordTransaction(database, transaction.Sender, transaction.Receiver,
		transaction.Sum, transaction.Status)
	if recordToDbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t record transaction info to database"))
		return
	}

	JsonBody, writingError := functions.WriteJson(transaction)
	if writingError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t write json response for you"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonBody)
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	database, dbError := db.Connect()
	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t connect to database"))
		return
	}

	user := &functions.User{}

	readError := functions.ParseJson(user, r)
	if readError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400: server can`t parse your request"))
		return
	}

	rows, selectError := db.Select(database, "SELECT id, balance FROM "+
		"wallets where id='"+user.Id+"';")
	if selectError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t get data from database"))
		return
	}

	scanError := functions.ScanRows(user, rows)
	if scanError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t scan data from database"))
		return
	}

	JsonBody, writingError := functions.WriteJson(user)
	if writingError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t write json response for you"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonBody)
}

func ChangeBalanceHandler(w http.ResponseWriter, r *http.Request) {
	database, dbError := db.Connect()
	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t connect to database"))
		return
	}

	transaction := &functions.Transaction{}

	readError := functions.ParseJson(transaction, r)
	if readError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400: server can`t parse your json"))
		return
	}

	transaction.Status = transaction.Status + "/pending"

	user := &functions.User{Id: transaction.Sender}

	rows, selectError := db.Select(database, "SELECT id, balance FROM "+
		"wallets where id='"+user.Id+"';")
	if selectError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t get data from database"))
		return
	}

	scanError := functions.ScanRows(user, rows)
	if scanError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t scan data from database"))
		return
	}

	if user.Balance < transaction.Sum {
		fmt.Println("Transaction failed, sender haven`t enough amount of money")
		transaction.Status = transaction.Status + "/canceled, sender haven`t enough amount of money"
	} else {
		removeSumStatus, removeSumError := db.RemoveSum(database, transaction.Sender, transaction.Sum)
		transaction.Status = transaction.Status + removeSumStatus
		if removeSumError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t operate transaction"))
			return
		}
	}

	recordError := db.RecordTransaction(database, transaction.Sender, transaction.Receiver,
		transaction.Sum, transaction.Status)
	if recordError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t record transaction to database"))
		return
	}

	JsonBody, writeError := functions.WriteJson(transaction)
	if writeError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t write json response for you"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonBody)
}
