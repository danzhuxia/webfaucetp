package main

import (
	"log"
	"webfaucetp/router"
	"webfaucetp/utils"
	"webfaucetp/utils/datastorage"
)

func main() {
	var err error
	err = datastorage.InitMysql()
	if err != nil {
		log.Fatalf("init mysql failed: %s", err.Error())
	}
	defer datastorage.SqlDb.Close()

	err = utils.InitFaucet()
	if err != nil {
		log.Fatalf("init faucet failed: %s", err.Error())
	}

	err = router.InitRouter()
	if err != nil {
		log.Fatalf("init router failed: %s", err.Error())
	}
}
