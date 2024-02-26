package db

import (
	"fmt"
	"io/ioutil"
	"nats/api"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"
)

var (
	DBConnection *sqlx.DB
)

func NewDBConnection() (*sqlx.DB, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parent := filepath.Dir(wd)

	err = godotenv.Load(parent + "./.env")
	if err != nil {
		panic("Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")
	driver := os.Getenv("DRIVER")

	connUrl := fmt.Sprintf("host=%s port=%v user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect(driver, connUrl)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	slog.Info("Successfully connected to DB")
	return db, nil
}

func ExecMigration(db *sqlx.DB, path string) error {
	fmt.Println(os.Getwd())
	query, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		panic(err)
	}
	return nil
}

func TestItem(db *sqlx.DB) {
	var order api.Order
	/* 	err := db.Get(&order, `SELECT order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	   	json_build_array(json_build_object("name", deliveries.name, "phone", deliveries.phone, "zip", deliveries.zip, "city", deliveries.city,"address", deliveries.address, "region", deliveries.region, "email", deliveries.email)) as Delivery
	   	FROM orders
	   	JOIN deliveries ON order_uid = deliveries.order_id
	   	WHERE order_uid=$1
	   `, "76fdee04-72aa-47eb-8c85-24069af468fb") */

	/* if err != nil {
		slog.Error(fmt.Sprintf("Error: %s", err))
	} */
	//fmt.Printf("\n %v", order)
	err := db.Get(&order, `SELECT order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	jsonb_agg(jsonb_build_object('name', deliveries.name, 'phone', deliveries.phone, 'zip', deliveries.zip, 'city', deliveries.city,'address', deliveries.address, 'region', deliveries.region, 'email', deliveries.email)) as Delivery
	FROM orders
	JOIN deliveries ON order_uid = deliveries.order_id
	WHERE order_uid=$1
	GROUP BY order_uid
`, "32f2ed8d-f3cc-41e0-ab66-00426f515ffa")
	if err != nil {
		slog.Error(fmt.Sprintf("Error while query: %s", err))
		return
	}
	fmt.Println(order)
	return
}

func GetItembyId(db *sqlx.DB, id string) api.Order {
	/* statement := fmt.Sprintf(`select order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	array_agg(ARRAY[deliveries.name, deliveries.phone, deliveries.zip::varchar, deliveries.city, deliveries.address, deliveries.region, deliveries.email]) as Delivery,
	array_agg(ARRAY[payments.transaction, payments.request_id::varchar, payments.currency, payments.provider, payments.amount::varchar, payments.payment_dt::varchar, payments.bank, payments.delivery_cost::varchar, payments.goods_total::varchar, payments.custom_fee::varchar]) AS Payment,
	array_agg(ARRAY[items.chrt_id::varchar, items.track_number, items.price::varchar, items.rid::varchar, items.name, items.sale::varchar, items.size::varchar, items.total_price::varchar, items.nm_id::varchar, items.brand, items.status::varchar]) as Item
	FROM orders
	JOIN deliveries ON order_uid = deliveries.order_id
	JOIN payments ON order_uid = payments.transaction
	JOIN items ON orders.track_number = items.track_number
	WHERE order_uid = %s
	Group by order_uid`, id) */
	return api.Order{}
}

func GetItems(db *sqlx.DB) []api.Order {
	/* statement := `select order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	array_agg(ARRAY[deliveries.name, deliveries.phone, deliveries.zip::varchar, deliveries.city, deliveries.address, deliveries.region, deliveries.email]) as Delivery,
	array_agg(ARRAY[payments.transaction, payments.request_id::varchar, payments.currency, payments.provider, payments.amount::varchar, payments.payment_dt::varchar, payments.bank, payments.delivery_cost::varchar, payments.goods_total::varchar, payments.custom_fee::varchar]) AS Payment,
	array_agg(ARRAY[items.chrt_id::varchar, items.track_number, items.price::varchar, items.rid::varchar, items.name, items.sale::varchar, items.size::varchar, items.total_price::varchar, items.nm_id::varchar, items.brand, items.status::varchar]) as Item
	FROM orders
	JOIN deliveries ON order_uid = deliveries.order_id
	JOIN payments ON order_uid = payments.transaction
	JOIN items ON orders.track_number = items.track_number
	Group by order_uid` */
	return nil
}

func InsertItem(db *sqlx.DB, api_order api.Order) {
	tx := db.MustBegin()

	_, err := tx.NamedExec(`INSERT INTO orders (order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) VALUES 
	(:order_uid,:track_number,:entry,:locale,:internal_signature,:customer_id,:delivery_service,:shardkey,:sm_id,:date_created,:oof_shard)`, api_order)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while inserting in orders %s", err))
		tx.Rollback()
		return
	}

	deliveries, payments, items := prepareEntityToDB(api_order)

	_, err = tx.NamedExec(`INSERT INTO deliveries (order_id, name,phone,zip,city,address,region,email) VALUES 
	(:order_id,:name,:phone,:zip,:city,:address,:region,:email)`, deliveries)

	if err != nil {
		slog.Error(fmt.Sprintf("Error while inserting in deliveries %s", err))
		tx.Rollback()
		return
	}

	_, err = tx.NamedExec(`INSERT INTO payments (transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) VALUES 
	(:transaction,:request_id,:currency,:provider,:amount,:payment_dt,:bank,:delivery_cost,:goods_total,:custom_fee)`, payments)

	if err != nil {
		slog.Error(fmt.Sprintf("Error while inserting in payments %s", err))
		tx.Rollback()
		return
	}

	_, err = tx.NamedExec(`INSERT INTO items (chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) VALUES 
	(:chrt_id,:track_number,:price,:rid,:name,:sale,:size,:total_price,:nm_id,:brand,:status)`, items)

	if err != nil {
		slog.Error(fmt.Sprintf("Error while inserting in items %s", err))
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		slog.Error(fmt.Sprintf("Error while completing transaction %s", err))
	}

}

func prepareEntityToDB(order api.Order) ([]Delivery, []Payment, []Item) {
	deliveries := []Delivery{}
	payments := []Payment{}
	items := []Item{}

	for _, fields := range order.Delivery {
		delivery := Delivery{
			Order_id: order.Order_uid,
			Name:     fields.Name,
			Phone:    fields.Phone,
			Zip:      fields.Zip,
			City:     fields.City,
			Address:  fields.Address,
			Region:   fields.Region,
			Email:    fields.Email,
		}
		deliveries = append(deliveries, delivery)

	}

	for _, fields := range order.Payment {
		payment := Payment{
			Transaction:   order.Order_uid,
			Request_id:    fields.Request_id,
			Currency:      fields.Currency,
			Provider:      fields.Provider,
			Amount:        fields.Amount,
			Payment_dt:    fields.Payment_dt,
			Bank:          fields.Bank,
			Delivery_cost: fields.Delivery_cost,
			Goods_total:   fields.Goods_total,
			Custom_fee:    fields.Custom_fee,
		}
		payments = append(payments, payment)

	}

	for _, fields := range order.Item {
		item := Item{
			Chrt_id:      fields.Chrt_id,
			Track_number: order.Track_number,
			Price:        fields.Price,
			Rid:          fields.Rid,
			Name:         fields.Name,
			Sale:         fields.Sale,
			Size:         fields.Size,
			Total_price:  fields.Total_price,
			Nm_id:        fields.Nm_id,
			Brand:        fields.Brand,
			Status:       fields.Status,
		}
		items = append(items, item)

	}
	return deliveries, payments, items
}
