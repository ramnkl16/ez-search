package models

// import (
// 	"fmt"

// 	productsyncmysql "github.com/ramnkl16/ez-search/datasources/productsyncmysql"
// 	"github.com/ramnkl16/ez-search/logger"
// 	"github.com/ramnkl16/ez-search/rest_errors"
// 	"github.com/ramnkl16/ez-search/utils/mysql_utils"
// )

// //Auto code generated with help of mysql table schema
// // table : AppLog

// // Insert the AppLog to the database.
// func (al *AppLog) Create() rest_errors.RestErr {
// 	var err error
// 	// sql insert query, primary key must be provided
// 	const sqlstr = `INSERT INTO AppLog ( ` +
// 		` id, ref1, ref2, level, message, createdAt, updatedAt ` +
// 		`) VALUES (` +
// 		`?, ?, ?, ?, ?, ?, ?` +
// 		`)`

// 	stmt, err := productsyncmysql.Client.Prepare(sqlstr)
// 	if err != nil {
// 		logger.Error("error when trying to prepare save AppLog statement", err)
// 		return rest_errors.NewInternalServerError("error when tying to save AppLog ", err)
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(al.ID, al.Ref1, al.Ref2, al.Level, al.Message, al.CreatedAt, al.UpdatedAt)

// 	if err != nil {
// 		logger.Error("error when trying to save user", err)
// 		return rest_errors.NewInternalServerError("error when tying to save AppLog", err)
// 	}
// 	return nil
// }

// // Update updates the AppLog in the database.
// func (al *AppLog) Update() rest_errors.RestErr {

// 	// sql query
// 	const sqlstr = `UPDATE AppLog SET ` +
// 		` ref1 = ?, ref2 = ?, level = ?, message = ?, updatedAt = ? ` +
// 		` WHERE id  = ? `

// 	stmt, err := productsyncmysql.Client.Prepare(sqlstr)
// 	if err != nil {
// 		logger.Error("error when trying to prepare update AppLog statement", err)
// 		return rest_errors.NewInternalServerError("error when tying to update AppLog ", err)
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(al.Ref1, al.Ref2, al.Level, al.Message, al.UpdatedAt, al.ID)
// 	if err != nil {
// 		logger.Error("error when trying to update AppLog", err)
// 		return rest_errors.NewInternalServerError("error when tying to update AppLog", err)
// 	}
// 	return nil
// }

// // Delete a recored ( AppLog) from the database.
// func (al *AppLog) Delete() rest_errors.RestErr {

// 	const sqlstr = `UPDATE AppLog SET ` +
// 		`  updatedAt=?  ` +
// 		` WHERE id  = ? `

// 	stmt, err := productsyncmysql.Client.Prepare(sqlstr)
// 	if err != nil {
// 		logger.Error("error when trying to prepare delete AppLog statement", err)
// 		return rest_errors.NewInternalServerError("error when tying to delete AppLog ", err)
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(al.UpdatedAt, al.ID)
// 	if err != nil {
// 		logger.Error("error when trying to delete AppLog", err)
// 		return rest_errors.NewInternalServerError("error when tying to delete AppLog", err)
// 	}
// 	return nil
// }

// // search  from ( AppLog).
// func (al *AppLog) Searh(start string, limit string) ([]AppLog, rest_errors.RestErr) {

// 	sqlstr := `SELECT ` +
// 		` id, ref1, ref2, level, message ` +
// 		` FROM AppLog  ` +
// 		"LIMIT " + start + "," + limit

// 	stmt, err := productsyncmysql.Client.Prepare(sqlstr)
// 	if err != nil {
// 		logger.Error("error when trying to prepare search AppLog", err)
// 		return nil, rest_errors.NewInternalServerError("error when tying search AppLog", err)
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query()
// 	if err != nil {
// 		logger.Error("error when trying to search AppLog", err)
// 		return nil, rest_errors.NewInternalServerError("error when tying to search AppLog", err)
// 	}
// 	defer rows.Close()
// 	// load results
// 	res := make([]AppLog, 0)
// 	for rows.Next() {
// 		var mal AppLog
// 		// scan
// 		err = rows.Scan(&mal.ID, &mal.Ref1, &mal.Ref2, &mal.Level, &mal.Message)
// 		if err != nil {
// 			logger.Error("error when scan user row into AppLog struct", err)
// 			return nil, rest_errors.NewInternalServerError("error when trying to search AppLog", err)
// 		}
// 		res = append(res, mal)
// 	}
// 	if len(res) == 0 {
// 		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no AppLog matching limit %s - %s ", start, limit))
// 	}
// 	return res, nil
// }

// // Get a record  from ( AppLog) .
// func (al *AppLog) Get() rest_errors.RestErr {
// 	const sqlstr = `SELECT ` +
// 		` id, ref1, ref2, level, message ` +
// 		` FROM AppLog WHERE id = ? `
// 	// run query
// 	stmt, err := productsyncmysql.Client.Prepare(sqlstr)
// 	if err != nil {
// 		logger.Error("error when trying to prepare get AppLog statement", err)
// 		return rest_errors.NewInternalServerError("error when tying to get AppLog", err)
// 	}
// 	defer stmt.Close()
// 	row := stmt.QueryRow(al.ID)
// 	// scan
// 	err = row.Scan(&al.ID, &al.Ref1, &al.Ref2, &al.Level, &al.Message)
// 	if err != nil {
// 		if sqlErr := mysql_utils.ParseError(err); sqlErr != nil {
// 			return sqlErr
// 		}
// 		logger.Error("error when scan user row into AppLog struct", err)
// 		return rest_errors.NewInternalServerError("error when tying to get AppLog", err)
// 	}
// 	return nil
// }
