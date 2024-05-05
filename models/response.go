package models

// Cette structure permet de renvoyer d'une manière correcte les reponse.
// On la place ici afin de pouvoir être utilisée partout
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
