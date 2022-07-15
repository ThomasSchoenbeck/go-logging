package main

import "time"

type (
	ClientLogs struct {
		LOG_ID     int
		SESSION_ID string
		LOG_LEVEL  string
		URL        string
		// MSG        string
		MSG        interface{}
		STACKTRACE string
		TIMESTAMP  time.Time
		USERAGENT  string
		CLIENT_IP  string
		REMOTE_IP  string
	}

	filters struct {
		Field string      `json:"field"`
		Value interface{} `json:"value"`
	}

	logsPaginationRequest struct {
		Parameters paginationParameters `json:"parameters"`
		Filters    *[]filters           `json:"filters"`
		Sorting    *ClientLogs          `json:"sorting"`
	}

	paginationParameters struct {
		Limit       int `json:"limit"`
		CurrentPage int `json:"currentPage"`
	}

	PaginationResponse struct {
		Data               interface{} `json:"data"`
		NumRecords         int         `json:"numRecords"`
		NumFilteredRecords int         `json:"numFilteredRecords"`
		PageCount          int         `json:"pageCount"`
	}
)
