package infra

import (
	"fmt"
	"realtime-batch-processing/config"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var Postgresql *gorm.DB

func Init() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		config.Config.PostgresqlUsername,
		config.Config.PostgresqlPassword,
		config.Config.PostgresqlHost,
		config.Config.PostgresqlPort,
		config.Config.PostgresqlDb)
	var err error
	fmt.Println(connectionString)
	Postgresql, err = gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = Postgresql.DB().Ping()
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	err := Postgresql.Close()
	if err != nil {
		panic(err)
	}
}
