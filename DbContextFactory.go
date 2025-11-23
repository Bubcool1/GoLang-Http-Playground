package main

import (
	"context"
	"log"

	"beardsall.xyz/golangHttpPlayground/config"
	"github.com/jmoiron/sqlx"
)

func SqlDbContextFactory(ctx context.Context) context.Context {
	println("Connecting to database")
	db, err := sqlx.Connect("postgres", config.ConnString)
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, config.DB_KEY, db)
	return ctx
}
