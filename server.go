package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"github.com/mikhailbuslaev/avito_task/db"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 8888
	user     = "postgres"
	password = "postgres"
	dbname   = "test"
)

type Transaction struct {

	SenderId	string
	RecieverId 	string
	Sum			string
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func TransactionHandler(w http.ResponseWriter, req *http.Request) {
	
	db = db.DatabaseConnect(connectionString)

	t := Transaction{}

	err := req.ParseForm()
	
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
		}
	
	err = json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
		}

		t.db.MakeTransaction(db)
}

func main() {

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/transactions", TransactionHandler).Methods("GET", "POST")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server run...")

	log.Fatal(server.ListenAndServe())
}
