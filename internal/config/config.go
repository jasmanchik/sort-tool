package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type KeyColumn struct {
	Num       int
	IsNumeric bool
}

type Config struct {
	logger     *slog.Logger
	FileName   string
	IsReverse  bool
	OnlyUnique bool
	KeyColumn  *KeyColumn
}

func New(logger *slog.Logger) *Config {
	return &Config{
		logger:     logger,
		FileName:   "",
		IsReverse:  false,
		OnlyUnique: false,
		KeyColumn:  &KeyColumn{Num: 1, IsNumeric: false},
	}
}

func (c *Config) ParseFlags() {
	column := flag.Int("k", 1, "column to sort on")
	numeric := flag.Bool("n", false, "numeric sort")
	reverse := flag.Bool("r", false, "reverse sort")
	unique := flag.Bool("u", false, "unique sort")
	flag.Parse()

	// Проверка оставшихся аргументов (имя файла)
	if flag.NArg() != 1 {
		c.logger.Error(fmt.Sprintf("Usage: %s [options] <file>\n", os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	c.FileName = filename
	c.IsReverse = *reverse
	c.KeyColumn.Num = *column
	c.KeyColumn.IsNumeric = *numeric
	c.OnlyUnique = *unique
}
