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
	source  *os.File
	params  *config.Config
	rawData string
	data    []Item
}

type Item struct {
	Line string
	Key  string
}

func New(config *config.Config) (*FileSort, error) {
	file, err := os.Open(config.FileName)
	if err != nil {
		return nil, err
	}

	return &FileSort{
		source:  file,
		params:  config,
		rawData: "",
		data:    make([]Item, 0),
	}, nil
}

func (f *FileSort) Close() error {
	err := f.source.Close()
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSort) read() error {
	scanner := bufio.NewScanner(f.source)
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
		return err
	}

	return nil
}

func (f *FileSort) Sort() error {
	if err := f.read(); err != nil {
		return fmt.Errorf("can't read file %v", err)
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
		return f.data[i].Key < f.data[j].Key
	}

	sort.Slice(f.data, sortFunc)

	if f.params.IsReverse {
		f.doReverse()
	}

	return nil
}

func (f *FileSort) doReverse() bool {
	if len(f.data) <= 1 {
		return false
	}
	for i, j := 0, len(f.data)-1; i < j; i, j = i+1, j-1 {
		f.data[i], f.data[j] = f.data[j], f.data[i]
	}
	return true
}

var NoDataError = errors.New("there is no data to write")

func (f *FileSort) Write(out *os.File) error {
	if len(f.data) <= 0 {
		return NoDataError
	}

	writer := bufio.NewWriter(out)
	hasWritten := make(map[string]struct{})
	for _, item := range f.data {
		if f.params.OnlyUnique {
			if _, ok := hasWritten[item.Line]; ok {
				continue
			}
			hasWritten[item.Line] = struct{}{}
		}

		if _, err := writer.WriteString(item.Line + "\n"); err != nil {
			return fmt.Errorf("can't write string %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("can't flush data %w", err)
	}

	return nil
}
