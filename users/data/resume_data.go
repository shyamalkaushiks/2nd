package data

import (
	model "users/models"

	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
)

// InsertUserResume :
func InsertUserResume(objUserResumesToUpload []interface{}) error {
	db := model.DBConn

	err := gormbulk.BulkInsert(db, objUserResumesToUpload, 100)
	if err != nil {
		return err
	}

	return nil
}
