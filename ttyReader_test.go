package main

import (
	"testing"
	"bufio"
	"os"
	"log"
	"github.com/stretchr/testify/assert"
)

func TestDataParsing(t *testing.T) {
	f, err := os.Open("./resources/test-data")
	reader := bufio.NewReader(f)
	reply, err := reader.ReadBytes('\x21')
	if err != nil {
		panic(err)
	}

	readData := string(reply)
	matches := matchData(readData)
	log.Print(len(matches))

	subMatchMap := map[string]string{
		"omis": "2.8.1",
		"data": "795.989",
		"unit": "kWh",
	}

	assert.Contains(t, matches, subMatchMap)
}
