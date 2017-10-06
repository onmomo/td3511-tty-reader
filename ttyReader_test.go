package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
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

	kwhMap := map[string]string{
		"omis": "2.8.1",
		"data": "795.989",
		"unit": "kWh",
	}

	assert.Contains(t, matches, kwhMap)

	kwMap := map[string]string{
		"omis": "1.7.0",
		"data": "0.348",
		"unit": "kW",
	}

	assert.Contains(t, matches, kwMap)

	kvarMap := map[string]string{
		"omis": "4.7.0",
		"data": "0.296",
		"unit": "kvar",
	}

	assert.Contains(t, matches, kvarMap)

	vMap := map[string]string{
		"omis": "32.7",
		"data": "232.4",
		"unit": "V",
	}

	assert.Contains(t, matches, vMap)

	aMap := map[string]string{
		"omis": "31.7",
		"data": "0.58",
		"unit": "A",
	}

	assert.Contains(t, matches, aMap)
}