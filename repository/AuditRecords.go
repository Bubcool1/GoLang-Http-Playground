package repository

import (
	"context"
	"log"

	"beardsall.xyz/golangHttpPlayground/config"
	"github.com/jmoiron/sqlx"
)

type ImportAudit struct {
	Id               int     `db:"id"`
	Filename         string  `db:"filename"`
	FilePath         string  `db:"file_path"`
	FileSize         *int    `db:"file_size"`
	FileModifiedDate *string `db:"file_modified_date"`
	ImportStartTime  string  `db:"import_start_time"`
	ImportEndTime    *string `db:"import_end_time"`
	RowCount         *int    `db:"row_count"`
	Status           *string `db:"status"`
	ErrorMessage     *string `db:"error_message"`
	TableName        *string `db:"table_name"`
	SchemaVersion    *string `db:"schema_version"`
	CreatedAt        *string `db:"created_at"`
}

func GetSingleRecord(ctx context.Context, query string) (*any, error) {
	db := ctx.Value(config.DB_KEY).(*sqlx.DB)

	var row any

	if err := db.SelectContext(ctx, &row, query); err != nil {
		log.Printf("error fetching audit rows: %v", err)
		return nil, err
	}
	return &row, nil
}

// This should be mapped instead of static structs, then it can be extenisbile
