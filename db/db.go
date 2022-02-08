package db

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
	Id string
}

type Transaction struct {

	SenderId	string
	RecieverId 	string
	Sum			string
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

func (t *Transaction) MakeTransaction(*sql.DB) {

	_, err := db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.SenderId + "') - " + t.Sum + " WHERE id = '" + t.SenderId + "';")

	if err != nil {
		fmt.Println("Debit query fail:")
		log.Fatal(err)
	} else {
		fmt.Println("Debit query successful")
	}

	_, err = db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.RecieverId + "') + " + t.Sum + " WHERE id = '" + t.RecieverId + "';")
	if err != nil {
		fmt.Println("Recieving funds fail:")
		log.Fatal(err)
	} else {
		fmt.Println("Recieving funds is successful")
	}

}

func DatabaseConnect(connectionString string) *sql.DB {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection fail:")
		log.Fatal(err)
	} else {
		fmt.Println("Database connection is successful")
	}
	return db
}
