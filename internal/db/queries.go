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

func GetItembyId(db *sqlx.DB, id string) api.Order {
	var order api.Order
	err := db.Get(&order, `SELECT order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	json_agg(json_build_object('name', deliveries.name, 'phone', deliveries.phone, 'zip', deliveries.zip, 'city', deliveries.city,'address', deliveries.address, 'region', deliveries.region, 'email', deliveries.email)) as Delivery,
	json_agg(json_build_object('transaction', payments.transaction, 'request_id', payments.request_id, 'currency', payments.currency, 'provider', payments.provider, 'amount', payments.amount, 'payment_dt', payments.payment_dt, 'bank', payments.bank, 'delivery_cost', payments.delivery_cost, 'goods_total', payments.goods_total, 'custom_fee', payments.custom_fee)) as Payment,
	json_agg(json_build_object('chrt_id', items.chrt_id, 'track_number', items.track_number, 'price', items.price, 'rid', items.rid, 'name', items.name, 'sale', items.sale, 'size', items.size, 'total_price', items.total_price, 'nm_id', items.nm_id, 'brand', items.brand, 'status', items.status)) as Item
	FROM orders
	JOIN deliveries ON order_uid = deliveries.order_id
	JOIN payments ON order_uid = payments.transaction
	JOIN items ON orders.track_number = items.track_number
	WHERE order_uid=$1
	GROUP BY order_uid
`, id)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while query: %s", err))
		return order
	}
	return order
}

func GetItems(db *sqlx.DB) []api.Order {
	var order []api.Order
	err := db.Select(&order, `SELECT order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	json_agg(json_build_object('name', deliveries.name, 'phone', deliveries.phone, 'zip', deliveries.zip, 'city', deliveries.city,'address', deliveries.address, 'region', deliveries.region, 'email', deliveries.email)) as Delivery,
	json_agg(json_build_object('transaction', payments.transaction, 'request_id', payments.request_id, 'currency', payments.currency, 'provider', payments.provider, 'amount', payments.amount, 'payment_dt', payments.payment_dt, 'bank', payments.bank, 'delivery_cost', payments.delivery_cost, 'goods_total', payments.goods_total, 'custom_fee', payments.custom_fee)) as Payment,
	json_agg(json_build_object('chrt_id', items.chrt_id, 'track_number', items.track_number, 'price', items.price, 'rid', items.rid, 'name', items.name, 'sale', items.sale, 'size', items.size, 'total_price', items.total_price, 'nm_id', items.nm_id, 'brand', items.brand, 'status', items.status)) as Item
	FROM orders
	JOIN deliveries ON order_uid = deliveries.order_id
	JOIN payments ON order_uid = payments.transaction
	JOIN items ON orders.track_number = items.track_number
	GROUP BY order_uid`)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while query: %s", err))
		return order
	}
	return order
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
