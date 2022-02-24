package main

import (
	"avitotask/apikey"
	"avitotask/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func KeyHandler(w http.ResponseWriter, req *http.Request) {
	t := apikey.Generate()
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]time.Time)
	resp["Your key"] = t

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func GetBalanceHandler(w http.ResponseWriter, req *http.Request) {

	database := db.Connect()

	var wallet db.Wallet
	var pointw *db.Wallet = &wallet
	pointw.Read(w, req)

	wallet.Balance = wallet.GetBalance(database)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]float32)
	resp["Balance"] = wallet.Balance

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func TransactionHandler(w http.ResponseWriter, req *http.Request) {

	database := db.Connect()
	var transaction db.TransactionTask
	var pointt *db.TransactionTask = &transaction
	pointt.Read(w, req)

	transaction.TransactionCheck(database)

	if transaction.Status == "approved" {
		transaction.MakeTransaction(database)
	}
	fmt.Println(transaction.Status)
}

func GetTransactionsHandler(w http.ResponseWriter, req *http.Request) {
	database := db.Connect()
	var wallet db.Wallet
	var pointw *db.Wallet = &wallet
	pointw.Read(w, req)
	var transactions []db.TransactionTask
	if wallet.Id != "" {
		transactions = wallet.GetTransactions(database)
	} else {
		fmt.Println("get transactions fail: wallet id is empty")
	}
	fmt.Println(transactions)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/transactions", TransactionHandler).Methods("POST")
	router.HandleFunc("/getbalance", GetBalanceHandler).Methods("POST")
	router.HandleFunc("/getkey", KeyHandler).Methods("POST")
	router.HandleFunc("/gettransactions", GetTransactionsHandler).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:1488",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server run...")

	log.Fatal(server.ListenAndServe())
}
