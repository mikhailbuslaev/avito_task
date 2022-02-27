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

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", config.Host,
		config.Port, config.User, config.Password, config.Dbname)

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

func Change(db *sql.DB, request string) string {
	_, err := db.Exec(request)

	if err != nil {
		fmt.Println("Change data failure:")
		log.Fatal(err)
		return "failed"
	} else {
		fmt.Println("Change data correct")
		return "successful"
	}
}

func Select(db *sql.DB, request string) *sql.Rows {
	result, err := db.Query(request)

	if err != nil {
		fmt.Println("Select query failure:")
		log.Fatal(err)
	} else {
		fmt.Println("Select query correct")
	}
	return result
}

func AddSum(db *sql.DB, receiver string, sum float32) string {
	request := fmt.Sprintf("UPDATE wallets set balance = "+
		"(SELECT balance FROM wallets WHERE id = '%s') + '%g' WHERE id = '%s';",
		receiver, sum, receiver)
	return "/adding money " + Change(db, request)
}

func RemoveSum(db *sql.DB, sender string, sum float32) string {
	request := fmt.Sprintf("UPDATE wallets set balance = "+
		"(SELECT balance FROM wallets WHERE id = '%s') - '%g' WHERE id = '%s';",
		sender, sum, sender)
	return "/removing money " + Change(db, request)
}

func GetTransactions(db *sql.DB, userid string, limit int) *sql.Rows {
	return Select(db, "SELECT SenderId, RecieverId, Sum, Status FROM "+
		"transactions where Senderid='"+userid+"' OR Recieverid='"+
		userid+"' LIMIT '"+fmt.Sprintf("%d", limit)+"';")
}

func RecordTransaction(db *sql.DB, sender string, receiver string,
	sum float32, status string) {
	request := fmt.Sprintf("INSERT INTO transactions VALUES"+
		"('%s', '%s', '%g', '%s');", sender, receiver, sum, status)
	_ = Change(db, request)
}
