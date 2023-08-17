package models

/* type OrderPrimaryKey struct {
	Id string `json:"id"`
} */

type Order struct {
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Count     int    `json:"count"`
	DateTime  string `json:"time"`
	Status    bool   `json:"status"`
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
	UserId    string `json:"user_id"`
	ProductId string `json:"product_id"`
	Count     int    `json:"count"`
	DateTime  string `json:"time"`
	Status    bool   `json:"status"`
}
