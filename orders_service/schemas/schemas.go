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

type GoodDB struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Cost        uint   `json:"cost"`
	ImageURL    string `json:"image_url"`
}

type OrderItemInput struct {
	GoodID uint `json:"good_id"`
	Amount uint `json:"amount"`
}

type OrderItemOutput struct {
	OrderUUID string `json:"order_uuid"`
	GoodID    uint   `json:"good_id"`
	Amount    uint   `json:"amount"`
}

type OrderItem struct {
	OrderUUID string `json:"order_uuid"`
	UserID    uint   `json:"user_id"`
	GoodID    uint   `json:"good_id"`
	Amount    uint   `json:"amount"`
}
