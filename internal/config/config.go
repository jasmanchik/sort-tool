package config

import (
	"flag"
	"log/slog"
	"os"
)

type KeyColumn struct {
	Num       int
	IsNumeric bool
}

type Config struct {
	logger      *slog.Logger
	InFile      *os.File
	OutFile     *os.File
	IsReverse   bool
	OnlyUnique  bool
	KeyColumn   *KeyColumn
	InFileName  string
	OutFileName string
}

func New(logger *slog.Logger) *Config {
	return &Config{
		logger:     logger,
		IsReverse:  false,
		OnlyUnique: false,
		KeyColumn:  &KeyColumn{Num: 1, IsNumeric: false},
	}
}

func (c *Config) ParseFlags() error {
	column := flag.Int("k", 1, "column to sort on")
	numeric := flag.Bool("n", false, "numeric sort")
	reverse := flag.Bool("r", false, "reverse sort")
	unique := flag.Bool("u", false, "unique sort")
	outSrc := flag.String("o", "", "file output name")
	inSrc := flag.String("i", "", "file output name")
	flag.Parse()

	c.InFileName = *inSrc
	c.OutFileName = *outSrc
	c.IsReverse = *reverse
	c.KeyColumn.Num = *column
	c.KeyColumn.IsNumeric = *numeric
	c.OnlyUnique = *unique

	return nil
}
