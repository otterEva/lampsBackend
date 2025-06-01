package schemas

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
