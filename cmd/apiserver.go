package main

import (
	"aipos/cmd/app"
	"log"
)

func main() {
	cmd := app.NewApiServerCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
