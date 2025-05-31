package schemas

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    string `json:"admin"`
}

type UserDB struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Good struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Cost        uint   `json:"cost"`
	ImageURL    string `json:"image_url"`
}
