package config

import (
	"github.com/joho/godotenv"
	"github.com/mbakhodurov/examples/week_4/di/ufo/internal/config/env"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	UFOGRPC UFOGRPCCONFIG
	Mongo   MongoConfig
}

func Load(path ...string) error {
	if err := godotenv.Load(path...); err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	ufoGRPCCfg, err := env.NewUFOGPRCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoCOnfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:  loggerCfg,
		UFOGRPC: ufoGRPCCfg,
		Mongo:   mongoCfg,
	}
	return nil
}

func AppConfing() *config {
	return appConfig
}
