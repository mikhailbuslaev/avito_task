package model

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
