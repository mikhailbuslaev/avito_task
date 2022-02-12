package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type TransactionTask struct {
	SenderId   string 	`json:"SenderId"`
	RecieverId string 	`json:"RecieverId"`
	Sum        float32 	`json:"Sum"`
	Status     string 	`json:"Status"`
}

var (
	db *sql.DB
)

type Wallet struct {
	Id string
}

func Connect() *sql.DB {
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

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password,
		config.Dbname)

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

func (w *Wallet) GetBalance(db *sql.DB) float32 {

	output, err := db.Query("SELECT balance FROM wallets where id = '" + w.Id + "';")
	if err != nil {
		fmt.Println("Select query failure:")
		log.Fatal(err)
	} else {
		fmt.Println("Select query correct:")
	}

	var result float32

	if err := output.Scan(&result); err != nil {
		log.Fatal(err)
	}
	log.Printf("Your balance is %g \n", result)
	return result
}

func (t *TransactionTask) TransactionCheck(db *sql.DB) {

	var wallet *Wallet
	wallet.Id = t.SenderId
	balance := wallet.GetBalance(db)
	if t.Sum > balance {
		fmt.Println(balance)
		t.Status = "rejected"
	} else {
		t.Status = "approved"
	}

}

func (t *TransactionTask) MakeTransaction(*sql.DB) {

	stringSum := fmt.Sprintf("%f", t.Sum)
	_, err := db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.SenderId + "') - " + stringSum + " WHERE id = '" + t.SenderId + "';")

	if err != nil {
		fmt.Println("Debit query fail:")
		log.Fatal(err)
	} else {
		fmt.Println("Debit query successful")
	}

	_, err = db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.RecieverId + "') + " + stringSum + " WHERE id = '" + t.RecieverId + "';")
	if err != nil {
		fmt.Println("Recieving funds fail:")
		log.Fatal(err)
	} else {
		fmt.Println("Recieving funds is successful")
	}

}
