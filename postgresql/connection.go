package postgresql

import (
	"fmt"
	"realtime-batch-processing/config"

	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

var Postgres *gorm.DB

func Init() {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Config.PostgresqlUsername,
		config.Config.PostgresqlPassword,
		config.Config.PostgresqlHost,
		config.Config.PostgresqlPort,
		config.Config.PostgresqlDb)
	var err error
	Postgres, err = gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = Postgres.DB().Ping()
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	err := Postgres.Close()
	if err != nil {
		panic(err)
	}
}
