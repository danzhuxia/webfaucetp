package router

import (
	"fmt"
	"webfaucetp/faucet"
	"webfaucetp/middleware"
	"webfaucetp/utils"

	"github.com/gin-gonic/gin"
)

func InitRouter() error {
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	auth := r.Group(`v1`)
	{
		auth.POST("/faucet", faucet.SprinkleAddr)
	}

	err := r.Run(fmt.Sprintf("%s:%s", utils.Host, utils.Port))
	if err != nil {
		return err
	}
	return nil
}
