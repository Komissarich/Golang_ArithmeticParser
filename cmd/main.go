package main

import (
	"calc/application"
)

func main() {
	app := application.New()

	app.RunServer()
}
