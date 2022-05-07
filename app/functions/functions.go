package functions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id      string  `json:"Id"`
	Balance float32 `json:"Balance"`
}

type Transaction struct {
	Sender   string  `json:"Sender"`
	Receiver string  `json:"Receiver"`
	Sum      float32 `json:"Sum"`
	Status   string  `json:"Status"`
}

type Transactions struct {
	Slice []Transaction `json:"Transactions"`
}

type JsonData interface {}

type Rows interface {
	Scan(rows *sql.Rows) error
}

func ParseReq(j JsonData, req *http.Request) error {
	err := req.ParseForm()

	if err != nil {
		fmt.Println(err)
	}

	err = json.NewDecoder(req.Body).Decode(j)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

func WriteJson(j JsonData) ([]byte, error) {

	jsonResult, err := json.Marshal(j)
	if err != nil {
		fmt.Println(err)
	}
	return jsonResult, err
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

func Validate(w http.ResponseWriter, req *http.Request) bool{
	var valid bool
	if req.Header["Key"][0] != "1111" {
		valid = false
	} else {
		valid = true
	}
	return valid
}
