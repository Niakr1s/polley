package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type getUUIDsResponse struct {
	UUIDs []string `json:"uuids"`
	Total int      `json:"total"`
}

type getUUIDsRequest struct {
	Page     int // default=0
	PageSize int // default=10
}

func newGetUUIDsRequestDefault() getUUIDsRequest {
	return getUUIDsRequest{Page: 0, PageSize: 10}
}

func newGetUUIDsRequestFromQuery(query url.Values) getUUIDsRequest {
	request := newGetUUIDsRequestDefault()
	if pageSize, err := getFirstIntValueFromQuery(query, "pageSize"); err == nil {
		request.PageSize = pageSize
	}
	if page, err := getFirstIntValueFromQuery(query, "page"); err == nil {
		request.Page = page
	}
	return request
}

func getFirstIntValueFromQuery(query url.Values, key string) (int, error) {
	values, ok := query[key]
	if !ok || len(values) == 0 {
		return 0, fmt.Errorf("no value exists for key=%s", key)
	}
	return strconv.Atoi(values[0])
}

func (s *Server) getUUIDsHandler(w http.ResponseWriter, r *http.Request) {
	request := newGetUUIDsRequestFromQuery(r.URL.Query())

	uuids, err := s.pollDB.GetNPollsUUIDs(request.PageSize, request.Page)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	response := getUUIDsResponse{
		UUIDs: uuids,
		Total: s.pollDB.GetTotal(),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
}
