package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"

	_ "github.com/lib/pq"
)

type Config struct {
	Host string		 	`json:"host"`
	Port int			`json:"port"`
	User string		 	`json:"user"`
	Password string		`json:"password"`
	Dbname string 		`json:"dbname"`
}

var (
	db *sql.DB
)

type Wallet struct {
	Id string
}

type Transaction struct {
	Id       string
	Sender   Wallet
	Recivier Wallet
	Sum      string
	Status   bool
}

func GetConfig() Config {
	file, err := os.Open("dbconfig.json")

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("Successfully opened dbconfig.json")
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	json.Unmarshal(byteValue, &config)

	if err != nil {
  		fmt.Println("error:", err)
	}
	return config
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
		t.Sender.Id + "') - " + t.Sum + " WHERE id = '" + t.Sender.Id + "';")

	if err != nil {
		fmt.Println("Debit query fail:")
		log.Fatal(err)
	} else {
		fmt.Println("Debit query successful")
	}

	_, err = db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.Recivier.Id + "') + " + t.Sum + " WHERE id = '" + t.Recivier.Id + "';")
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

func main() {

	var c Config
	c = GetConfig()
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Dbname)

	db = DatabaseConnect(connectionString)

	// var test_wallet Wallet
	// test_wallet.Id = "1"
	// test_wallet.GetBalance(db)

	// var test_transaction Transaction
	// test_transaction.Sum = "-20"
	// test_transaction.Sender.Id = "1"
	// test_transaction.Recivier.Id = "2"

	// test_transaction.MakeTransaction(db)
	// test_wallet.GetBalance(db)

}