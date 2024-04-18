package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pgxpool_prometheus "github.com/cmackenzie1/pgxpool-prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

var (
	applicationDbName     = viper.GetString("DB_NAME")
	applicationDbUsername = viper.GetString("DB_USERNAME")
	applicationDbPassword = viper.GetString("DB_PASSWORD")
	applicationDbHost     = viper.GetString("DB_HOST")
	applicationDbPort     = viper.GetString("DB_PORT")
	applicationDbParams   = viper.GetString("DB_PARAMS")
	maxTimeout            = 10 * time.Second
	maxConnLifeTime       = 60 * time.Minute
	maxConnIdleTime       = 5 * time.Minute
	maxConns              = int32(100)
	minConns              = int32(0)
	dbPool                *pgxpool.Pool
	once                  sync.Once
)

func GetConnectionPool() *pgxpool.Pool {
	once.Do(func() {
		dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
			applicationDbUsername,
			applicationDbPassword,
			applicationDbHost,
			applicationDbPort,
			applicationDbName,
			applicationDbParams,
		)
		config, err := pgxpool.ParseConfig(dbUrl)
		if err != nil {
			log.Fatal(err)
		}

		config.MaxConnLifetime = maxConnLifeTime
		config.MaxConnIdleTime = maxConnIdleTime
		config.MaxConns = maxConns
		config.MinConns = minConns
		config.ConnConfig.ConnectTimeout = maxTimeout

		dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			log.Fatal(err)
		}

		err = dbPool.Ping(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		prometheus.MustRegister(pgxpool_prometheus.NewPgxPoolStatsCollector(dbPool, viper.GetString("DB_NAME")))
	})

	return dbPool
}
