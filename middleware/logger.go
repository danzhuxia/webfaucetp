package middleware

import (
	"fmt"
	"math"
	"os"
	"time"

	rotate "github.com/lestrrat-go/file-rotatelogs"

	"github.com/gin-gonic/gin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	filePath := "log/faucet"
	linkName := "latest_log.log"
	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)

	logWriter, _ := rotate.New(
		filePath+"%Y%m%d.log",
		rotate.WithMaxAge(7*24*time.Hour),
		rotate.WithRotationTime(24*time.Hour),
		rotate.WithLinkName(linkName),
	)
	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.DebugLevel: logWriter,
		logrus.TraceLevel: logWriter,
	}
	hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.AddHook(hook)

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds())/1000000.0)))
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "Unknown"
		}
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"Status":    statusCode,
			"SpendTime": spendTime,
			"IP":        clientIP,
			"Method":    method,
			"Path":      path,
			"Agent":     userAgent,
			"DataSize":  dataSize,
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
