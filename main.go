package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func setup() context.Context {
	ctx := context.Background()
	ctx = SqlDbContextFactory(ctx)

	return ctx
}

func main() {
	println("Starting server on port 8080")
	ctx := setup()

	for route, handler := range Routes {
		http.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
			println("Received request for route: " + route)
			handlerRes := handler(ctx, req)

			if handlerRes == nil {
				FormatResponse(w, map[string]any{
					"error":      "Something went wrong",
					"statusCode": 500,
				}, req)
			}
			FormatResponse(w, handlerRes, req)
		})
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO: Reformat error handling. Should be a tuple response response, error across all funcs and the interface
