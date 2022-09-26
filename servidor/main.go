package main

import (
	"github.com/soyReymundus/social/domain"
	"github.com/soyReymundus/social/web"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	errEnv := godotenv.Load(".env")
	Check(errEnv)

	dmn := domain.Domain{}
	server := web.Server{}

	dmn.Go()
	server.Open(&dmn)
}
