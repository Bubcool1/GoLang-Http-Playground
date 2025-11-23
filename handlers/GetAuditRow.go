package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

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

func GetLatestAuditRow(ctx context.Context, req *http.Request) any {
	db := ctx.Value(config.DB_KEY).(*sqlx.DB)

	// Parse and validate itemCount
	itemCountStr := req.URL.Query().Get("itemCount")

	itemCount := 1
	if itemCountStr != "" {
		parsed, err := strconv.Atoi(itemCountStr)
		if err != nil || parsed < 1 {
			log.Printf("invalid itemCount: %s", itemCountStr)
			return nil // or return an error response
		}
		itemCount = parsed
	}

	// Cap the limit to prevent abuse
	if itemCount > 100 {
		itemCount = 100
	}

	query := `SELECT id, filename, file_path, file_size, file_modified_date, 
			  import_start_time, import_end_time, row_count, status, 
			  error_message, table_name, schema_version, created_at 
			  FROM import_audit ORDER BY id DESC LIMIT $1`

	if itemCount == 1 {
		var row ImportAudit
		if err := db.GetContext(ctx, &row, query, itemCount); err != nil {
			log.Printf("error fetching audit row: %v", err)
			return nil
		}

		return row
	}

	var rows []ImportAudit
	if err := db.SelectContext(ctx, &rows, query, itemCount); err != nil {
		log.Printf("error fetching audit rows: %v", err)
		return nil
	}
	return rows
}

// This should be mapped instead of static structs, then it can be extenisbile
