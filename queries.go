package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"
)

func getLogsPaginated(ctx context.Context, lpr PaginationRequest, appID string) (*PaginationResponse, error) {

	var cls []ClientLogs

	var whereClause []string
	var numTotalRecords int
	var numFilteredRecords int
	var sqlNumOfRecords = "SELECT COUNT(*) from TSC_CLIENT_LOGS "

	var pr PaginationResponse

	// get total number of Records
	err := db.QueryRowContext(ctx, sqlNumOfRecords).Scan(&numTotalRecords)
	if err != nil {
		log.Println("error collecting number of Total Records for Client Logs", err)
		return nil, err
	}
	pr.NumRecords = numTotalRecords

	if lpr.Filters != nil {

		for _, f := range *lpr.Filters {

			if f.Field == "SESSION_ID" {
				whereClause = append(whereClause, fmt.Sprintf("lower(%s) LIKE '%%%s%%'", f.Field, strings.ToLower(fmt.Sprintf("%v", f.Value))))
			}

			if f.Field == "LOG_LEVEL" {
				whereClause = append(whereClause, fmt.Sprintf("lower(%s) LIKE '%%%s%%'", f.Field, strings.ToLower(fmt.Sprintf("%v", f.Value))))
			}

			// log.Println("filter", f, f.Value)

			// if f.Value == reflect.TypeOf("string") {
			// 	whereClause = append(whereClause, fmt.Sprintf("lower(%s) LIKE '%%%s%%'", f.Field, strings.ToLower(fmt.Sprintf("%v", f.Value))))
			// }

			// if f.Value == reflect.TypeOf("int") {
			// 	whereClause = append(whereClause, fmt.Sprintf("%s LIKE '%%%s%%'", f.Field, fmt.Sprintf("%v", f.Value)))
			// }

		}

		// filters := reflect.TypeOf(lpr.Filters)
		// filterValues := reflect.ValueOf(lpr.Filters)

		// // for i := 0; i < filters.NumField(); i++ {
		// // 	fieldName := filters.Field(i).Name
		// // 	fieldType := filters.Field(i).Type
		// // 	log.Println("i", i, ", Name:", fieldName, ", Type:", fieldType, ", value:", filterValues.Field(i).Interface())

		// // 	if fieldType == reflect.TypeOf("string") {
		// // 		whereClause = append(whereClause, fmt.Sprintf("lower(%s) LIKE '%%%s%%'", fieldName, strings.ToLower(fmt.Sprintf("%v", filterValues.Field(i).Interface()))))
		// // 	}

		// // 	if fieldType == reflect.TypeOf("int") {
		// // 		whereClause = append(whereClause, fmt.Sprintf("%s LIKE '%%%s%%'", fieldName, fmt.Sprintf("%v", filterValues.Field(i).Interface())))
		// // 	}

		// // }
	}
	sqlStr := fmt.Sprintf(`SELECT
  		LOG_ID
		, APP_ID
  	, SESSION_ID
		, LOG_LEVEL
		, URL
		, MSG
		, STACKTRACE
		, TIMESTAMP
		, USERAGENT
		, CLIENT_IP
		, REMOTE_IP
	FROM TSC_CLIENT_LOGS
	`)

	var sqlNumFilteredRecords string = sqlNumOfRecords
	whereClause = append(whereClause, fmt.Sprintf("APP_ID = '%s'", appID))
	if len(whereClause) > 0 {
		sqlStr += fmt.Sprintf("where %s", strings.Join(whereClause, " AND "))
		sqlNumFilteredRecords += fmt.Sprintf("where %s", strings.Join(whereClause, " AND "))
	}

	log.Println("sqlNumFilteredRecords", sqlNumFilteredRecords)

	err = db.QueryRowContext(ctx, sqlNumFilteredRecords).Scan(&numFilteredRecords)
	if err != nil {
		log.Println("error counting number of filtered records", err)
		return nil, err
	}
	pr.NumFilteredRecords = numFilteredRecords

	sqlStr += " ORDER BY LOG_ID"

	if lpr.Parameters.Limit > 0 {
		pr.PageCount = int(math.Max(float64(numFilteredRecords/lpr.Parameters.Limit), float64(1)))
		sqlStr += fmt.Sprintf(" LIMIT %d", lpr.Parameters.Limit)
		if lpr.Parameters.CurrentPage > 1 {
			sqlStr += fmt.Sprintf(" OFFSET %d", lpr.Parameters.Limit*(lpr.Parameters.CurrentPage-1))
		}
	}

	log.Println("sqlStr", sqlStr)

	rows, err := db.QueryContext(ctx, sqlStr)
	if err != nil {
		log.Println("Error query Client Logs paginated", err)
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var l ClientLogs

		//For CLOB Column
		// b := new(bytes.Buffer)
		// var lobRuleJSON driver.NullLob
		// lobRuleJSON.Lob = new(driver.Lob)
		// lobRuleJSON.Lob.SetWriter(b)

		// if err := rows.Scan(&l.LOG_ID, &l.SESSION_ID, &l.LOG_LEVEL, &l.URL, &lobRuleJSON, &l.STACKTRACE, &l.TIMESTAMP, &l.USERAGENT, &l.CLIENT_IP, &l.REMOTE_IP); err != nil {
		if err := rows.Scan(&l.LOG_ID, &l.APP_ID, &l.SESSION_ID, &l.LOG_LEVEL, &l.URL, &l.MSG, &l.STACKTRACE, &l.TIMESTAMP, &l.USERAGENT, &l.CLIENT_IP, &l.REMOTE_IP); err != nil {
			log.Println("error scanning Client Logs, index:", i, err)
			return nil, err
		}

		// if lobRuleJSON.Valid { // only valid if data from DB is not NULL
		// l.MSG = string(b.Bytes())
		// }

		cls = append(cls, l)
		i++
	}

	pr.Data = cls

	return &pr, nil
}

