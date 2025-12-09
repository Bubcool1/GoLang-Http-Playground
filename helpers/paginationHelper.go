package helpers

import (
	"net/http"

	"beardsall.xyz/golanghttpplayground/config"
)

func GetPaginationDetails(req *http.Request) (*http.Request, int, int, error) {
	pageNumberStr := req.URL.Query().Get("pageNumber")
	if req.URL.Query().Has("pageNumber") {
		delete(req.URL.Query(), "pageNumber")
	}

	itemCountStr := req.URL.Query().Get("pageSize")
	if req.URL.Query().Has("pageSize") {
		delete(req.URL.Query(), "pageSize")
	}

	PageNumber := 1
	ItemCount := config.ITEMS_PER_PAGE

	var err error

	if pageNumberStr != "" {
		PageNumber, err = SafeConvertToInt(pageNumberStr)
		if err != nil {
			return req, 0, 0, err
		}
	}

	if itemCountStr != "" {
		ItemCount, err = SafeConvertToInt(itemCountStr)
		if err != nil {
			return req, 0, 0, err
		}
	}

	offset, limit := CalculatePagination(PageNumber, ItemCount, true)
	return req, offset, limit, nil
}
