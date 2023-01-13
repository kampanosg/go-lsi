package domain

type Auth struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
	Server string `json:"server"`
}
