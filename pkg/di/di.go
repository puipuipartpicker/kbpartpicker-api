package di

import (
	"log"

	"github.com/puipuipartpicker/kbpartpicker/api/pkg/env"
	"github.com/puipuipartpicker/kbpartpicker/api/pkg/logging"
)

const (
	envAppEnv                              env.VarName = "APP_ENV"
	envLogLevel                            env.VarName = "LOG_LEVEL"
)

func GetAppEnv() env.AppEnv {
	return env.Env(envAppEnv, env.EnvTest)
}

func GetLogger() logging.Logger {
	if !GetAppEnv().IsCloud() {
		l, err := logging.NewDevelopmentLogger()
		if err != nil {
			log.Fatalf("failed to init logger: %+v", err)
		}

		return l
	}

	l, err := logging.NewLogger(env.String(envLogLevel))
	if err != nil {
		log.Fatalf("failed to init logger: %+v", err)
	}

	return l
}

func GetMainLogger() logging.Logger {
	return GetLogger().Named("main")
}

func GetContextLogger() logging.ContextLogger {
	if !GetAppEnv().IsCloud() {
		l, err := logging.NewDevelopmentContextLogger(logging.ContextParser)
		if err != nil {
			log.Fatalf("failed to init logger: %+v", err)
		}

		return l
	}

	l, err := logging.NewContextLogger(
		env.String(envLogLevel),
		logging.ContextParser,
	)
	if err != nil {
		log.Fatalf("failed to init logger: %+v", err)
	}

	return l
}
