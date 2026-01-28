package hosthatch

type ApiResult[T any] struct {
	Result string `json:"result"`
	Data   T      `json:"data"`
}

type Servers struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Id       int     `json:"id"`
	Hostname string  `json:"hostname"`
	State    string  `json:"state"`
	Created  string  `json:"created"`
	Product  Product `json:"product"`
	Network  Network `json:"network"`
	Billing  Billing `json:"billing"`
}

type Product struct {
	Id       int    `json:"_id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Image    string `json:"image"`
}

type Network struct {
	PrimaryIPAddress string `json:"primary_ipaddress"`
}

type Billing struct {
	BillingCycle  string `json:"billing_cycle"`
	NextDue       string `json:"next_due"`
	InitialCost   int    `json:"initial_cost"`
	RecurringCost int    `json:"recurring_cost"`
	Currency      string `json:"currency"`
}
