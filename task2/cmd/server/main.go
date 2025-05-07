package main

import (
	"vk-Service/task2/internal"
	"vk-Service/task2/internal/config"
)

func main() {
	logger := config.NewLogger()
	err := config.SetUp()
	if err != nil {
		logger.Error("failed to setup config", "err", err)
	}
	internal.RunServer(logger)
}
