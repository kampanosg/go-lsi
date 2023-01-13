package response

import "time"

type Auth struct {
	CustomerID            int         `json:"CustomerId"`
	FullName              string      `json:"FullName"`
	Company               string      `json:"Company"`
	ProductName           string      `json:"ProductName"`
	ExpirationDate        time.Time   `json:"ExpirationDate"`
	IsAccountHolder       bool        `json:"IsAccountHolder"`
	SessionUserID         int         `json:"SessionUserId"`
	ID                    string      `json:"Id"`
	EntityID              string      `json:"EntityId"`
	DatabaseName          string      `json:"DatabaseName"`
	DatabaseServer        interface{} `json:"DatabaseServer"`
	PrivateDatabaseServer interface{} `json:"PrivateDatabaseServer"`
	DatabaseUser          interface{} `json:"DatabaseUser"`
	DatabasePassword      interface{} `json:"DatabasePassword"`
	AppName               interface{} `json:"AppName"`
	SidRegistration       string      `json:"sid_registration"`
	UserName              string      `json:"UserName"`
	Md5Hash               string      `json:"Md5Hash"`
	Locality              string      `json:"Locality"`
	SuperAdmin            bool        `json:"SuperAdmin"`
	TTL                   int         `json:"TTL"`
	Token                 string      `json:"Token"`
	AccessToken           interface{} `json:"AccessToken"`
	GroupName             string      `json:"GroupName"`
	Device                string      `json:"Device"`
	DeviceType            string      `json:"DeviceType"`
	UserType              string      `json:"UserType"`
	Status                struct {
		State      string   `json:"State"`
		Reason     string   `json:"Reason"`
		Parameters struct{} `json:"Parameters"`
	} `json:"Status"`
	UserID     string `json:"UserId"`
	Properties struct {
	} `json:"Properties"`
	Email      string `json:"Email"`
	Server     string `json:"Server"`
	PushServer string `json:"PushServer"`
}
