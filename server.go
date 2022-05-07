package main

import (
	"avito_task/app/db"
	functions "avito_task/app/functions"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getbalance", GetBalanceHandler).Methods("POST")
	r.HandleFunc("/maketransaction", MakeTransactionHandler).Methods("POST")
	r.HandleFunc("/gettransactions", GetTransactionsHandler).Methods("POST")
	r.HandleFunc("/changebalance", ChangeBalanceHandler).Methods("POST")

	s := &http.Server{
		Addr:           ":1111",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	if functions.Validate(w, r) == true {

		database, dbError := db.Connect() //connect to database, config in dbconfig.json
		if dbError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t connect to database"))
			return
		}
	
		transactions := &functions.Transactions{}
	
		user := &functions.User{}
	
		readError := functions.ParseReq(user, r) //parsing json body of request
		if readError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("400: server can`t parse your request"))
			return
		}
	
		rows, getError := db.GetTransactions(database, user.Id, 10) //getting transacctions from database
		if getError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t get transactions info from database"))
			return
		}
	
		scanError := functions.ScanRows(transactions, rows) //scan sql result to transaction slice
		if scanError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t operate results from database"))
			return
		}
	
		JsonBody, writingJsonError := functions.WriteJson(transactions) //converting transaction slice to json
		if writingJsonError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t write json response for you"))
			return
		}
	
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(JsonBody) //send response
	}
}

func MakeTransactionHandler(w http.ResponseWriter, r *http.Request) {

	if functions.Validate(w, r) == true {
		database, dbError := db.Connect() //connecting to database
		if dbError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t connect to database"))
			return
		}
	
		transaction := &functions.Transaction{}
	
		readError := functions.ParseReq(transaction, r) //parsing json body request
		if readError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("400: server can`t parse your json"))
			return
		}
	
		transaction.Status = transaction.Status + "/pending" //updating status
	
		sender := &functions.User{Id: transaction.Sender}
	
		rows, dbError := db.Select(database, "SELECT id, balance FROM "+ //getting balance from database
			"wallets where id='"+sender.Id+"';")
		if dbError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t get balance from database"))
			return
		}
	
		scanError := functions.ScanRows(sender, rows) //scan sql result to User struct
		if scanError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t scan result from database"))
			return
		}
	
		if sender.Balance < transaction.Sum { //checking that sender balance more than transaction sum
			fmt.Println("Transaction failed, sender haven`t enough amount of money")
			transaction.Status = transaction.Status + "/canceled, sender haven`t enough amount of money"
		} else {
			addSumStatus, addSumError := db.AddSum(database, transaction.Receiver, transaction.Sum) //change receiver balance in database
			transaction.Status = transaction.Status + addSumStatus                                  //updating status
			if addSumError != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: server can`t operate transaction"))
				return
			}
	
			removeSumStatus, removeSumError := db.AddSum(database, transaction.Sender, transaction.Sum) //change sender balance in database
			transaction.Status = transaction.Status + removeSumStatus                                   //updating status
			if removeSumError != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: server can`t operate transaction"))
				return
			}
		}
	
		recordToDbError := db.RecordTransaction(database, transaction.Sender, transaction.Receiver,
			transaction.Sum, transaction.Status) //saving information about transaction in database
		if recordToDbError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t record transaction info to database"))
			return
		}
	
		JsonBody, writingError := functions.WriteJson(transaction) //convert transaction to json
		if writingError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t write json response for you"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(JsonBody)// send response
	}
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	if functions.Validate(w, r) == true {

	database, dbError := db.Connect()//connect to database
	if dbError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t connect to database"))
		return
	}

	user := &functions.User{}

	readError := functions.ParseReq(user, r)//parse request to user struct
	if readError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("400: server can`t parse your request"))
		return
	}

	rows, selectError := db.Select(database, "SELECT id, balance FROM "+
		"wallets where id='"+user.Id+"';")//get balance from database
	if selectError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t get data from database"))
		return
	}

	scanError := functions.ScanRows(user, rows)//scan sql data to user struct
	if scanError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t scan data from database"))
		return
	}

	JsonBody, writingError := functions.WriteJson(user)//convert user struct to json
	if writingError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500: server can`t write json response for you"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(JsonBody)//send response
	}
}

func ChangeBalanceHandler(w http.ResponseWriter, r *http.Request) {

	if functions.Validate(w, r) == true {
		
		database, dbError := db.Connect()//connect to database
		if dbError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t connect to database"))
			return
		}
	
		transaction := &functions.Transaction{}
	
		readError := functions.ParseReq(transaction, r)//parse json to transaction struct
		if readError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("400: server can`t parse your json"))
			return
		}
	
		transaction.Status = transaction.Status + "/pending"//update status
	
		user := &functions.User{Id: transaction.Sender}
	
		rows, selectError := db.Select(database, "SELECT id, balance FROM "+
			"wallets where id='"+user.Id+"';")//get balance of user from database
		if selectError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t get data from database"))
			return
		}
	
		scanError := functions.ScanRows(user, rows)//scan sql data to user struct
		if scanError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t scan data from database"))
			return
		}
	
		if user.Balance < transaction.Sum {
			fmt.Println("Transaction failed, sender haven`t enough amount of money")
			transaction.Status = transaction.Status + "/canceled, sender haven`t enough amount of money"//update status
		} else {
			removeSumStatus, removeSumError := db.RemoveSum(database, transaction.Sender, transaction.Sum)//databbase change
			transaction.Status = transaction.Status + removeSumStatus//update status
			if removeSumError != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500: server can`t operate transaction"))
				return
			}
		}
	
		recordError := db.RecordTransaction(database, transaction.Sender, transaction.Receiver,//record transaction info to database
			transaction.Sum, transaction.Status)
		if recordError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t record transaction to database"))
			return
		}
	
		JsonBody, writeError := functions.WriteJson(transaction)//convert transaction struct to json
		if writeError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500: server can`t write json response for you"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		w.Write(JsonBody)//send response
	}
}
