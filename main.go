package main

import (
	"mw-project/model"
	"mw-project/routes"
)

func main() {
	router := routes.Setup()

	model.ConnectDatabase()

	router.Run("localhost:9000")
}