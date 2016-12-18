package infra

import (
	"fmt"
	"realtime-batch-processing/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MySQL *gorm.DB

func Init() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		config.Config.MySqlUsername,
		config.Config.MySqlPassword,
		config.Config.MySqlHost,
		config.Config.MySqlPort,
		config.Config.MySqlDb)
	var err error
	fmt.Println(connectionString)
	MySQL, err = gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	err = MySQL.DB().Ping()
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	err := MySQL.Close()
	if err != nil {
		panic(err)
	}
}
