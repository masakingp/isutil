package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

var db *sqlx.DB

func main() {
	ConnectDB()

	Sample()

	db.Close()
}


////////////////////////////////////////////////////////////////////////////////
// Sample
////////////////////////////////////////////////////////////////////////////////

type Chair struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Thumbnail   string `db:"thumbnail" json:"thumbnail"`
	Price       int64  `db:"price" json:"price"`
	Height      int64  `db:"height" json:"height"`
	Width       int64  `db:"width" json:"width"`
	Depth       int64  `db:"depth" json:"depth"`
	Color       string `db:"color" json:"color"`
	Features    string `db:"features" json:"features"`
	Kind        string `db:"kind" json:"kind"`
	Popularity  int64  `db:"popularity" json:"-"`
	PopularityDesc  int64  `db:"popularity_desc" json:"-"`
	Stock       int64  `db:"stock" json:"-"`
}

func Sample() {
	query := `SELECT * FROM chair WHERE id = ?`
	chair := Chair{}
	err := db.Get(&chair, query, 18938)
	if err != nil {
		panic(err)
	}

	log.Print(chair.Name)

	var chairs []Chair
	query = `SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT ?`
	err = db.Select(&chairs, query, 1)
	if err != nil {
		panic(err)
	}

	for _, c := range chairs {
		log.Print(c.Name)
	}
}



////////////////////////////////////////////////////////////////////////////////
// DB
////////////////////////////////////////////////////////////////////////////////

type MySQLConnectionEnv struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

var mySQLConnectionData *MySQLConnectionEnv

func NewMySQLConnectionEnv() *MySQLConnectionEnv {
	return &MySQLConnectionEnv{
		Host:     getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		User:     getEnv("MYSQL_USER", "isucon"),
		DBName:   getEnv("MYSQL_DBNAME", "isuumo"),
		Password: getEnv("MYSQL_PASS", "isucon"),
	}
}

func getEnv(key, defaultValue string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultValue
}

func (mc *MySQLConnectionEnv) ConnectDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?interpolateParams=true", mc.User, mc.Password, mc.Host, mc.Port, mc.DBName)
	return sqlx.Open("mysql", dsn)
}

func ConnectDB() {
	mySQLConnectionData = NewMySQLConnectionEnv()

	var err error
	db, err = mySQLConnectionData.ConnectDB()
	if err != nil {
		log.Printf("DB connection failed : %v", err)
	}
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxIdleConns(512)
	db.SetMaxOpenConns(512)
}

