package config

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func InitLogger() (*zap.SugaredLogger, error) {
	if Log == nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatalln(err)
		}

		Log = logger.Sugar()
	}
	return Log, nil
}
