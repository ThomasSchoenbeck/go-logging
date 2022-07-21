package main

type (
	Application struct {
		APP_ID   string
		APP_NAME string
		APP_URL  *string
		APP_DESC *string
		APP_LOGO []byte
	}

	Feedback_Channel struct {
		CHANNEL_ID       int
		APP_ID           string
		CHANNEL_NAME     string
		CHANNEL_DESC     string
		CHANNEL_ENDPOINT string
	}

	Feedback struct {
		CHANNEL_ID                 int
		FEEDBACK_ID                int
		APP_ID                     string
		FEEDBACK_TITLE             string
		FEEDBACK_MESSAGE           string
		FEEDBACK_POSITIVE_NEGATIVE *bool
		FEEDBACK_RAITING           *string
	}

	ClientLogs struct {
		LOG_ID     int
		APP_ID     string
		SESSION_ID string
		LOG_LEVEL  string
		URL        string
		// MSG        string
		MSG        interface{}
		STACKTRACE string
		// TIMESTAMP  time.Time
		TIMESTAMP string
		USERAGENT string
		CLIENT_IP string
		REMOTE_IP string
	}

	filters struct {
		Field string      `json:"field"`
		Value interface{} `json:"value"`
	}

	sorting struct {
		Field         string
		SortDirection string
	}

	PaginationRequest struct {
		Parameters paginationParameters `json:"parameters"`
		Filters    *[]filters           `json:"filters"`
		Sorting    *[]sorting           `json:"sorting"`
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
