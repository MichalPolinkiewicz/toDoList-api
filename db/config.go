package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func init() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	charset := os.Getenv("charset")
	pt := os.Getenv("parseTime")
	dbUri := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=%s", username, password, dbName, charset, pt)

	conn, err := gorm.Open("mysql", dbUri)

	if err != nil {
		fmt.Print("ERROR:", err)
	}

	db = conn
	//db.Debug().AutoMigrate(&models.Task{})
}
