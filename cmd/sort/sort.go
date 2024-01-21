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
	cfg.ParseFlags()

	fileSort, err := sort.New(cfg)
	defer func() {
		if err := fileSort.Close(); err != nil {
			logger.Error(fmt.Sprintf("can't close fileSort %v", err))
			os.Exit(1)
		}
	}()

	if err := fileSort.Sort(); err != nil {
		logger.Error(fmt.Sprintf("can't sort data %v", err))
		os.Exit(1)
	}

	// Создание файла для записи результатов
	outputFile, err := os.Create("sorted.txt")
	if err != nil {
		logger.Error(fmt.Sprintf("can't create a file: %v", err))
		os.Exit(1)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			logger.Error(fmt.Sprintf("can't close a new file: %v", err))
			os.Exit(1)
		}
	}()

	if err := fileSort.Write(outputFile); err != nil {
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
