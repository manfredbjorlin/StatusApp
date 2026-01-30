package exaroton

type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    T      `json:"data"`
}

type Server struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Status   int      `json:"status"`
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Players  Players  `json:"players"`
	Software Software `json:"software"`
	Shared   bool     `json:"shared"`
}

type Players struct {
	Max   int      `json:"max"`
	Count int      `json:"count"`
	List  []string `json:"list"`
}

type Software struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Account struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Verified bool    `json:"verified"`
	Credits  float64 `json:"credits"`
}

type Ram struct {
	Ram int `json:"ram"`
}
