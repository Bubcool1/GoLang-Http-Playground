package handlers

import (
	"context"
	"log"
	"net/http"

	"beardsall.xyz/golangHttpPlayground/config"
	"beardsall.xyz/golangHttpPlayground/helpers"
	"beardsall.xyz/golangHttpPlayground/repository"
)

type import_audit struct {
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

func GetPaginatedAuditRows(ctx context.Context, req *http.Request) (any, error) {
	querySuffix := `LIMIT $1 OFFSET $2`

	pageNumberStr := req.URL.Query().Get("pageNumber")
	itemCountStr := req.URL.Query().Get("itemCount")

	PageNumber := 1
	ItemCount := config.ITEMS_PER_PAGE

	var err error

	if pageNumberStr != "" {
		PageNumber, err = helpers.SafeConvertToInt(pageNumberStr)
		if err != nil {
			return nil, err
		}
	}

	if itemCountStr != "" {
		ItemCount, err = helpers.SafeConvertToInt(itemCountStr)
		if err != nil {
			return nil, err
		}
	}

	offset, limit := helpers.CalculatePagination(PageNumber, ItemCount, true)

	return repository.ListRecords[import_audit](ctx, querySuffix, limit, offset)
}

func GetLatestAuditRow(ctx context.Context, req *http.Request) (any, error) {
	row, err := repository.GetRecord[import_audit](ctx)

	if err != nil {
		log.Printf("error fetching audit row: %v", err)
		return nil, err
	}
	return row, nil
}

// This should be mapped instead of static structs, then it can be extenisbile
