package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var instance *pgxpool.Pool
var instanceOnce = &sync.Once{}

var (
	DatabaseUser     = "username"
	DatabasePassword = "password"
	DatabaseUrl      = "localhost"
	DatabasePort     = 5432
	DatabaseName     = "default_database"
)

func GetDB() *pgxpool.Pool {
	instanceOnce.Do(func() {
		dbStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", DatabaseUser, DatabasePassword, DatabaseUrl, DatabasePort, DatabaseName)
		var err error
		for instance, err = pgxpool.Connect(context.Background(), dbStr); err != nil; instance, err = pgxpool.Connect(context.Background(), dbStr) {
			logrus.Errorf("Can not establish db connection: %s", err.Error())
			time.Sleep(12)
		}
	})
	return instance
}
