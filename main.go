package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

type rowSearcher struct {
	chRowsInBulk <-chan []string
	current []string
}

func (rs *rowSearcher) Search(row string) (found bool, inRange bool, exhausted bool) {
	if len(rs.current) == 0 {
		var ok bool
		select {
		case rs.current, ok = <-rs.chRowsInBulk:
			if !ok {
				exhausted = true
			}
		}
	}
	i := sort.SearchStrings(rs.current, row)
	if i < len(rs.current) {
		inRange = true
		if rs.current[i] == row {
			found = true
		}
	} else {
		rs.current = nil
	}
	return
}

const usage = "Usage: rose <op:[and,sub]> <file1> <file2>"

func main() {
	flag.Parse()
	if flag.NArg() != 3 {
		fmt.Println(usage)
		os.Exit(1)
	}
	var files [2]*os.File
	for i := 1; i <= 2; i++ {
		name := flag.Arg(i)
		file, err := os.Open(name)
		if err != nil {
			fmt.Printf("Failed to open file %d: %s\n", i, err)
			os.Exit(1)
		}
		files[i-1] = file
	}
	switch flag.Arg(0) {
	case "and":
		intersect(files[0], files[1])
	case "sub":
		complement(files[0], files[1])
	default:
		 fmt.Println(usage)
		 os.Exit(1)
	}
}

func intersect(f1, f2 *os.File) {
	scanner1 := bufio.NewScanner(f1)
	searcher := &rowSearcher{
		chRowsInBulk: readLinesInBulk(f2, 64),
	}
	var lastLine string
	for scanner1.Scan() {
		line := scanner1.Text()
		if len(line) == 0 {
			continue
		}
		if line == lastLine {
			 continue
		}
		lastLine = line
		for {
			found, inRange, exhausted := searcher.Search(line)
			if found {
				fmt.Println(line)
				break
			}
			if inRange {
				break
			}
			if exhausted {
				return
			}
		}
	}
}

func complement(f1, f2 *os.File) {
	scanner1 := bufio.NewScanner(f1)
	searcher := &rowSearcher{
		chRowsInBulk: readLinesInBulk(f2, 64),
	}
	var lastLine string
	var f2Exhausted bool
	for scanner1.Scan() {
		line := scanner1.Text()
		if len(line) == 0 {
			continue
		}
		if line == lastLine {
			continue
		}
		lastLine = line
		if f2Exhausted {
			fmt.Println(line)
			continue
		}
		for {
			found, inRange, exhausted := searcher.Search(line)
			if found {
				break
			}
			if inRange || exhausted {
				fmt.Println(line)
				f2Exhausted = exhausted
				break
			}
		}
	}
}

func readLinesInBulk(reader io.Reader, bulkSize int) <-chan []string {
	scanner2 := bufio.NewScanner(reader)
	chReadLines2 := make(chan []string, 2)
	go func() {
		readLines2 := make([]string, 0, bulkSize)
		for scanner2.Scan() {
			line := scanner2.Text()
			if len(line) == 0 {
				continue
			}
			readLines2 = append(readLines2, line)
			if len(readLines2) >= bulkSize {
				chReadLines2 <- readLines2
				readLines2 = make([]string, 0, bulkSize)
			}
		}
		if len(readLines2) > 0 {
			chReadLines2 <- readLines2
		}
		close(chReadLines2)
	}()
	return chReadLines2
}