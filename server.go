package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	. "github.com/mikhailbuslaev/avito_task/db"
	"github.com/mikhailbuslaev/avito_task/greet"
)

func HomeHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func (t *TransactionTask) TransactionCheck(db *sql.DB) {

	var wallet *Wallet
	wallet.Id = t.SenderId
	balance := wallet.GetBalance(db)
	if t.Sum > balance {
		fmt.Println(balance)
	}

}

func TransactionHandler(w http.ResponseWriter, req *http.Request) {

	transaction := TransactionTask{}

	err := req.ParseForm()

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(transaction.SenderId)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/transactions", TransactionHandler).Methods("GET", "POST")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:1488",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	greet.Hello()
	connstring := GetConfig()
	fmt.Println(connstring)
	DatabaseConnect(connstring)

	fmt.Println("Server run...")

	log.Fatal(server.ListenAndServe())
}
