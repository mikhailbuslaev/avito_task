package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func Connect() (*sql.DB, error) {
	file, err := os.Open("dbconfig.json")

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("Successfully opened dbconfig.json")
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	var config Config
	json.Unmarshal(byteValue, &config)

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", config.Host,
		config.Port, config.User, config.Password, config.Dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Database opening fail")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection fail")
	} else {
		fmt.Println("Database connection is successful")
	}
	return db, err
}

func Change(db *sql.DB, request string) (string, error) {
	_, err := db.Exec(request)
	if err != nil {
		fmt.Println("Change data failure")
		return "failed", err
	} else {
		fmt.Println("Change data correct")
		return "successful", err
	}
}

func Select(db *sql.DB, request string) (*sql.Rows, error) {
	result, err := db.Query(request)

	if err != nil {
		fmt.Println("Select query failure")
	} else {
		fmt.Println("Select query correct")
	}
	return result, err
}

func AddSum(db *sql.DB, receiver string, sum float32) (string, error) {
	request := fmt.Sprintf("UPDATE wallets set balance = "+
		"(SELECT balance FROM wallets WHERE id = '%s') + '%g' WHERE id = '%s';",
		receiver, sum, receiver)
	status, err := Change(db, request)
	status = "/adding money " + status
	return status, err
}

func RemoveSum(db *sql.DB, sender string, sum float32) (string, error) {
	request := fmt.Sprintf("UPDATE wallets set balance = "+
		"(SELECT balance FROM wallets WHERE id = '%s') - '%g' WHERE id = '%s';",
		sender, sum, sender)
	status, err := Change(db, request)
	status = "/removing money " + status
	return status, err
}

func GetTransactions(db *sql.DB, userid string, limit int) (*sql.Rows, error) {
	return Select(db, "SELECT SenderId, RecieverId, Sum, Status FROM "+
		"transactions where Senderid='"+userid+"' OR Recieverid='"+
		userid+"' LIMIT '"+fmt.Sprintf("%d", limit)+"';")
}

func RecordTransaction(db *sql.DB, sender string, receiver string,
	sum float32, status string) error {
	request := fmt.Sprintf("INSERT INTO transactions VALUES"+
		"('%s', '%s', '%g', '%s');", sender, receiver, sum, status)
	_, err := Change(db, request)
	return err
}
