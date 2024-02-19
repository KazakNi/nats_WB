package db

import "time"

type Order struct {
	Order_uid          string    `db:"order_uid"`
	Track_number       string    `db:"track_number"`
	Entry              string    `db:"entry"`
	Locale             string    `db:"locale"`
	Internal_signature string    `db:"internal_signature"`
	Customer_id        int       `db:"customer_id"`
	Delivery_service   string    `db:"delivery_service"`
	Shardkey           int       `db:"shardkey"`
	Sm_id              int       `db:"sm_id"`
	Date_created       time.Time `db:"date_created"`
	Oof_shard          int       `db:"oof_shard"`
}

type Delivery struct {
	Order_id string `db:"order_id"` // FK
	Name     string `db:"name"`
	Phone    string `db:"phone"`
	Zip      int    `db:"zip"`
	City     string `db:"city"`
	Address  string `db:"address"`
	Region   string `db:"region"`
	Email    string `db:"email"`
}

type Payment struct {
	Order_id      string `db:"order_id"` // FK
	Transaction   string `db:"transaction"`
	Request_id    int    `db:"request_id"`
	Currency      string `db:"currency"`
	Provider      string `db:"provider"`
	Amount        int    `db:"amount"`
	Payment_dt    int    `db:"payment_dt"`
	Bank          string `db:"bank"`
	Delivery_cost int    `db:"delivery_cost"`
	Goods_total   int    `db:"goods_total"`
	Custom_fee    int    `db:"custom_fee"`
}

type Item struct {
	Order_id     string `db:"order_id"` // FK
	Chrt_id      int    `db:"chrt_id"`
	Track_number string `db:"track_number"`
	Price        int    `db:"price"`
	Rid          string `db:"rid"`
	Name         string `db:"name"`
	Sale         int    `db:"sale"`
	Size         int    `db:"size"`
	Total_price  int    `db:"total_price"`
	Nm_id        int    `db:"nm_id"`
	Brand        string `db:"brand"`
	Status       int    `db:"status"`
}
