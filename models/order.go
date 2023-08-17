package models

type OrderPrimaryKey struct {
	Id string `json:"id"`
}

type Order struct {
	Id         string             `json:"id"`
	UserId     string             `json:"user_id"`
	ProductId  string             `json:"product_id"`
	Count      int                `json:"count"`
	DateTime   string             `json:"time"`
	Status     bool               `json:"status"`
	OrderItems []*CreateOrderItem `json:"order_items"`
	Sum        int                `json:"sum"`
}

type CreateOrderItem struct {
	Id         string `json:"id"`
	ProductId  string `json:"product_id"`
	OrderId    string `json:"order_id"`
	Count      int    `json:"count"`
	TotalPrice int    `json:"total_price"`
}

type OrderGetList struct {
	Count  int
	Orders []*Order
}
type CreateOrder struct {
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Count     int    `json:"count"`
	DateTime  string `json:"time"`
	Status    bool   `json:"status"`
}

type OrderGetListRequest struct {
	Offset int
	Limit  int
}

type UpdateOrder struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Count     int    `json:"count"`
	DateTime  string `json:"time"`
	Status    bool   `json:"status"`
}
