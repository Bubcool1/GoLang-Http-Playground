package repository

import (
	"context"
	"log"
	"reflect"
	"strings"

	"beardsall.xyz/golangHttpPlayground/config"
	"github.com/jmoiron/sqlx"
)

func GetRecordFromQuery[T any](ctx context.Context, query string, args ...any) (*T, error) {
	db := ctx.Value(config.DB_KEY).(*sqlx.DB)

	var row T

	if err := db.GetContext(ctx, &row, query, args...); err != nil {
		log.Printf("error fetching audit rows: %v", err)
		return nil, err
	}
	return &row, nil
}

func ListRecordsFromQuery[T any](ctx context.Context, query string, args ...any) ([]T, error) {
	db := ctx.Value(config.DB_KEY).(*sqlx.DB)

	var rows []T

	if err := db.SelectContext(ctx, &rows, query, args...); err != nil {
		log.Printf("error fetching audit rows: %v", err)
		return nil, err
	}
	return rows, nil
}

// This should be mapped instead of static structs, then it can be extenisbile

func buildSqlQueryForType[T any]() string {
	var row T

	rowType := reflect.TypeOf(row)
	// If it's a pointer, get the element type
	if rowType.Kind() == reflect.Ptr {
		rowType = rowType.Elem()
	}

	columns := make([]string, 0)

	for i := 0; i < rowType.NumField(); i++ {
		field := rowType.Field(i)

		// Use db tag if present, otherwise use field name
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			columns = append(columns, dbTag)
		} else {
			columns = append(columns, field.Name)
		}
	}

	queryString := "SELECT "

	for i, colName := range columns {
		queryString += colName
		if i < len(columns)-1 {
			queryString += ", "
		}
	}

	queryString += " FROM " + rowType.Name() // TODO: Replace with actual table name

	return queryString
}

func GetRecord[T any](ctx context.Context, args ...any) (*T, error) {
	queryString := buildSqlQueryForType[T]()
	queryString += " LIMIT 1"

	return GetRecordFromQuery[T](ctx, queryString, args...)
}

func ListRecords[T any](ctx context.Context, querySuffix string, args ...any) ([]T, error) {
	queryString := buildSqlQueryForType[T]()

	if !strings.HasPrefix(querySuffix, " ") {
		querySuffix = " " + querySuffix
	}

	queryString += querySuffix

	return ListRecordsFromQuery[T](ctx, queryString, args...)
}
