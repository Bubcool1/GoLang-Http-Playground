package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"beardsall.xyz/golanghttpplayground/config"
	"beardsall.xyz/golanghttpplayground/helpers"
	"beardsall.xyz/golanghttpplayground/repository"
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
	req, offset, limit, err := helpers.GetPaginationDetails(req)
	if err != nil {
		return nil, err
	}
	// params := []repository.QueryFilter{}
	params, linkOperator := helpers.ExtractQueryParams(req)
	// raw_params, link_operator := helpers.ExtractQueryParams(req)
	// for param, value := range raw_params {
	// 	params = append(params, repository.QueryFilter{
	// 		FieldName: param,
	// 		Operator:  "=",
	// 		Value:     value.(string),
	// 	})
	// }
	log.Print(linkOperator)
	params = append(params, repository.QueryFilter{
		FieldName: config.LIMIT_PARAM,
		Operator:  "",
		Value:     strconv.Itoa(limit),
	})
	params = append(params, repository.QueryFilter{
		FieldName: config.OFFSET_PARAM,
		Operator:  "",
		Value:     strconv.Itoa(offset),
	})
	return repository.PaginatedListRecords[import_audit](ctx, params, linkOperator)
}

func GetLatestAuditRow(ctx context.Context, req *http.Request) (any, error) {
	// queryFilters := []repository.QueryFilter{
	// 	{FieldName: "id", Operator: "=", Value: "50"},
	// }
	row, err := repository.GetRecord[import_audit](ctx)
	// row, err := repository.GetRecord[import_audit](ctx, queryFilters...)

	if err != nil {
		log.Printf("error fetching audit row: %v", err)
		return nil, err
	}
	return row, nil
}

// This should be mapped instead of static structs, then it can be extenisbile
