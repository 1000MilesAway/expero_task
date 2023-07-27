package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var instance *sql.DB
var instanceOnce = &sync.Once{}

var (
	DatabaseUser     = "username"
	DatabasePassword = "password"
	DatabaseUrl      = "localhost"
	DatabasePort     = 5432
	DatabaseName     = "default_database"
)

// func GetDB() *pgxpool.Pool {
// 	instanceOnce.Do(func() {
// 		dbStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", DatabaseUser, DatabasePassword, DatabaseUrl, DatabasePort, DatabaseName)
// 		var err error
// 		for instance, err = pgxpool.Connect(context.Background(), dbStr); err != nil; instance, err = pgxpool.Connect(context.Background(), dbStr) {
// 			logrus.Errorf("Can not establish db connection: %s", err.Error())
// 			time.Sleep(12)
// 		}
// 	})
// 	return instance
// }

func GetDB() *sql.DB {
	instanceOnce.Do(func() {
		// dsn := fmt.Sprintf("%s:%s@%s:%d/%s", DatabaseUser, DatabasePassword, DatabaseUrl, DatabasePort, DatabaseName)
		dsn := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable", DatabaseUrl, DatabasePort, DatabaseUser, DatabasePassword, DatabaseName)

		var err error
		instance, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("Failed to connect to the database:", err)
		}
	})
	return instance
}
