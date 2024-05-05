package models

type User struct {
	Username   string `json:"username"`
	Name       string `json:"name"`
	Passphrase string `json:"passphrase"`
	Avatar     string `json:"avatar"`
}

type Users struct {
	Users []User `json:"users"`
}
