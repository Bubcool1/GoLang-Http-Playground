package config

import "os"

type contextKey string

const (
	DB_KEY contextKey = "db"
)

var ConnString = os.Getenv("DB_CONN_STRING")
