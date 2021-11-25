package env

import (
	"fmt"
	"os"
	"strconv"
)

// VarName is the name of environment variable.
type VarName string

// AppEnv represents application runtime environment.
type AppEnv string

func get(name VarName) (value string, err error) {
	val := os.Getenv(string(name))
	if val == "" {
		return val, fmt.Errorf("failed to load environment variable. key: %s", name)
	}

	return val, nil
}

func String(name VarName) string {
	val, err := get(name)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	return val
}

func Int(name VarName) int {
	val, err := get(name)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("env value is not parsabl. key: %s, val: %s: %w", name, val, err))
	}

	return intVal
}

const (
	EnvTest AppEnv = "test"
	EnvDev  AppEnv = "dev"
	EnvStg  AppEnv = "stg"
	EnvPrd  AppEnv = "prd"
)

func Env(name VarName, fallback AppEnv) AppEnv {
	val, err := get(name)
	if err != nil {
		return fallback
	}

	return AppEnv(val)
}

func (e AppEnv) IsTest() bool {
	return e == EnvTest
}

func (e AppEnv) IsDev() bool {
	return e == EnvDev
}

func (e AppEnv) IsStg() bool {
	return e == EnvStg
}

func (e AppEnv) IsPrd() bool {
	return e == EnvPrd
}

func (e AppEnv) IsCloud() bool {
	return e.IsDev() || e.IsStg() || e.IsPrd()
}
