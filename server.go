package main

import (
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

func TransactionHandler(w http.ResponseWriter, req *http.Request) {

	database := db.Connect()
	t := db.TransactionTask{}
	t = ReadHttpRequest(w, req)

	t.TransactionCheck(database)

	if t.Status == "approved" {
		t.MakeTransaction(database)
	}
}

func ReadHttpRequest(w http.ResponseWriter, req *http.Request) db.TransactionTask {

	t := db.TransactionTask{}
	err := req.ParseForm()

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		log.Fatal(err)
	}
	return t
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

	fmt.Println("Server run...")

	log.Fatal(server.ListenAndServe())
}
