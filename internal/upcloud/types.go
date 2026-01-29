package upcloud

import (
	"github.com/google/uuid"
)

type Servers struct {
	Data struct {
		Server []Server `json:"server"`
	} `json:"servers"`
}

type Server struct {
	CoreNumber   string    `json:"core_number"`
	Created      int       `json:"created"`
	Hostname     string    `json:"hostname"`
	MemoryAmount string    `json:"memory_amount"`
	Plan         string    `json:"plan"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Uuid         uuid.UUID `json:"uuid"`
	Zone         string    `json:"zone"`
}

type ServerDetails struct {
	Server Server `json:"server"`
}

type Account struct {
	Data struct {
		Credits  float64 `json:"credits"`
		Username string  `json:"username"`
	} `json:"account"`
}

type Billing struct {
	Currency    string  `json:"currency"`
	TotalAmount float64 `json:"total_amount"`
}