func getSQLPaginationFilters(filters *[]filters) []string {

	var whereClause []string

	if filters != nil {

		for _, f := range *filters {

			log.Println("filter", f, f.Value)

			if f.Value == reflect.TypeOf("string") {
				whereClause = append(whereClause, fmt.Sprintf("lower(%s) LIKE '%%%s%%'", f.Field, strings.ToLower(fmt.Sprintf("%v", f.Value))))
			}

			if f.Value == reflect.TypeOf("int") {
				whereClause = append(whereClause, fmt.Sprintf("%s LIKE '%%%s%%'", f.Field, fmt.Sprintf("%v", f.Value)))
			}

		}

	}

	return whereClause
}

func getTotalAmountOfRecords(ctx context.Context, sqlNumOfRecords string) (int, error) {

	var numTotalRecords int
	// get total number of Records
	err := db.QueryRowContext(ctx, sqlNumOfRecords).Scan(&numTotalRecords)
	if err != nil {
		log.Println("error collecting number of Total Records for apps", err)
		return -1, err
	}

	return numTotalRecords, nil
}

func createApp(ctx context.Context, a Application) (*Application, error) {

	newAppId := NewID()

	_, err := db.ExecContext(ctx, "INSERT INTO TSC_APPLICATIONS (APP_ID, APP_NAME) VALUES (?,?)", newAppId, a.APP_NAME)
	if err != nil {
		log.Println("error creating app with id", newAppId, err)
		return nil, err
	}

	err = db.QueryRowContext(ctx, `SELECT
		APP_ID
	, APP_NAME
	, APP_URL
	, APP_DESC
	, APP_LOGO
	, INSERT_TS
	, UPDATE_TS
	FROM TSC_APPLICATIONS WHERE APP_ID = ?`, newAppId).Scan(&a.APP_ID, &a.APP_NAME, &a.APP_URL, &a.APP_DESC, &a.APP_LOGO, &a.INSERT_TS, &a.UPDATE_TS)
	if err != nil {
		log.Println("error retrieving newly created app with id", newAppId, err)
		return nil, err
	}

	return &a, nil
}

func getAppByID(ctx context.Context, appID string) (*Application, error) {
	var a Application

	err := db.QueryRowContext(ctx, fmt.Sprintf(`SELECT
		APP_ID
	, APP_NAME
	, APP_URL
	, APP_DESC
	, APP_LOGO
	, INSERT_TS
	, UPDATE_TS
	FROM TSC_APPLICATIONS WHERE APP_ID = '%s'`, appID)).Scan(&a.APP_ID, &a.APP_NAME, &a.APP_URL, &a.APP_DESC, &a.APP_LOGO, &a.INSERT_TS, &a.UPDATE_TS)
	if err != nil {
		log.Println("error getting app by id", appID, err)
		return nil, err
	}

	return &a, nil
}

