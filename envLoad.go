package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type EnvStruct struct {
	PasswordDb string
	PortDb     string
	Namedb     string
	Jwt_secret string
	DbHost     string
}

var Env EnvStruct

func EnvLoad(loger LogFile) (EnvStruct, error) {

	err := godotenv.Load()
	if err != nil {
		loger.LogEnv.Println("Ошибка загрузки .env", err)
	}

	Env.PasswordDb = os.Getenv("PASSWORDB")
	Env.PortDb = os.Getenv("PORTDB")
	Env.Namedb = os.Getenv("NAMEDB")
	Env.Jwt_secret = os.Getenv("JWT_SECRET")
	Env.DbHost = os.Getenv("DBHOST")
	if Env.Jwt_secret == "" {
		err := errors.New("Пустой Secret")
		return EnvStruct{}, err
	}
	return Env, nil
}
