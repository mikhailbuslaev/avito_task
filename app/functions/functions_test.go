package functions

import (
	"testing"
	"encoding/json"
	"net/http"
	"bytes"
	"database/sql"
	"avito_task/app/db"
	"fmt"
)

func TestParseTransaction(t *testing.T) {

	var expected = [3]Transaction{{"0","0",0,"a"}, {"0","0",0,"a"}, {"0","0",0,"a"}}
	var jsonData [3][]byte
	var requests [3]*http.Request

	for i := range expected {

		var err error
		jsonData[i], err =	json.Marshal(expected[i])
		if err != nil {
			t.Errorf("Problem while marshlling json")
		} else {
			requests[i], err = http.NewRequest("POST", "/", bytes.NewReader(jsonData[i]))
		}
	}

	for i := range expected {

		have := &Transaction{}

		err := have.Parse(requests[i])

		if err != nil {
			t.Errorf("Error while scan json to struct")
		}
		if have.Sender != expected[i].Sender && 
		have.Receiver != expected[i].Receiver && 
		have.Sum != expected[i].Sum && 
		have.Status != expected[i].Status {
			t.Errorf("handler returned wrong value: got %+v want %+v", 
			have, expected[i])
		}
	}
}

func TestParseUser(t *testing.T) {

	var expected = [3]User{{"0", 1}, {"0", 2}, {"1", 4}}
	var jsonData [3][]byte
	var requests [3]*http.Request

	for i := range expected {

		var err error
		jsonData[i], err =	json.Marshal(expected[i])
		if err != nil {
			t.Errorf("Problem while marshlling json")
		} else {
			requests[i], err = http.NewRequest("POST", "/", bytes.NewReader(jsonData[i]))
		}
	}

	for i := range expected {

		have := &User{}

		err := have.Parse(requests[i])

		if err != nil {
			t.Errorf("Error while scan json to struct")
		}
		if have.Id != expected[i].Id && 
		have.Balance != expected[i].Balance {
			t.Errorf("handler returned wrong value: got %+v want %+v", 
			have, expected[i])
		}
	}
}

func TestParseTransactions(t *testing.T) {

	var expected = [3]Transactions{}
	expected[0].Slice = []Transaction{{"1","7",90,"a"}, {"2","3",10,"ak"}, {"0","0",0,"a"}}
	expected[1].Slice = []Transaction{{"2","8",10,"av"}, {"4","9",70,"aa"}, {"0","0",0,"a"}}
	expected[2].Slice = []Transaction{{"3","9",10,"aq"}, {"8","1",80,"ax"}, {"0","0",0,"a"}}
	var jsonData [3][]byte
	var requests [3]*http.Request

	for i := range expected {

		var err error
		jsonData[i], err =	json.Marshal(expected[i])
		if err != nil {
			t.Errorf("Problem while marshlling json")
		} else {
			requests[i], err = http.NewRequest("POST", "/", bytes.NewReader(jsonData[i]))
		}
	}

	for i := range expected {

		have := Transactions{}

		err := have.Parse(requests[i])

		if err != nil {
			t.Errorf("Error while scan json to struct")
		}
		for j := range expected[0].Slice {

			if have.Slice[j].Sender != expected[i].Slice[j].Sender && 
			have.Slice[j].Receiver != expected[i].Slice[j].Receiver && 
			have.Slice[j].Sum != expected[i].Slice[j].Sum && 
			have.Slice[j].Status != expected[i].Slice[j].Status {
				t.Errorf("handler returned wrong value: got %+v want %+v", 
				have.Slice[j], expected[i].Slice[j])
			}
		}
	}
}

