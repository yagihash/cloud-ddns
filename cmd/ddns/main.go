package main

import (
	"fmt"
	"os"

	"github.com/yagihash/cloud-ddns/logger"

	"github.com/yagihash/cloud-ddns/config"
)

const (
	exitOK = iota
	exitError
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	env, err := config.Load()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[error]", err)
		return exitError
	}

	log, sync, err := logger.New(env.LogPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "[error]", err)
		return exitError
	}
	defer sync()

	log.Info("app started watching global ip address")

	log.Info("app shutting down")
	return exitOK
}
