package main

import (
	"bufio"
	"fmt"
	"os"
)

func readPasswordFile(f string) (string, error) {
	var line string
	var fd *os.File
	var err error

	fd, err = os.Open(f)
	if err != nil {
		return line, err
	}

	scanner := bufio.NewScanner(fd)
	scanner.Scan()
	line = scanner.Text()
	fd.Close()
	if line == "" {
		return line, fmt.Errorf("Empty password read from file %s", f)
	}

	return line, nil
}