func getAppsPaginated(ctx context.Context, lpr PaginationRequest) (*PaginationResponse, error) {
	var apps []Application
	var numFilteredRecords int
	var pr PaginationResponse
	var sqlNumOfRecords = "SELECT COUNT(*) from TSC_APPLICATIONS "

	var err error
	pr.NumRecords, err = getTotalAmountOfRecords(ctx, sqlNumOfRecords)
	if err != nil {
		log.Println("error getting total number of records", err)
		return nil, err
	}
	if pr.NumRecords == -1 {
		log.Println("invalid number of records", pr)
	}

	whereClause := getSQLPaginationFilters(lpr.Filters)

	sqlStr := fmt.Sprintf(`SELECT
  		APP_ID
  	, APP_NAME
		, APP_URL
		, APP_DESC
		, APP_LOGO
		, INSERT_TS
		, UPDATE_TS
	FROM TSC_APPLICATIONS
	`)

	var sqlNumFilteredRecords string = sqlNumOfRecords
	if len(whereClause) > 0 {
		sqlStr += fmt.Sprintf("where %s", strings.Join(whereClause, " AND "))
		sqlNumFilteredRecords += fmt.Sprintf("where %s", strings.Join(whereClause, " AND "))
	}

	log.Println("sqlNumFilteredRecords", sqlNumFilteredRecords)

	err = db.QueryRowContext(ctx, sqlNumFilteredRecords).Scan(&numFilteredRecords)
	if err != nil {
		log.Println("error counting number of filtered records", err)
		return nil, err
	}
	pr.NumFilteredRecords = numFilteredRecords

	// sqlStr += " ORDER BY LOG_ID"

	if lpr.Parameters.Limit > 0 {
		pr.PageCount = int(math.Max(float64(numFilteredRecords/lpr.Parameters.Limit), float64(1)))
		sqlStr += fmt.Sprintf(" LIMIT %d", lpr.Parameters.Limit)
		if lpr.Parameters.CurrentPage > 1 {
			sqlStr += fmt.Sprintf(" OFFSET %d", lpr.Parameters.Limit*(lpr.Parameters.CurrentPage-1))
		}
	}

	log.Println("sqlStr", sqlStr)

	rows, err := db.QueryContext(ctx, sqlStr)
	if err != nil {
		log.Println("Error query TSC_Applications paginated", err)
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var a Application

		if err := rows.Scan(&a.APP_ID, &a.APP_NAME, &a.APP_URL, &a.APP_DESC, &a.APP_LOGO, &a.INSERT_TS, &a.UPDATE_TS); err != nil {
			log.Println("error scanning function block, index:", i, err)
			return nil, err
		}

		apps = append(apps, a)
		i++
	}

	pr.Data = apps

	return &pr, nil
}

