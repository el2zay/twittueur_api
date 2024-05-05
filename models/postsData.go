package models

type Post struct {
	ID         string   `json:"id"`
	Body       string   `json:"body"`
	Image      string   `json:"image"`
	Date       string   `json:"date"`
	Device     string   `json:"device"`
	Passphrase string   `json:"passphrase"`
	Likedby    []string `json:"likedby"`
	Comments   []string `json:"comments"`
	IsComment  bool     `json:"isComment"`
}

type PostRequest struct {
	Posts []Post `json:"posts"`
	LikedBy []string `json:"likedby"`
}
