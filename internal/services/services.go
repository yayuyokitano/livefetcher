package services

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var IsTesting bool

func Start() (err error) {
	if os.Getenv("TESTING") == "true" {
		IsTesting = true
	}

	var domain string
	if os.Getenv("CONTAINERIZED") == "true" {
		domain = "db"
	} else {
		domain = "localhost"
	}

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?pool_max_conns=100",
		os.Getenv("POSTGRES_USER"), url.QueryEscape(os.Getenv("POSTGRES_PASSWORD")), domain, os.Getenv("POSTGRES_DB"))

	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return
	}
	return
}

func Stop() {
	if Pool != nil {
		Pool.Close()
	}
}
