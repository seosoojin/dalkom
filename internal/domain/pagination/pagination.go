package pagination

import (
	"net/http"
	"strconv"
)

type Page struct {
	Offset int
	Limit  int
}

func NewPageFromRequest(r *http.Request) (Page, error) {
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	if offsetStr == "" {
		offsetStr = "0"
	}

	if limitStr == "" {
		limitStr = "-1"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return Page{}, err
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return Page{}, err
	}

	return Page{
		Offset: offset,
		Limit:  limit,
	}, nil
}