func TestWriteUser(t *testing.T) {
	var users = [3]User{{"0", 1}, {"0", 2}, {"1", 4}}
	var expected [3][]byte

	for i := range expected {

		var err error
		expected[i], err =	json.Marshal(users[i])
		if err != nil {
			t.Errorf("Problem while marshlling json")
		}
	}

	for i := range expected {

		have, err := users[i].Write()

		if err != nil {
			t.Errorf("Error while scan json to struct")
		}
		if string(have) != string(expected[i]) {
			t.Errorf("handler returned wrong value: got %s want %s", 
			have, expected[i])
		}
	}	
}

func TestWriteTransaction(t *testing.T) {
	var transactions = [3]Transaction{{"0","0",0,"a"}, {"0","0",0,"a"}, {"0","0",0,"a"}}
	var expected [3][]byte

	for i := range expected {

		var err error
		expected[i], err =	json.Marshal(transactions[i])
		if err != nil {
			t.Errorf("Problem while marshlling json")
		}
	}

	for i := range expected {

		have, err := transactions[i].Write()

		if err != nil {
			t.Errorf("Error while scan json to struct")
		}
		if string(have) != string(expected[i]) {
			t.Errorf("handler returned wrong value: got %s want %s", 
			have, expected[i])
		}
	}	
}

func TestWriteTransactions(t *testing.T) {
	var transactions = [3]Transactions{}
	transactions[0].Slice = []Transaction{{"1","7",90,"a"}, {"2","3",10,"ak"}, {"0","0",0,"a"}}
	transactions[1].Slice = []Transaction{{"2","8",10,"av"}, {"4","9",70,"aa"}, {"0","0",0,"a"}}
	transactions[2].Slice = []Transaction{{"3","9",10,"aq"}, {"8","1",80,"ax"}, {"0","0",0,"a"}}
	var expected [3][]byte

	for i := range expected {

		var err error
		expected[i], err =	json.Marshal(transactions[i].Slice)
		if err != nil {
			t.Errorf("Problem while marshlling json")
		}
	}

	for i := range expected {

		have, err := transactions[i].Write()

		if err != nil {
			t.Errorf("Error while scan json to struct")
		}
		if string(have) != string(expected[i]) {
			t.Errorf("handler returned wrong value: got %s want %s", 
			have, expected[i])
		}
	}
}

func TestScanUser (t *testing.T) {
	database, dbError := db.Connect()
	if dbError != nil {
		t.Errorf("Error while connecting to database")
	}

	var expected = [3]User{{"1", 374.12}, {"2", 576.78}, {"3", 59}}
	var rows [3]*sql.Rows

	for i:= range rows {

		request := "SELECT id, balance FROM wallets where id='"+expected[i].Id+"';"
		var err error
		rows[i], err = database.Query(request)

		if err != nil {
			t.Errorf("Select query failure")
		}

		have := &User{}
		
		var scanError error
		scanError = have.Scan(rows[i])
		if scanError != nil {
			t.Errorf("Error while scan data to struct")
		}

		if have.Id != expected[i].Id && 
		have.Balance != expected[i].Balance {
			t.Errorf("handler returned wrong value: got %+v want %+v", 
			have, expected[i])
		}
	}
}

func TestScanTransaction (t *testing.T) {
	database, dbError := db.Connect()
	if dbError != nil {
		t.Errorf("Error while connecting to database")
	}

	var expected = [2]Transaction{{"","",0,""}, {"1","2",39.88,"completed"}}
	var rows [2]*sql.Rows
	var user User

	for i:= range rows {
		user.Id = fmt.Sprintf("%d", i)
		var err error
		rows[i], err = db.GetTransactions(database, user.Id, 1)
		if err != nil {
			t.Errorf("Select query failure")
		}
		
		var have = &Transaction{}
		var scanError error
		scanError = have.Scan(rows[i])
		if scanError != nil {
			t.Errorf("Error while scan data to struct")
		}

		if have.Sender != expected[i].Sender && 
		have.Receiver != expected[i].Receiver &&
		have.Sum != expected[i].Sum && 
		have.Status != expected[i].Status {
			t.Errorf("handler returned wrong value: got %+v want %+v", 
			*have, expected[i])
		}
	}
}