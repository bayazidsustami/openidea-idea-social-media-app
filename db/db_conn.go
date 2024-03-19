package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

var (
	applicationDbName     = viper.GetString("DB_NAME")
	applicationDbUsername = viper.GetString("DB_USERNAME")
	applicationDbPassword = viper.GetString("DB_PASSWORD")
	applicationDbHost     = viper.GetString("DB_HOST")
	applicationDbPort     = viper.GetString("DB_PORT")
	maxConnLifeTime       = 60 * time.Minute
	maxConnIdleTime       = 5 * time.Minute
	maxConns              = int32(100)
	minConns              = int32(10)
)

func GetConnectionPool() *pgxpool.Pool {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", applicationDbUsername, applicationDbPassword, applicationDbHost, applicationDbPort, applicationDbName)
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	config.MaxConnLifetime = maxConnLifeTime
	config.MaxConnIdleTime = maxConnIdleTime
	config.MaxConns = maxConns
	config.MinConns = minConns
	config.ConnConfig.ConnectTimeout = 10 * time.Second

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return dbPool
}
