package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SortItem struct {
	Line string
	Key  string
}

func main() {
	// Парсинг флагов командной строки
	column := flag.Int("k", 1, "column to sort on")
	numeric := flag.Bool("n", false, "numeric sort")
	reverse := flag.Bool("r", false, "reverse sort")
	unique := flag.Bool("u", false, "unique sort")
	flag.Parse()

	// Проверка оставшихся аргументов (имя файла)
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	// Чтение файла
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		os.Exit(1)
	}
	defer file.Close()

	items := make([]SortItem, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		columns := strings.Fields(line)
		key := ""
		if *column-1 < len(columns) {
			key = columns[*column-1]
		}
		items = append(items, SortItem{Line: line, Key: key})
	}

	// Функция сортировки
	sortFunc := func(i, j int) bool {
		if *numeric {
			numI, errI := strconv.Atoi(items[i].Key)
			numJ, errJ := strconv.Atoi(items[j].Key)
			if errI == nil && errJ == nil {
				return numI < numJ
			}
		}
		return items[i].Key < items[j].Key
	}

	// Применение сортировки
	sort.Slice(items, sortFunc)

	if *reverse {
		for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
			items[i], items[j] = items[j], items[i]
		}
	}

	// Создание файла для записи результатов
	outputFile, err := os.Create("sorted.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Запись в файл
	writer := bufio.NewWriter(outputFile)
	seen := make(map[string]bool)
	for _, item := range items {
		if *unique {
			if _, ok := seen[item.Line]; ok {
				continue
			}
			seen[item.Line] = true
		}
		fmt.Fprintln(writer, item.Line)
	}
	writer.Flush()

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}
}
