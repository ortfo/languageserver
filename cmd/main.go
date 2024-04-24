package main

import (
	"fmt"
	"os"

	"github.com/ortfo/languageserver"
	"go.uber.org/zap"
)

func main() {
	os.MkdirAll("./logs", os.ModePerm)
	logconf := zap.NewDevelopmentConfig()
	logconf.OutputPaths = []string{"./logs/server.log"}
	logger, _ := logconf.Build()

	if len(os.Args) < 2 || os.Args[1] == "" {
		fmt.Println("set the path to the ortfodb.yaml path as the first argument")
		os.Exit(1)
	}

	languageserver.StartServer(logger, os.Args[1], "./logs/")
}
