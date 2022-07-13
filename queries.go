package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/SAP/go-hdb/driver"
)

func getLogsPaginated(ctx context.Context, lpr logsPaginationRequest) (*PaginationResponse, error) {

	var cls []ClientLogs

	var whereClause []string
	var numTotalRecords int
	var numFilteredRecords int
	var sqlNumOfRecords = "SELECT COUNT(*) from TSC_CLIENT_LOGS "

	var pr PaginationResponse

	// get total number of Records
	err := db.QueryRowContext(ctx, sqlNumOfRecords).Scan(&numTotalRecords)
	if err != nil {
		log.Println("error collecting number of Total REcords for Function Blocks", err)
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
  		SESSION_ID
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

	sqlStr += " ORDER BY SESSION_ID, TIMESTAMP"

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
		log.Println("Error query function blocks paginated", err)
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var l ClientLogs

		//For CLOB Column
		b := new(bytes.Buffer)
		var lobRuleJSON driver.NullLob
		lobRuleJSON.Lob = new(driver.Lob)
		lobRuleJSON.Lob.SetWriter(b)

		// var ioStringFunctionJSON *strings.Reader
		// lob := &driver.Lob{}
		// lob.SetReader(ioStringFunctionJSON)

		// lob := new(driver.Lob)
		// b := new(bytes.Buffer)
		// lob.SetWriter(b)

		if err := rows.Scan(&l.SESSION_ID, &l.LOG_LEVEL, &l.URL, &lobRuleJSON, &l.STACKTRACE, &l.TIMESTAMP, &l.USERAGENT, &l.CLIENT_IP, &l.REMOTE_IP); err != nil {
			log.Println("error scanning function block, index:", i, err)
			return nil, err
		}

		// b, err := json.Marshal(lob)
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(string(b))
		// l.MSG = lob

		if lobRuleJSON.Valid { // only valid if data from DB is not NULL
			// stringValue := string(b.Bytes())
			// ruleVersion.RULE_JSON = &stringValue

			l.MSG = string(b.Bytes())
		}

		// l.MSG = string(b)
		// wr := &bytes.Buffer{}
		// lob.SetWriter(wr)

		// l.MSG = *lob

		cls = append(cls, l)
		i++
	}

	pr.Data = cls

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
		, LOG_LEVEL
		, URL
		, MSG
		, STACKTRACE
		, TIMESTAMP
		, USERAGENT
		, CLIENT_IP
		, REMOTE_IP
		) values (?,?,?,?,?,?,?,?,?)`)
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

		var ioStringFunctionJSON *strings.Reader
		var lobFunctionJSON *driver.Lob
		if l.MSG != nil {
			b, err := json.Marshal(l.MSG)
			if err != nil {
				log.Println(fmt.Sprint("error marshalling dynamic json structure: ", err))
				return err
			}
			fmt.Println(string(b))
			ioStringFunctionJSON = strings.NewReader(string(b))
			// ioStringFunctionJSON = strings.NewReader(l.MSG)
			lobFunctionJSON = new(driver.Lob)
			lobFunctionJSON.SetReader(ioStringFunctionJSON)
		}

		var inputData []interface{}
		inputData = append(inputData, l.SESSION_ID)
		inputData = append(inputData, l.LOG_LEVEL)
		inputData = append(inputData, l.URL)
		inputData = append(inputData, lobFunctionJSON)
		// inputData = append(inputData, l.MSG)
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
