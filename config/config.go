package config

import "os"

type contextKey string

const (
	DB_KEY             contextKey = "db"
	PORT               string     = "8080"
	ITEMS_PER_PAGE     int        = 25
	MAX_ITEMS_PER_PAGE int        = 1000
)

var ConnString = os.Getenv("DB_CONN_STRING")
