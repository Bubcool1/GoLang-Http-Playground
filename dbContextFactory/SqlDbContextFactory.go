package DbContextFactory

import (
	"context"
	"log"

	"beardsall.xyz/golanghttpplayground/config"
	"github.com/jmoiron/sqlx"
)

func SqlDbContextFactory(ctx context.Context) context.Context {
	log.Println("Connecting to database")
	db, err := sqlx.Connect("postgres", config.ConnString)
	log.Println("Connection Successful")
	if err != nil {
		log.Fatal(err)
	}
	ctx = context.WithValue(ctx, config.DB_KEY, db)
	return ctx
}
