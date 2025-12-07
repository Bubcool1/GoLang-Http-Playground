package main

import (
	"context"
	"log"
	"net/http"

	"beardsall.xyz/golangHttpPlayground/config"
	DbContextFactory "beardsall.xyz/golangHttpPlayground/dbContextFactory"
	ResponseHandler "beardsall.xyz/golangHttpPlayground/helpers"
	_ "github.com/lib/pq"
)

func setup() context.Context {
	ctx := context.Background()
	ctx = DbContextFactory.SqlDbContextFactory(ctx)

	return ctx
}

func main() {
	println("Starting server on port 8080")
	ctx := setup()

	for route, handler := range Routes {
		http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
			println("Received request for route: " + route)
			handlerRes, err := handler(ctx, req)

			if err != nil {
				// This needs to return the value and status code properly instead of just panik 500 errors
				ErrVal := map[string]any{
					"error":             "Something went wrong",
					"statusCode":        500,
					"additionalDetails": err.Error(),
				}
				ResponseHandler.FormatResponse(w, ErrVal, req)
			}
			ResponseHandler.FormatResponse(w, handlerRes, req)
		})
	}

	log.Fatal(http.ListenAndServe(":"+config.PORT, nil))
}

// TODO: Reformat error handling. Should be a tuple response response, error across all funcs and the interface
