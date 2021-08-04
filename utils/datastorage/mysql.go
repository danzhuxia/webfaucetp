package datastorage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"webfaucetp/utils"
	"webfaucetp/utils/errmsg"

	"database/sql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SqlDb *sql.DB
var db *gorm.DB
var err error

//github.com/ethereum/go-ethereum
func InitMysql() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		utils.DbUser, utils.DbPassword, utils.DbHost, utils.DbPort, utils.DbName)
	SqlDb, err = sql.Open(utils.DbType, dsn)
	db, err = gorm.Open(mysql.New(mysql.Config{
		// DriverName:                utils.DbType,
		Conn:                      SqlDb,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: true,
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&SprinkleRecord{})
	SqlDb.SetMaxIdleConns(10)
	SqlDb.SetMaxOpenConns(100)
	SqlDb.SetConnMaxLifetime(time.Hour)

	return nil
}

// Fetch Address
func CheckRecordInMysql(SpRo *SprinkleRecord) int {
	var addrs []string
	etherrumAddr := SpRo.EtherumAddress
	err = db.Select("id").Where("ip_address=?", SpRo.IPAddress).First(&SpRo).Error
	if err != nil {
		return errmsg.DB_OPTIONS_FAILED
	}
	//先check一下Id,看有没有记录
	if SpRo.ID > 0 {
		//先检查是否达到次数
		if SpRo.Sprinkles >= 5 {
			return errmsg.REACH_MAX_SPRINKLE_TIMES
		}

		addrs, err = SpRo.SprinkledToSlice()
		if err != nil {
			return errmsg.DECODE_SPRINKLED_ADDR_FAILED
		}
		for _, address := range addrs {
			if address == etherrumAddr {
				return errmsg.ETHERUM_ADDRESS_HAS_SPRINKLED
			}
		}
		return errmsg.UPDATE_ACCESS
	}

	return errmsg.SUCCESS
}

func insertRecordMysql(dbtx *gorm.DB, sr *SprinkleRecord) int {
	spRo := &SprinkleRecord{
		IPAddress:      sr.IPAddress,
		EtherumAddress: sr.EtherumAddress,
	}
	err = dbtx.Create(spRo).Error
	if err != nil {
		dbtx.Rollback()
		return errmsg.INSERT_RECORD_FAILED
	}
	return errmsg.SUCCESS
}

func UpdateRecordMysql(dbtx *gorm.DB, sr *SprinkleRecord) int {
	var spRo SprinkleRecord
	var sprinkles []string
	err = dbtx.Where("ip_address=?", sr.IPAddress).First(&spRo).Error
	if err != nil {
		dbtx.Callback()
		return errmsg.DB_OPTIONS_FAILED
	}
	sprinkles, err = spRo.SprinkledToSlice()
	if err != nil {
		dbtx.Callback()
		return errmsg.DECODE_SPRINKLED_ADDR_FAILED
	}
	sprinkles = append(sprinkles, sr.EtherumAddress)
	byts, _ := json.Marshal(&sprinkles)
	spRo.EtherumAddress = sr.EtherumAddress
	spRo.Sprinkles++
	spRo.SprinkledAddrs = bytes.NewBuffer(byts).String()
	err = dbtx.Save(&spRo).Error
	if err != nil {
		dbtx.Callback()
		return errmsg.DB_OPTIONS_FAILED
	}
	return errmsg.SUCCESS
}
