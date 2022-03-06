package functions

import (
	"avito_task/app/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type User model.User
type Transaction model.Transaction

type Transactions struct {
	Slice []Transaction
}

type JsonData interface {
	Parse(req *http.Request) error
	Write() ([]byte, error)
}

type Rows interface {
	Scan(rows *sql.Rows) error
}

func (user *User) Parse(req *http.Request) error {
	err := req.ParseForm()

	if err != nil {
		fmt.Println(err)
	}

	err = json.NewDecoder(req.Body).Decode(user)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (transaction *Transaction) Parse(req *http.Request) error {
	err := req.ParseForm()

	if err != nil {
		fmt.Println(err)
	}

	err = json.NewDecoder(req.Body).Decode(transaction)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (transactions *Transactions) Parse(req *http.Request) error {
	err := req.ParseForm()

	if err != nil {
		fmt.Println(err)
	}

	err = json.NewDecoder(req.Body).Decode(transactions)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func ParseJson(data JsonData, req *http.Request) error {
	return data.Parse(req)
}

func (user *User) Write() ([]byte, error) {

	jsonResult, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}
	return jsonResult, err
}

func (transaction *Transaction) Write() ([]byte, error) {

	jsonResult, err := json.Marshal(transaction)
	if err != nil {
		fmt.Println(err)
	}
	return jsonResult, err
}

func (transactions Transactions) Write() ([]byte, error) {

	jsonResult, err := json.Marshal(transactions.Slice)
	if err != nil {
		fmt.Println(err)
	}
	return jsonResult, err
}

func WriteJson(data JsonData) ([]byte, error) {
	return data.Write()
}

func (user *User) Scan(rows *sql.Rows) error {
	var err error = nil

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Balance)

		if err != nil {
			fmt.Println("Result Scan fail")
		} else {
			fmt.Println("Result Scan successful")
		}
	}
	return err
}

func (transaction *Transaction) Scan(rows *sql.Rows) error {
	var err error = nil

	for rows.Next() {
		err = rows.Scan(&transaction.Sender, &transaction.Receiver,
			&transaction.Sum, &transaction.Status)

		if err != nil {
			fmt.Println("Result Scan fail")
		} else {
			fmt.Println("Result Scan successful")
		}
	}
	return err
}

func (transactions *Transactions) Scan(rows *sql.Rows) error {
	var err error = nil

	for rows.Next() {
		transaction := Transaction{}
		err = rows.Scan(&transaction.Sender, &transaction.Receiver,
			&transaction.Sum, &transaction.Status)

		if err != nil {
			fmt.Println("Result Scan fail")
		} else {
			fmt.Println("Result Scan successful")
			transactions.Slice = append(transactions.Slice, transaction)
		}
	}
	return err
}

func ScanRows(r Rows, rows *sql.Rows) error {
	return r.Scan(rows)
}
