package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func LoadCredentialsFromCsv(authFile string) (map[string]string, error) {
	if len(authFile) <= 0 {
		return nil, errors.New("No auth credentials file specified")
	}
	csvFile, _ := os.Open(authFile)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	log.Infof("Parsing auth csv file: %s", authFile)
	result := map[string]string{}
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if len(line) == 0 {
			// Empty line
			continue
		}
		result[line[0]] = line[1]
	}
	log.Infof("Found %d credentials in file specified", len(result))
	return result, nil
}
