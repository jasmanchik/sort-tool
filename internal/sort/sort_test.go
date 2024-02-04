package sort

import (
	"bufio"
	"errors"
	"os"
	"sort-tool/internal/config"
	"testing"
)

func getDataFromFile(s string) []string {
	f, err := os.Open(s)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	var data []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	return data
}

func TestNew(t *testing.T) {
	cfg := &config.Config{
		KeyColumn: &config.KeyColumn{Num: 1, IsNumeric: false},
	}
	fs, err := New(cfg)
	if err != nil {
		t.Errorf("New() error = %v, wantErr %v", err, nil)
	}
	if fs == nil {
		t.Errorf("New() = %v, want non-nil FileSort", fs)
	}
}

func TestReverseSort(t *testing.T) {
	cfg := &config.Config{
		InFileName: "../../test/testdata/input.txt",
		IsReverse:  true,
		KeyColumn:  &config.KeyColumn{Num: 1, IsNumeric: false},
	}

	fileSort, err := New(cfg)
	err = fileSort.Sort()
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}

	sortedData := fileSort.data
	sortedSlice := make([]string, 0)
	for _, str := range sortedData {
		sortedSlice = append(sortedSlice, str.Line)
	}

	expectData := getDataFromFile("../../test/testdata/output_reverse.txt")
	for i, str := range expectData {
		if str != sortedSlice[i] {
			t.Fatalf("got = %v, want %v", sortedSlice[i], str)
		}
	}
}

func TestColumnSort(t *testing.T) {
	cfg := &config.Config{
		InFileName: "../../test/testdata/input.txt",
		KeyColumn:  &config.KeyColumn{Num: 4, IsNumeric: false},
	}

	fileSort, err := New(cfg)
	err = fileSort.Sort()
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}

	sortedData := fileSort.data
	sortedSlice := make([]string, 0)
	for _, str := range sortedData {
		sortedSlice = append(sortedSlice, str.Line)
	}

	expectData := getDataFromFile("../../test/testdata/output_column_sort_4.txt")
	for i, str := range expectData {
		if str != sortedSlice[i] {
			t.Fatalf("got = %v, want %v", sortedSlice[i], str)
		}
	}
}

func TestColumnNumericSort(t *testing.T) {
	cfg := &config.Config{
		InFileName: "../../test/testdata/input.txt",
		KeyColumn:  &config.KeyColumn{Num: 3, IsNumeric: true},
	}

	fileSort, err := New(cfg)
	err = fileSort.Sort()
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}

	sortedData := fileSort.data
	sortedSlice := make([]string, 0)
	for _, str := range sortedData {
		sortedSlice = append(sortedSlice, str.Line)
	}

	expectData := getDataFromFile("../../test/testdata/output_column_numeric_sort.txt")
	for i, str := range expectData {
		if str != sortedSlice[i] {
			t.Fatalf("got = %v, want %v", sortedSlice[i], str)
		}
	}
}

func TestUniqueSort(t *testing.T) {
	cfg := &config.Config{
		InFileName: "../../test/testdata/input.txt",
		OnlyUnique: true,
		KeyColumn:  &config.KeyColumn{Num: 1, IsNumeric: false},
	}

	fileSort, err := New(cfg)
	err = fileSort.Sort()
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}

	sortedData := fileSort.data
	sortedSlice := make([]string, 0)
	for _, str := range sortedData {
		sortedSlice = append(sortedSlice, str.Line)
	}

	expectData := getDataFromFile("../../test/testdata/output_unique_sort.txt")
	for i, str := range expectData {
		if str != sortedSlice[i] {
			t.Fatalf("got = %v, want %v", sortedSlice[i], str)
		}
	}
}

func TestOutFileSort(t *testing.T) {
	cfg := &config.Config{
		OutFileName: "../../test/testdata/output_sort_file.txt",
		InFileName:  "../../test/testdata/input.txt",
		KeyColumn:   &config.KeyColumn{Num: 1, IsNumeric: false},
	}

	fileSort, err := New(cfg)
	err = fileSort.Sort()
	if err != nil {
		t.Fatalf("Sort() error = %v", err)
	}
	err = fileSort.Write()
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	_, err = os.Stat("../../test/testdata/output_sort_file.txt")
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("output_sort_file.txt is not exist")
	}

	err = os.Remove("../../test/testdata/output_sort_file.txt")
	if err != nil {
		t.Fatalf("failed to remove output_sort_file.txt: %v", err)
	}
}
