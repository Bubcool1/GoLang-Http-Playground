package repository

import (
	"context"
	"log"
	"reflect"
	"strconv"
	"strings"

	"beardsall.xyz/golanghttpplayground/config"
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

	if config.SQL_LOGGING {
		log.Printf("Executing SQL Query: \n`%s`\n with args: %v", query, args)
	}

	if err := db.SelectContext(ctx, &rows, query, args...); err != nil {
		log.Printf("error fetching audit rows: %v", err)
		return nil, err
	}
	return rows, nil
}

// This should be mapped instead of static structs, then it can be extenisbile

func buildSqlQueryForType[T any](filters []QueryFilter, joinOperator string) (string, []any) {
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

	params := []any{}
	joinOperator = strings.Trim(joinOperator, "")
	joinOperator = " " + joinOperator + " "

	for filterIndex := range filters {
		if filters[filterIndex].FieldName == "" || filters[filterIndex].Operator == "" {
			log.Printf("Invalid filter at index %d: %+v", filterIndex, filters[filterIndex])
			continue
		}
		if filterIndex == 0 {
			queryString += " WHERE "
		} else {
			queryString += joinOperator
		}
		queryString += filters[filterIndex].FieldName + " " + filters[filterIndex].Operator + " $" + strconv.Itoa(filterIndex+1)
		params = append(params, filters[filterIndex].Value)
	}

	return queryString, params
}

type QueryFilter struct {
	FieldName string
	Operator  string
	Value     string
}

func GetRecord[T any](ctx context.Context, filters ...QueryFilter) (*T, error) {
	queryString, params := buildSqlQueryForType[T](filters, "AND")

	queryString += " LIMIT 1"

	return GetRecordFromQuery[T](ctx, queryString, params...)
}

func ListRecords[T any](ctx context.Context, querySuffix string, filters ...QueryFilter) ([]T, error) {
	queryString, args := buildSqlQueryForType[T](filters, "AND")

	if !strings.HasPrefix(querySuffix, " ") {
		querySuffix = " " + querySuffix
	}

	queryString += querySuffix

	return ListRecordsFromQuery[T](ctx, queryString, args...)
}

func PaginatedListRecordsAdvanced[T any](ctx context.Context, querySuffix string, filters []QueryFilter, linkOperator string) ([]T, error) {
	FilterLimit := strconv.Itoa(config.ITEMS_PER_PAGE)
	FilterOffset := "0"
	for filterIndex := range filters {
		switch filters[filterIndex].FieldName {
		case config.LIMIT_PARAM:
			FilterLimit = filters[filterIndex].Value
		case config.OFFSET_PARAM:
			FilterOffset = filters[filterIndex].Value
		}
	}
	queryString, args := buildSqlQueryForType[T](filters, linkOperator)

	if !strings.HasPrefix(querySuffix, " ") {
		querySuffix = " " + querySuffix
	}

	queryString += querySuffix

	// Here is the problem, args turns into a list of strings, when it needs to be a list of ints, either that or everything needs to be string, but I think that would break for bools etc, unless we covert bools to 0/1
	args = append(args, FilterLimit)
	args = append(args, FilterOffset)

	// len args - 2 thats where we start limit and offset
	args_len := len(args)
	queryString += " LIMIT $" + strconv.Itoa(args_len-1) + " OFFSET $" + strconv.Itoa(args_len)

	return ListRecordsFromQuery[T](ctx, queryString, args...)
}

func PaginatedListRecords[T any](ctx context.Context, filters []QueryFilter, linkOperator string) ([]T, error) {
	queryStringPlaceholder := ""
	return PaginatedListRecordsAdvanced[T](ctx, queryStringPlaceholder, filters, linkOperator)
}
