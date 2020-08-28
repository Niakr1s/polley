package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"polley/db"
	"polley/models"
	"strconv"
)

type getNPollsHandlerResponse struct {
	Polls      []*models.Poll `json:"polls"`
	TotalPolls int            `json:"totalPolls"`
}

type getNPollsHandlerRequest struct {
	Page     int // default=0
	PageSize int // default=10
}

func newGetNPollsHandlerRequestDefault() getNPollsHandlerRequest {
	return getNPollsHandlerRequest{Page: 0, PageSize: 10}
}

func newGetNPollsHandlerRequestFromQuery(query url.Values) getNPollsHandlerRequest {
	request := newGetNPollsHandlerRequestDefault()
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

func (s *Server) getNPollsHandler(w http.ResponseWriter, r *http.Request) {
	request := newGetNPollsHandlerRequestFromQuery(r.URL.Query())

	polls, err := db.GetNPolls(s.pollDB, request.PageSize, request.Page)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}

	pollsResponses := make([]pollResponse, len(polls))
	for i, poll := range polls {
		pollsResponses[i] = newPollResponse(s, r, poll)
	}

	err = json.NewEncoder(w).Encode(pollsResponses)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
}
