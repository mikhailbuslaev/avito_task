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
	SenderId   string  `json:"SenderId"`
	RecieverId string  `json:"RecieverId"`
	Sum        float32 `json:"Sum"`
	Status     string  `json:"Status"`
	Key        string  `json:"Key"`
}

type Wallet struct {
	Id      string `json:"Id"`
	Balance float32
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

	rows, err := db.Query("SELECT balance FROM wallets where id = '" + w.Id + "';")
	if err != nil {
		fmt.Println("Select query failure:")
		log.Fatal(err)
	} else {
		fmt.Println("Select query correct:")
	}

	var result float32

	for rows.Next() {

		err := rows.Scan(&result)

		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Your balance is %g \n", result)
	}
	return result
}

func (t *TransactionTask) TransactionCheck(db *sql.DB) {

	t.Status = "uncompleted"

	wallet := &Wallet{}
	wallet.Id = t.SenderId
	wallet.Balance = wallet.GetBalance(db)

	if t.Sum > wallet.Balance {
		t.Status = "rejected: not enough money"

		fmt.Println("Transaction rejected: not enough money")

	} else if t.Sum < 0 {
		t.Status = "rejected: negative sum"

		fmt.Println("Transaction rejected: negative sum")

	} else {
		t.Status = "approved"

		fmt.Println("Transaction approved")

	}

}

func (t *TransactionTask) MakeTransaction(db *sql.DB) {

	stringSum := fmt.Sprintf("%g", t.Sum)

	if t.SenderId != "1" {
		t.RemoveSum(db, stringSum)
	}

	if t.RecieverId != "1" && t.Status != "failed: sum removing" {
		t.AddSum(db, stringSum)
	}

}

func (t *TransactionTask) RemoveSum(db *sql.DB, stringSum string) {

	_, err := db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.SenderId + "') - '" + stringSum + "' WHERE id = '" + t.SenderId + "';")

	if err != nil {
		t.Status = "failed: sum removing"
		fmt.Println("Debit query fail:")
		log.Fatal(err)
		return

	} else {
		t.Status = "completed"
		fmt.Println("Debit query successful")
	}
}

func (t *TransactionTask) AddSum(db *sql.DB, stringSum string) {

	_, err := db.Exec("UPDATE wallets set balance = (SELECT balance FROM wallets WHERE id = '" +
		t.RecieverId + "') + '" + stringSum + "' WHERE id = '" + t.RecieverId + "';")
	if err != nil {
		t.Status = "failed: sum recieving"
		fmt.Println("Recieving funds fail:")
		log.Fatal(err)
	} else {
		t.Status = "completed"
		fmt.Println("Recieving funds is successful")
	}
}
