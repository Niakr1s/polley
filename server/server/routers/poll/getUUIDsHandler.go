package poll

import (
	"encoding/json"
	"net/http"
	"net/url"
	"polley/server/helpers"
	"polley/server/storage.go"
)

type getUUIDsResponse struct {
	UUIDs []string `json:"uuids"`
	Total int      `json:"total"`
}

type getUUIDsRequest struct {
	Page     int
	PageSize int
}

func newGetUUIDsRequestFromQuery(query url.Values) getUUIDsRequest {
	defaultRequest := getUUIDsRequest{Page: 0, PageSize: 10}
	if pageSize, err := helpers.GetFirstIntValueFromQuery(query, "pageSize"); err == nil {
		defaultRequest.PageSize = pageSize
	}
	if page, err := helpers.GetFirstIntValueFromQuery(query, "page"); err == nil {
		defaultRequest.Page = page
	}
	return defaultRequest
}

func getUUIDs(storage *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := newGetUUIDsRequestFromQuery(r.URL.Query())

		uuids, err := storage.Polls.GetNPollsUUIDs(request.PageSize, request.Page)
		if err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}

		response := getUUIDsResponse{
			UUIDs: uuids,
			Total: storage.Polls.GetTotal(),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			helpers.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	}

}
