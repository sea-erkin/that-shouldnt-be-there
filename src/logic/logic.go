package logic

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/sea-erkin/that-shouldnt-be-there/src/repo"
)

var (
	print = fmt.Println
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func containsHost(s []repo.HostDb, e repo.HostDb) bool {
	for _, a := range s {
		if a.Host == e.Host {
			return true
		}
	}
	return false
}

func GetHostCompareTimestamps(hostMapTimestamp map[int64][]repo.HostDb) (int64, int64) {
	var keys = make([]int64, 0)
	for k, _ := range hostMapTimestamp {
		keys = append(keys, k)
	}
	var newerTimestamp int64
	var olderTimestamp int64
	if keys[0] > keys[1] {
		newerTimestamp = keys[0]
		olderTimestamp = keys[1]
	} else {
		newerTimestamp = keys[1]
		olderTimestamp = keys[0]
	}
	return olderTimestamp, newerTimestamp
}

func GetNewApproveHosts(hostMapTimestamp map[int64][]repo.HostDb, olderTimestamp, newerTimestamp int64) ([]repo.HostDb, []repo.HostDb) {
	var newHosts = make([]repo.HostDb, 0)
	var hostsToApprove = make([]repo.HostDb, 0)
	for _, host := range hostMapTimestamp[newerTimestamp] {
		if !containsHost(hostMapTimestamp[olderTimestamp], host) {
			newHosts = append(newHosts, host)
		} else {
			hostsToApprove = append(hostsToApprove, host)
		}
	}
	return newHosts, hostsToApprove
}

func GetMissingInactivateHosts(hostMapTimestamp map[int64][]repo.HostDb, olderTimestamp, newerTimestamp int64) ([]repo.HostDb, []repo.HostDb) {
	var missingHosts = make([]repo.HostDb, 0)
	var hostsToInactivate = make([]repo.HostDb, 0)
	for _, host := range hostMapTimestamp[olderTimestamp] {
		if !containsHost(hostMapTimestamp[newerTimestamp], host) {
			missingHosts = append(missingHosts, host)
		} else {
			hostsToInactivate = append(hostsToInactivate, host)
		}
	}
	return missingHosts, hostsToInactivate
}

// Could create readLines
func ParseSubdomain(fileLocation string) ([]string, string) {
	fileBytes, err := ioutil.ReadFile(fileLocation)
	checkErr(err)

	fileString := string(fileBytes[:])

	fileLines := strings.Split(fileString, "\n")
	// remove last line b\c new line
	fileLines = fileLines[0 : len(fileLines)-1]

	domainNameLine := fileLines[0]
	domainName := strings.TrimPrefix(domainNameLine, "www.")

	return fileLines, domainName
}
