package main

import (
	"context"
	"net/http"

	"beardsall.xyz/golangHttpPlayground/handlers"
)

type HttpRequestHandler func(ctx context.Context, req *http.Request) any

var Routes = map[string]HttpRequestHandler{
	"/auditLatest": handlers.GetLatestAuditRow,
}
