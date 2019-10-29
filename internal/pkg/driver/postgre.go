package driver

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gopkg.in/gorp.v2"
)

type DBPostgreOption struct {
	Host        string
	Port        int
	Username    string
	Password    string
	DBName      string
	MaxPoolSize int
}

func NewPostgreDatabase(option DBPostgreOption) *gorp.DbMap {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", option.Host, option.Port, option.Username, option.DBName, option.Password))
	if err != nil {
		panic(fmt.Errorf("ERROR connect to DB Postgre: %s | %v", option.DBName, err))
	}

	db.SetMaxOpenConns(option.MaxPoolSize)
	gorp := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	return gorp
}
