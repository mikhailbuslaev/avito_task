package main

import (
	"avitotask/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"avitotask/apikey"

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
	resp["Your key:"] = t

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func GetBalanceHandler(w http.ResponseWriter, req *http.Request) {

	database := db.Connect()

	wallet := ReadHttpGetBalanceRequest(w, req)

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

func ReadHttpGetBalanceRequest(w http.ResponseWriter, req *http.Request) db.Wallet {

	var wallet db.Wallet
	err := req.ParseForm()

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = json.NewDecoder(req.Body).Decode(&wallet)
	if err != nil {
		log.Fatal(err)
	}
	return wallet
}

func TransactionHandler(w http.ResponseWriter, req *http.Request) {

	database := db.Connect()
	t := db.TransactionTask{}
	t = ReadHttpTransactionRequest(w, req)

	t.TransactionCheck(database)

	if t.Status == "approved" {
		t.MakeTransaction(database)
	}
	fmt.Println(t.Status)
}

func ReadHttpTransactionRequest(w http.ResponseWriter, req *http.Request) db.TransactionTask {

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
	router.HandleFunc("/transactions", TransactionHandler).Methods("POST")
	router.HandleFunc("/getbalance", GetBalanceHandler).Methods("POST")
	router.HandleFunc("/getkey", KeyHandler).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:1488",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server run...")

	log.Fatal(server.ListenAndServe())
}
