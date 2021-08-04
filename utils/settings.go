package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"webfaucetp/utils/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/ini.v1"
)

var (
	Host string
	Port string

	DbType     string
	DbAddr     string
	DbName     string
	DbUser     string
	DbPassword string

	RawUrl     string
	Contract   string
	SanaAmount string
	PrivateKey string
)

var (
	EthClient *ethclient.Client
	Amount    *big.Int
	Sana      *token.Token
	Auth      *bind.TransactOpts
)

const key = `{"address":"5cf1d1b737071fc6f190a2b7f10189107d445aab","crypto":{"cipher":"aes-128-ctr","ciphertext":"4b86af2288d127d3b64f4002b0ddf14a68b51855d017f803ea0ceaddccc41666","cipherparams":{"iv":"22f01583503d6ee0bb0bd875a5e2d40c"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"95489436fffcebe9db0b0361e3710c092bd2b3678e2e6e5fdbad742d78ba6e99"},"mac":"bcdbe94264b26d2dd0239f9a0665689fe1a049388d64d31a8266bb2975b777aa"},"id":"f0c88354-7604-4a1b-94a7-0f92f259036f","version":3}`

func init() {
	file, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Read Configure From File Failed...", err)
	}

	LoadServer(file)
	LoadDB(file)
	LoadOthers(file)
}

func LoadServer(file *ini.File) {
	Host = file.Section("server").Key("host").String()
	Port = file.Section("server").Key("port").String()
}

func LoadDB(file *ini.File) {
	DbType = file.Section("db").Key("db_type").String()
	DbAddr = file.Section("db").Key("db_addr").String()
	DbName = file.Section("db").Key("db_name").String()
	DbUser = file.Section("db").Key("db_user").String()
	DbPassword = file.Section("db").Key("db_password").String()
}

func LoadOthers(file *ini.File) {
	RawUrl = file.Section("others").Key("raw_url").String()
	Contract = file.Section("others").Key("contract").String()
	SanaAmount = file.Section("others").Key("sana_amount").String()
	PrivateKey = file.Section("others").Key("private_key").String()
}

func InitFaucet() (err error) {
	var ok bool
	if Amount, ok = big.NewInt(0).SetString(SanaAmount, 10); !ok {
		return fmt.Errorf("get amount error")
	}

	EthClient, err = ethclient.Dial(RawUrl)
	if err != nil {
		// log.Printf("Failed to connect to the Ethereum client: %v", err)
		return fmt.Errorf("failed to connect to the Ethereum client: %s", err.Error())
	}

	Sana, err = token.NewToken(common.HexToAddress(Contract), EthClient)
	if err != nil {
		return fmt.Errorf("failed to instantiate a Token contract: %s", err.Error())
	}

	var chainId *big.Int
	chainId, err = EthClient.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch chain_ID: %s", err.Error())
	}

	Auth, err = bind.NewTransactorWithChainID(strings.NewReader(PrivateKey), key, chainId)
	// Auth, err = bind.NewTransactor(strings.NewReader(PrivateKey), key)
	if err != nil {
		return fmt.Errorf("failed to create authorized transactor: %s", err.Error())
	}
	return nil
}
