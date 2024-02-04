package sort

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"sort-tool/internal/config"
	"strconv"
	"strings"
)

type FileSort struct {
	params  *config.Config
	rawData string
	data    []Item
}

type Item struct {
	Line string
	Key  string
}

func New(config *config.Config) (*FileSort, error) {
	return &FileSort{
		params: config,
		data:   make([]Item, 0),
	}, nil
}

func (f *FileSort) Close() error {

	err := f.params.OutFile.Close()
	if err != nil {
		return fmt.Errorf("can't close file %v", err)
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}

func (f *FileSort) read() error {
	if fileExists(f.params.InFileName) {
		var err error
		f.params.InFile, err = os.Open(f.params.InFileName)
		if err != nil {
			return fmt.Errorf("failed open output file %v", err)
		}
	} else {
		var err error
		fmt.Println(f.params.InFileName)
		f.params.InFile, err = os.Create(f.params.InFileName)
		if err != nil {
			return fmt.Errorf("failed create output file %v", err)
		}
	}

	scanner := bufio.NewScanner(f.params.InFile)
	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.Fields(line)
		key := ""
		if f.params.KeyColumn.Num-1 < len(columns) {
			key = columns[f.params.KeyColumn.Num-1]
		}
		f.data = append(f.data, Item{Line: line, Key: key})
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("can't scan file %v", err)
	}

	return nil
}

func (f *FileSort) Sort() error {
	if err := f.read(); err != nil {
		return fmt.Errorf("can't read file %v", err)
	}

	var sortedData []Item
	if f.params.OnlyUnique {
		uniqueKeys := make(map[string]struct{})
		for _, item := range f.data {
			key := item.Key
			if _, exists := uniqueKeys[key]; !exists {
				uniqueKeys[key] = struct{}{}
				sortedData = append(sortedData, item)
			}
		}
	} else {
		sortedData = f.data
	}

	sortFunc := func(i, j int) bool {
		if f.params.KeyColumn.IsNumeric {
			numI, errI := strconv.Atoi(f.data[i].Key)
			numJ, errJ := strconv.Atoi(f.data[j].Key)
			if errI == nil && errJ == nil {
				return numI < numJ
			}
			return false
		}
		if f.params.IsReverse {
			return sortedData[i].Key > sortedData[j].Key
		}
		return sortedData[i].Key < sortedData[j].Key
	}

	sort.Slice(sortedData, sortFunc)
	f.data = sortedData

	return nil
}

var NoDataError = errors.New("there is no data to write")

func (f *FileSort) Write() error {
	if len(f.data) <= 0 {
		return NoDataError
	}

	if f.params.OutFileName != "" {
		var err error
		f.params.OutFile, err = os.Create(f.params.OutFileName)
		if err != nil {
			return fmt.Errorf("failed create file %v", err)
		}
	} else {
		f.params.OutFile = os.Stdout
	}

	writer := bufio.NewWriter(f.params.OutFile)
	for _, item := range f.data {
		if _, err := writer.WriteString(item.Line + "\n"); err != nil {
			return fmt.Errorf("can't write string %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("can't flush data %w", err)
	}

	return nil
}
