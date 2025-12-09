package config

import "os"

type contextKey string

const (
	DB_KEY              contextKey = "db"
	PORT                string     = "8080"
	ITEMS_PER_PAGE      int        = 25
	MAX_ITEMS_PER_PAGE  int        = 1000
	PAGE_NUMBER_PARAM   string     = "pageNumber"
	PAGE_SIZE_PARAM     string     = "pageSize"
	SQL_LOGGING         bool       = true
	LIMIT_PARAM         string     = "limit"
	OFFSET_PARAM        string     = "offset"
	DEFAULT_OPERATOR    string     = "="
	OPERATOR_SUFFIX     string     = "Operator"
	LINK_OPERATOR_PARAM string     = "linkOperator"
)

var RESERVED_PARAMS = []string{PAGE_NUMBER_PARAM, PAGE_SIZE_PARAM, "operator"}

// var RESERVED_PARAMS = []string{"operator"}

var ConnString = os.Getenv("DB_CONN_STRING")
