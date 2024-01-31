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
	logger      *slog.Logger
	FileName    *os.File
	OutFileName *os.File
	IsReverse   bool
	OnlyUnique  bool
	KeyColumn   *KeyColumn
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

	var file *os.File
	if fileExists(*inSrc) {
		var err error
		file, err = os.Open(*inSrc)
		if err != nil {
			flag.PrintDefaults()
			return fmt.Errorf("failed open output file %v", err)
		}
	} else {
		var err error
		file, err = os.Create(*inSrc)
		if err != nil {
			flag.PrintDefaults()
			return fmt.Errorf("failed create output file %v", err)
		}
	}
	var fileOut *os.File
	if *outSrc != "" {
		var err error
		fileOut, err = os.Create(*outSrc)
		if err != nil {
			flag.PrintDefaults()
			return fmt.Errorf("failed create file %v", err)
		}
	} else {
		fileOut = os.Stdout
	}

	c.FileName = file
	c.OutFileName = fileOut
	c.IsReverse = *reverse
	c.KeyColumn.Num = *column
	c.KeyColumn.IsNumeric = *numeric
	c.OnlyUnique = *unique

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}
