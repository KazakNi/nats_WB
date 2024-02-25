package services

import (
	"encoding/json"
	"nats/api"

	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Order struct {
	Order_uid          string     `fake:"{uuid}"`
	Track_number       string     `fake:"{numerify:#######}"`
	Entry              string     `fake:"{hackerabbreviation}"`
	Delivery           []Delivery `fakesize:"1"`
	Payment            []Payment  `fakesize:"1"`
	Item               []Item     `fakesize:"1"`
	Locale             string     `fake:"{countryabr}"`
	Internal_signature string     `fake:"{buzzword}"`
	Customer_id        int        `fake:"{number:1,1000000}"`
	Delivery_service   string     `fake:"{company}"`
	Shardkey           int        `fake:"{number:1,9}"`
	Sm_id              int        `fake:"{number:1,99}"`
	Date_created       time.Time  `fake:"{date}"`
	Oof_shard          int        `fake:"{number:1,9}"`
}

type Delivery struct {
	Name    string `fake:"{firstname}"`
	Phone   string `fake:"{phoneformatted}"`
	Zip     int    `fake:"{number:000000,999999}"`
	City    string `fake:"{city}"`
	Address string `fake:"Wildberries{streetsuffix}"`
	Region  string `fake:"Region{streetsuffix}"`
	Email   string `fake:"{email}"`
}

type Payment struct {
	Transaction   string `fake:"skip"`
	Request_id    int    `fake:"{number:1,9999}"`
	Currency      string `fake:"{currencyshort}"`
	Provider      string `fake:"{hackernoun}"`
	Amount        int    `fake:"{number:1,999}"`
	Payment_dt    int    `fake:"{number:1,1000000}"`
	Bank          string `fake:"{hackerabbreviation}"`
	Delivery_cost int    `fake:"{number:1,5000}"`
	Goods_total   int    `fake:"{number:1,999}"`
	Custom_fee    int    `fake:"{number:1,99}"`
}

type Item struct {
	Chrt_id      int    `fake:"{number:1,1000000}"`
	Track_number string `fake:"{buzzword}"`
	Price        int    `fake:"{number:1,1000000}"`
	Rid          string `fake:"{uuid}"`
	Name         string `fake:"{name}"`
	Sale         int    `fake:"{number:1,99}"`
	Size         int    `fake:"{number:1,99}"`
	Total_price  int    `fake:"{number:1,1000000}"`
	Nm_id        int    `fake:"{number:1,1000000}"`
	Brand        string `fake:"{company}"`
	Status       int    `fake:"{number:100,599}"`
}

func Json_generator(i int) ([]byte, error) {
	gofakeit.Seed(i)

	var fakeorder Order
	var b []byte
	var order api.Order

	gofakeit.Struct(&fakeorder)

	b, err := json.MarshalIndent(&fakeorder, "", "    ")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &order)

	if err != nil {
		return nil, err
	}

	b, err = json.MarshalIndent(&order, "", "    ")

	if err != nil {
		return nil, err
	}

	return b, nil
}
