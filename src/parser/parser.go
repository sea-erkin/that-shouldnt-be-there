package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

var (
	print = fmt.Println
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ReadLines(fileLocation string) []string {
	fileBytes, err := ioutil.ReadFile(fileLocation)
	checkErr(err)

	fileString := string(fileBytes[:])

	fileLines := strings.Split(fileString, "\n")
	// remove last line b\c new line
	fileLines = fileLines[0 : len(fileLines)-1]

	return fileLines
}
