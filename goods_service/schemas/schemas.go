package schemas

type Good struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Cost        uint   `json:"cost"`
	ImageURL    string `json:"image_url"`
}

type GoodDB struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Cost        uint   `json:"cost"`
	ImageURL    string `json:"image_url"`
}
