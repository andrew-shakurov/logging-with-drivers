package main

import (
	"os"

	"example.com/log"
	_ "example.com/logdriverjson"
)

func main() {
	err := log.ConfigureFromFile("config.json")
	if err != nil {
		print(err.Error())
		os.Exit(1)
	}
	logger := log.NewLog()
	defer log.Close()

	logger.Log("Hello World!", nil)
}
