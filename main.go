package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 8888
	user     = "postgres"
	password = "postgres"
	dbname   = "test"
)

var (
	db *sql.DB
)

type Wallet struct {
	Id      string
}

type Transaction struct {
	Id string
	Sender Wallet
	Recipient Wallet
	Sum float32
	Status bool
}

func (w *Wallet) GetBalance(db *sql.DB) {

	output, err := db.Query("SELECT balance FROM wallets where id = '" + w.Id + "';")
	if err != nil {
		fmt.Println("Select query failure:")
		log.Fatal(err)
	} else {
		fmt.Println("Select query correct:")
	}

	for output.Next() {
		var (
			result string
		)
		if err := output.Scan(&result); err != nil {
			log.Fatal(err)
		}
		log.Printf("Your balance is %s \n", result)
	}
}

func (t *Transaction) SendMoney(*sql.DB) {

	output, err := db.Query("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '"
	+t.Sender.Id"') - "+t.Sum" WHERE id = '" + t.Sender.Id + "';")
	if err != nil {
		fmt.Println("Select query failure:")
		log.Fatal(err)
	} else {
		fmt.Println("Select query correct:")
	}
}

func DatabaseConnect(connectionString string) *sql.DB {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Connection failure:")
		log.Fatal(err)
	} else {
		fmt.Println("Connection is correct")
	}
	return db
}

func main() {

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",2
		host, port, user, password, dbname)

	db = DatabaseConnect(connectionString)

	var test_wallet Wallet
	test_wallet.Id = "1"
	test_wallet.GetBalance(db)
}
