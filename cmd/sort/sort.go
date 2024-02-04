package main

import (
	"fmt"
	"log/slog"
	"os"
	"sort-tool/internal/config"
	"sort-tool/internal/sort"
)

func main() {
	logger := SetupLogger()
	cfg := config.New(logger)
	err := cfg.ParseFlags()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to parse flags: %v", err))
	}

	var fileSort *sort.FileSort
	fileSort, err = sort.New(cfg)
	if err = fileSort.Sort(); err != nil {
		logger.Error(fmt.Sprintf("can't sort data %v", err))
		os.Exit(1)
	}

	if err = fileSort.Write(); err != nil {
		logger.Error(fmt.Sprintf("can't write data: %v", err))
		os.Exit(1)
	}
}

func SetupLogger() *slog.Logger {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
