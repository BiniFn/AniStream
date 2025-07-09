package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func parsePagination(r *http.Request, defaultPage, defaultSize int) (page, size int, err error) {
	q := r.URL.Query()

	page = defaultPage
	if v := q.Get("page"); v != "" {
		pageVal, err2 := strconv.Atoi(v)
		if err2 != nil || pageVal < 1 {
			return 0, 0, fmt.Errorf("invalid page")
		}
		page = pageVal
	}

	size = defaultSize
	if v := q.Get("itemsPerPage"); v != "" {
		sizeVal, err2 := strconv.Atoi(v)
		if err2 != nil || sizeVal < 1 {
			return 0, 0, fmt.Errorf("invalid itemsPerPage")
		}
		size = sizeVal
	}

	return page, size, nil
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func jsonOK(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
