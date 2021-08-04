package datastorage

import (
	"webfaucetp/utils/errmsg"

	"gorm.io/gorm"
)

type Transaction struct {
	DB       *gorm.DB
	Record   *SprinkleRecord
	CallFunc func(*SprinkleRecord) int
}

func NewTx(record *SprinkleRecord, callbackfunc func(*SprinkleRecord) int) *Transaction {
	return &Transaction{
		DB:       db,
		Record:   record,
		CallFunc: callbackfunc,
	}
}

func (tx *Transaction) DoInsertTransaction() (code int) {
	dbtx := tx.DB.Begin()
	code = insertRecordMysql(dbtx, tx.Record)
	if code != errmsg.SUCCESS {
		return
	}
	code = tx.CallFunc(tx.Record)
	if code != errmsg.SUCCESS {
		dbtx.Rollback()
		return
	}
	err = dbtx.Commit().Error
	if err != nil {
		code = errmsg.ERROR
	}
	return
}

// Duplicate code, needs to be improved
func (tx *Transaction) DoUpdateTransaction() (code int) {
	dbtx := tx.DB.Begin()
	code = insertRecordMysql(dbtx, tx.Record)
	if code != errmsg.SUCCESS {
		return
	}
	code = tx.CallFunc(tx.Record)
	if code != errmsg.SUCCESS {
		dbtx.Rollback()
		return
	}
	err = dbtx.Commit().Error
	if err != nil {
		code = errmsg.ERROR
	}
	return
}
