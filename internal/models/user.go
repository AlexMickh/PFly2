package models

type User struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	ImageUrl    string   `json:"image_url"`
	Description string   `json:"description"`
	Interests   []string `json:"interests"`
}
