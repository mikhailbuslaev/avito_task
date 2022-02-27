package functions

import (
	"avito_task/app/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User model.User
type Transaction model.Transaction

type Transactions struct {
	Slice []Transaction
}

type JsonData interface {
	ParseJson(w http.ResponseWriter, req *http.Request)
	WriteJson(w http.ResponseWriter, req *http.Request)
}

type Rows interface {
	Scan(rows *sql.Rows)
}

func (user *User) ParseJson(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(req.Body).Decode(user)

	if err != nil {
		log.Fatal(err)
	}
}

func (transaction *Transaction) ParseJson(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(req.Body).Decode(transaction)

	if err != nil {
		log.Fatal(err)
	}
}

func (transactions *Transactions) ParseJson(w http.ResponseWriter, req *http.Request) {

}

func Read(data JsonData, w http.ResponseWriter, req *http.Request) {
	data.ParseJson(w, req)
}

func (user *User) WriteJson(w http.ResponseWriter, req *http.Request) {

	jsonResult, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult) 
}

func (transaction *Transaction) WriteJson(w http.ResponseWriter, req *http.Request) {

	jsonResult, err := json.Marshal(transaction)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func (transactions Transactions) WriteJson(w http.ResponseWriter, req *http.Request) {

	jsonResult, err := json.Marshal(transactions.Slice)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult) 
}

func Write(data JsonData, w http.ResponseWriter, req *http.Request) {
	data.WriteJson(w, req)
}

func (user *User) Scan(rows *sql.Rows) {
	for rows.Next() {

		err := rows.Scan(&user.Id, &user.Balance)

		if err != nil {
			fmt.Println("Result Scan fail:")
			log.Fatal(err)
		} else {
			fmt.Println("Result Scan successful")
		}
	}
}

func (transaction *Transaction) Scan(rows *sql.Rows) {
	for rows.Next() {

		err := rows.Scan(&transaction.Sender, &transaction.Receiver,
			&transaction.Sum, &transaction.Status)

		if err != nil {
			fmt.Println("Result Scan fail:")
			log.Fatal(err)
		} else {
			fmt.Println("Result Scan successful")
		}
	}
}

func (transactions *Transactions) Scan(rows *sql.Rows) {
	for rows.Next() {

		transaction := Transaction{}
		err := rows.Scan(&transaction.Sender, &transaction.Receiver,
			&transaction.Sum, &transaction.Status)

		if err != nil {
			fmt.Println("Result Scan fail:")
			log.Fatal(err)
		} else {
			fmt.Println("Result Scan successful")
			transactions.Slice = append(transactions.Slice, transaction)
		}
	}
}

func ScanRows(r Rows, rows *sql.Rows) {
	r.Scan(rows)
}

