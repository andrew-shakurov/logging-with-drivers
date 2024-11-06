package main

import (
	"os"

	"example.com/log"
	_ "example.com/logdriverjson"
)

func main() {
	err := log.GlobalLogFactory.ConfigureFromFile("config.json")
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
	logger := log.GlobalLogFactory.NewLog()
	logger.Log("Hello World!", nil)
}
