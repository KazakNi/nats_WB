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

func getItembyId(db *sqlx.DB, id string) api.Order {
	// from controller
}

func getItems(db *sqlx.DB) []api.Order {
	// from Cache
}

func insertItem(db *sqlx.DB, order api.Order) {
	// from Pub
}
