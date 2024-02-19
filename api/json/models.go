package json

import "time"

type Order struct {
	Order_uid          string `json:"order_uid" validate:"required"`
	Track_number       string `json:"track_number" validate:"required,numeric"`
	Entry              string `json:"entry" validate:"required"`
	Delivery           []Delivery
	Payment            []Payment
	Items              []Item
	Locale             string    `json:"locale" validate:"required"`
	Internal_signature string    `json:"internal_signature" validate:"required"`
	Customer_id        int       `json:"customer_id" validate:"required,numeric"`
	Delivery_service   string    `json:"delivery_service" validate:"required"`
	Shardkey           int       `json:"shardkey" validate:"required,numeric"`
	Sm_id              int       `json:"sm_id" validate:"required,numeric"`
	Date_created       time.Time `json:"date_created"  validate:"required,datetime"`
	Oof_shard          int       `json:"oof_shard" validate:"required,numeric"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required,alpha"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     int    `json:"zip" validate:"required,numeric,lte=999999"`
	City    string `json:"city" validate:"required,alpha"`
	Address string `json:"address" validate:"required,alpha"`
	Region  string `json:"region" validate:"required,alpha"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction   string `json:"transaction" validate:"required,alphanum"`
	Request_id    int    `json:"request_id" validate:"required,numeric"`
	Currency      string `json:"currency" validate:"required,alpha"`
	Provider      string `json:"provider" validate:"required,alpha"`
	Amount        int    `json:"amount" validate:"required,numeric"`
	Payment_dt    int    `json:"payment_dt" validate:"required,numeric"`
	Bank          string `json:"bank" validate:"required,alpha"`
	Delivery_cost int    `json:"delivery_cost" validate:"required,numeric,gt=0"`
	Goods_total   int    `json:"goods_total" validate:"required,numeric,gt=0"`
	Custom_fee    int    `json:"custom_fee" validate:"required,numeric"`
}

type Item struct {
	Chrt_id      int    `json:"chrt_id" validate:"required,numeric"`
	Track_number string `json:"track_number" validate:"required,alpha"`
	Price        int    `json:"price" validate:"required,numeric,gt=0"`
	Rid          string `json:"rid" validate:"required,alphanum"`
	Name         string `json:"name" validate:"required,alpha"`
	Sale         int    `json:"sale" validate:"required,numeric,gt=0"`
	Size         int    `json:"size" validate:"required,numeric"`
	Total_price  int    `json:"total_price" validate:"required,numeric,gt=0"`
	Nm_id        int    `json:"nm_id" validate:"required,numeric,gte=0"`
	Brand        string `json:"brand" validate:"required,alpha"`
	Status       int    `json:"status" validate:"required,numeric"`
}
