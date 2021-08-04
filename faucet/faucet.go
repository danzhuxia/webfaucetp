package faucet

import (
	"context"
	"net/http"
	"webfaucetp/utils"
	"webfaucetp/utils/datastorage"
	"webfaucetp/utils/errmsg"
	"webfaucetp/utils/validator"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SprinkleAddr(c *gin.Context) {
	var code int
	var ipAddress string
	if ipAddr, ok := c.RemoteIP(); !ok {
		code = errmsg.CAN_NOT_GET_ACCESS_IP
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
	} else {
		ipAddress = ipAddr.String()
	}

	var spRo datastorage.SprinkleRecord
	err := c.ShouldBindJSON(&spRo)
	if err != nil {
		code = errmsg.CAN_NOT_BIND_PARAMS
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		})
		return
	}
	spRo.IPAddress = ipAddress
	//check params
	var msg string
	msg, code = validator.Validate(&spRo)
	if code != errmsg.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": msg,
		})
		return
	}

	// check database
	check := datastorage.CheckRecordInMysql(&spRo)
	switch check {
	case errmsg.SUCCESS:
		tx := datastorage.NewTx(&spRo, sprinkle)
		code = tx.DoInsertTransaction()
	case errmsg.UPDATE_ACCESS:
		tx := datastorage.NewTx(&spRo, sprinkle)
		code = tx.DoUpdateTransaction()
	default:
		code = check
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": errmsg.GetErrMsg(code),
	})
}

func sprinkle(record *datastorage.SprinkleRecord) (code int) {
	// Create an IPC based RPC connection to a remote node and instantiate a contract binding

	toAddress := common.HexToAddress(record.EtherumAddress) //目标地址
	val, _ := utils.Sana.BalanceOf(nil, toAddress)
	logrus.Infof("before transfer :%s", val)

	tx, err := utils.Sana.Transfer(utils.Auth, toAddress, utils.Amount)
	if err != nil {
		logrus.Errorf("Failed to request token transfer: %s", err.Error())
		return errmsg.REQUEST_TOKEN_TRANSFER_FAILED
	}
	ctx := context.Background()

	receipt, err := bind.WaitMined(ctx, utils.EthClient, tx)
	if err != nil {
		logrus.Errorf("Failed to Mine Transaction: %s", err.Error())
		return errmsg.TRANSACTION_MINING_FAILED
	}
	val, _ = utils.Sana.BalanceOf(nil, toAddress)
	logrus.Infof("After Transfer: %s", val)
	logrus.Infof("Transaction Hash: %s", tx.Hash())
	logrus.Infof("Receipt :%s", receipt.TxHash)

	return errmsg.SUCCESS
}