func createFeedbackChannel(ctx context.Context, appID string, fc Feedback_Channel) (*Feedback_Channel, error) {

	var nfc Feedback_Channel
	var currentMaxChannelID *int
	var newMaxChannelID int
	err := db.QueryRowContext(ctx, "SELECT MAX(CHANNEL_ID) FROM TSC_FEEDBACK_CHANNELS").Scan(&currentMaxChannelID)
	if err != nil {
		log.Println("error getting feddback channel max id", err)
		return nil, err
	}
	if currentMaxChannelID == nil {
		newMaxChannelID = 1
	} else {
		newMaxChannelID = *currentMaxChannelID + 1
	}
	log.Println("next newMaxChannelID", newMaxChannelID)

	fc.APP_ID = appID
	fc.CHANNEL_ENDPOINT = fmt.Sprintf("/feedback/%s-%d", appID, newMaxChannelID)

	res, err := db.ExecContext(ctx, `INSERT INTO TSC_FEEDBACK_CHANNELS (
		  CHANNEL_ID
		, APP_ID
		, CHANNEL_NAME
		, CHANNEL_DESC
		, CHANNEL_ENDPOINT
) VALUES (?,?,?,?,?)
		`, newMaxChannelID, appID, fc.CHANNEL_NAME, fc.CHANNEL_DESC, fc.CHANNEL_ENDPOINT)
	if err != nil {
		log.Println("error creating feedback channel for appID", appID, err)
		return nil, err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		log.Println("last inserted id error", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		log.Println("last rowsaffected error", err)
	}

	log.Println("created new channel", lid, ra)

	err = db.QueryRowContext(ctx, `SELECT
		CHANNEL_ID
	, APP_ID
	, CHANNEL_NAME
	, CHANNEL_DESC
	, CHANNEL_ENDPOINT
	, INSERT_TS
	, UPDATE_TS
	FROM TSC_FEEDBACK_CHANNELS WHERE APP_ID = ? and CHANNEL_ID = ?`, appID, newMaxChannelID).Scan(&nfc.CHANNEL_ID, &nfc.APP_ID, &nfc.CHANNEL_NAME, &nfc.CHANNEL_DESC, &nfc.CHANNEL_ENDPOINT, &nfc.INSERT_TS, &nfc.UPDATE_TS)
	if err != nil {
		log.Println("error retrieving newly created feedback channel for appID", appID, err)
		return nil, err
	}

	log.Printf("retrieved newly inserted feedback channel %#v\n", nfc)

	return &nfc, nil
}

func getFeedbackChannelPaginated(ctx context.Context, lpr PaginationRequest, appID string) (*PaginationResponse, error) {
	var channels []Feedback_Channel
	var numFilteredRecords int
	var pr PaginationResponse
	var sqlNumOfRecords = "SELECT COUNT(*) from TSC_FEEDBACK_CHANNELS " //keep the space at the end for a following where or order by statement

	var err error
	pr.NumRecords, err = getTotalAmountOfRecords(ctx, sqlNumOfRecords)
	if err != nil {
		log.Println("error getting total number of records", err)
		return nil, err
	}
	if pr.NumRecords == -1 {
		log.Println("invalid number of records", pr)
	}

	whereClause := getSQLPaginationFilters(lpr.Filters)

	sqlStr := fmt.Sprintf(`SELECT
		  CHANNEL_ID
		, APP_ID
		, CHANNEL_NAME
  	, CHANNEL_DESC
		, CHANNEL_ENDPOINT
		, INSERT_TS
		, UPDATE_TS
	FROM TSC_FEEDBACK_CHANNELS WHERE APP_ID = '%s' `, appID) //keep the space for the following statements

	var sqlNumFilteredRecords string = sqlNumOfRecords
	if len(whereClause) > 0 {
		sqlStr += strings.Join(whereClause, " AND ")
		sqlNumFilteredRecords += strings.Join(whereClause, " AND ")
	}

	log.Println("sqlNumFilteredRecords", sqlNumFilteredRecords)

	err = db.QueryRowContext(ctx, sqlNumFilteredRecords).Scan(&numFilteredRecords)
	if err != nil {
		log.Println("error counting number of filtered records", err)
		return nil, err
	}
	pr.NumFilteredRecords = numFilteredRecords

	// sqlStr += " ORDER BY LOG_ID"

	if lpr.Parameters.Limit > 0 {
		pr.PageCount = int(math.Max(float64(numFilteredRecords/lpr.Parameters.Limit), float64(1)))
		sqlStr += fmt.Sprintf(" LIMIT %d", lpr.Parameters.Limit)
		if lpr.Parameters.CurrentPage > 1 {
			sqlStr += fmt.Sprintf(" OFFSET %d", lpr.Parameters.Limit*(lpr.Parameters.CurrentPage-1))
		}
	}

	log.Println("sqlStr", sqlStr)

	rows, err := db.QueryContext(ctx, sqlStr)
	if err != nil {
		log.Println("Error query TSC_Applications paginated", err)
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var fc Feedback_Channel

		if err := rows.Scan(&fc.CHANNEL_ID, &fc.APP_ID, &fc.CHANNEL_NAME, &fc.CHANNEL_DESC, &fc.CHANNEL_ENDPOINT, &fc.INSERT_TS, &fc.UPDATE_TS); err != nil {
			log.Println("error scanning function block, index:", i, err)
			return nil, err
		}

		channels = append(channels, fc)
		i++
	}

	pr.Data = channels

	return &pr, nil
}

func getFeedbackPaginated(ctx context.Context, lpr PaginationRequest, appID string, channelID int) (*PaginationResponse, error) {
	var feedbacks []Feedback
	var numFilteredRecords int
	var pr PaginationResponse
	var sqlNumOfRecords = "SELECT COUNT(*) from TSC_FEEDBACK " //keep the space at the end for a following where or order by statement

	var err error
	pr.NumRecords, err = getTotalAmountOfRecords(ctx, sqlNumOfRecords)
	if err != nil {
		log.Println("error getting total number of records", err)
		return nil, err
	}
	if pr.NumRecords == -1 {
		log.Println("invalid number of records", pr)
	}

	whereClause := getSQLPaginationFilters(lpr.Filters)

	sqlStr := fmt.Sprintf(`SELECT
  	  CHANNEL_ID
		, FEEDBACK_ID
  	, APP_ID
  	, FEEDBACK_TITLE
		, FEEDBACK_MESSAGE
		, FEEDBACK_POSITIVE_NEGATIVE
		, FEEDBACK_RAITING
		, REVIEWED
		, INSERT_TS
		, UPDATE_TS
	FROM TSC_FEEDBACK WHERE APP_ID = '%s' AND CHANNEL_ID = %d `, appID, channelID) //keep the space for the following statements

	var sqlNumFilteredRecords string = sqlNumOfRecords
	if len(whereClause) > 0 {
		sqlStr += strings.Join(whereClause, " AND ")
		sqlNumFilteredRecords += strings.Join(whereClause, " AND ")
	}

	log.Println("sqlNumFilteredRecords", sqlNumFilteredRecords)

	err = db.QueryRowContext(ctx, sqlNumFilteredRecords).Scan(&numFilteredRecords)
	if err != nil {
		log.Println("error counting number of filtered records", err)
		return nil, err
	}
	pr.NumFilteredRecords = numFilteredRecords

	// sqlStr += " ORDER BY LOG_ID"

	if lpr.Parameters.Limit > 0 {
		pr.PageCount = int(math.Max(float64(numFilteredRecords/lpr.Parameters.Limit), float64(1)))
		sqlStr += fmt.Sprintf(" LIMIT %d", lpr.Parameters.Limit)
		if lpr.Parameters.CurrentPage > 1 {
			sqlStr += fmt.Sprintf(" OFFSET %d", lpr.Parameters.Limit*(lpr.Parameters.CurrentPage-1))
		}
	}

	log.Println("sqlStr", sqlStr)

	rows, err := db.QueryContext(ctx, sqlStr)
	if err != nil {
		log.Println("Error query TSC_Applications paginated", err)
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var f Feedback

		if err := rows.Scan(&f.CHANNEL_ID, &f.FEEDBACK_ID, &f.APP_ID, &f.FEEDBACK_TITLE, &f.FEEDBACK_MESSAGE, &f.FEEDBACK_POSITIVE_NEGATIVE, &f.FEEDBACK_RAITING, &f.REVIEWED, &f.INSERT_TS, &f.UPDATE_TS); err != nil {
			log.Println("error scanning function block, index:", i, err)
			return nil, err
		}

		feedbacks = append(feedbacks, f)
		i++
	}

	pr.Data = feedbacks

	return &pr, nil
}

func saveLogMessages(ctx context.Context, logs []ClientLogs) error {

	var err error

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("error in begin tx", err)
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `bulk insert into TSC_CLIENT_LOGS (
		  SESSION_ID
		, APP_ID
		, LOG_LEVEL
		, URL
		, MSG
		, STACKTRACE
		, TIMESTAMP
		, USERAGENT
		, CLIENT_IP
		, REMOTE_IP
		) values (?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("error in prepare context bulk insert into TSC_CLIENT_LOGS", err)
		return err
	}
	defer stmt.Close()

	for i, l := range logs {

		// clobMsg := new(driver.Lob)
		// if l.MSG != "" {
		// 	clobMsg.SetReader(strings.NewReader(l.MSG))
		// } else {
		// 	clobMsg.SetReader(strings.NewReader(""))
		// }

		// var ioStringFunctionJSON *strings.Reader
		// var lobFunctionJSON *driver.Lob
		// if l.MSG != nil {
		// 	b, err := json.Marshal(l.MSG)
		// 	if err != nil {
		// 		log.Println(fmt.Sprint("error marshalling dynamic json structure: ", err))
		// 		return err
		// 	}
		// 	fmt.Println(string(b))
		// 	ioStringFunctionJSON = strings.NewReader(string(b))
		// 	// ioStringFunctionJSON = strings.NewReader(l.MSG)
		// 	lobFunctionJSON = new(driver.Lob)
		// 	lobFunctionJSON.SetReader(ioStringFunctionJSON)
		// }

		var inputData []interface{}
		inputData = append(inputData, l.SESSION_ID)
		inputData = append(inputData, l.APP_ID)
		inputData = append(inputData, l.LOG_LEVEL)
		inputData = append(inputData, l.URL)
		// inputData = append(inputData, lobFunctionJSON)
		inputData = append(inputData, l.MSG)
		inputData = append(inputData, l.STACKTRACE)
		inputData = append(inputData, l.TIMESTAMP)
		inputData = append(inputData, l.USERAGENT)
		inputData = append(inputData, l.CLIENT_IP)
		inputData = append(inputData, l.REMOTE_IP)

		if _, err := stmt.ExecContext(ctx, inputData...); err != nil {
			log.Println("error in bulk insert -> exec context for TSC_CLIENT_LOGS", i, err)
			return err
		}

	}
	// end bulk statements
	if _, err := stmt.ExecContext(ctx); err != nil {
		log.Println("error in exec context after loop for TSC_CLIENT_LOGS", err)
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		errMsg := fmt.Sprintln("error in commit", err)
		log.Println(errMsg)
		return errors.New(errMsg)
	}

	return err
}
