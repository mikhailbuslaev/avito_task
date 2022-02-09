package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/mikhailbuslaev/avito_task/greet"
	"github.com/mikhailbuslaev/avito_task/db"
)

type TransactionTask struct {

	SenderId	string
	RecieverId 	string
	Sum			string
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func TransactionHandler(w http.ResponseWriter, req *http.Request) {

		transaction := TransactionTask{}

		err := req.ParseForm()
	
		if err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
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
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	greet.Hello()
	fmt.Println("Server run...")

	log.Fatal(server.ListenAndServe())
}