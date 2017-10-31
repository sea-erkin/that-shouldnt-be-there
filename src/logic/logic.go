package logic

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sea-erkin/that-shouldnt-be-there/src/alert"
	"github.com/sea-erkin/that-shouldnt-be-there/src/common"
	"github.com/sea-erkin/that-shouldnt-be-there/src/dto"
	"github.com/sea-erkin/that-shouldnt-be-there/src/parser"
	"github.com/sea-erkin/that-shouldnt-be-there/src/repo"
)

var (
	Debug = false
)

func print(params ...interface{}) {
	if Debug == true {
		fmt.Println(params)
	}
}

func ApproveApprovalCode(approvalCodeFlag string) {
	print("Trying to approve with the following uuid:", approvalCodeFlag)
	if repo.DbSetIpPortAlertApprovedByUuid(approvalCodeFlag) {
		print("Approved with the following uuid:", approvalCodeFlag)
	}
	if repo.DbSetHostAlertApprovedByUuid(approvalCodeFlag) {
		print("Approved with the following uuid:", approvalCodeFlag)
	}
}

func CreateIpHostsFromResolvedSubdomains(subdomainDoneResolveDirectory, subdomainArchiveResolveDirectory string) {
	files, err := ioutil.ReadDir(subdomainDoneResolveDirectory)
	common.CheckErr(err)

	for _, file := range files {
		if string(file.Name()[0]) == "." {
			continue
		}
		ipHosts := parser.ReadLines(subdomainDoneResolveDirectory + file.Name())

		fileSplit := strings.Split(file.Name(), "_")
		domainName := strings.ToLower(fileSplit[0])
		timeScanned := fileSplit[2]

		print("Inserting IPhosts: ", ipHosts, domainName)

		// create function to insert the ip hosts as IP_host into sqlite
		// once that is done create the alerting logic.
		// once that is done, automate the nmap scans for these IPs and automate screenshots.
		repo.DbInsertIpHostsForDomain(ipHosts, domainName, timeScanned)

		err = os.Rename(subdomainDoneResolveDirectory+file.Name(), subdomainArchiveResolveDirectory+file.Name())
	}
}

func SendCompletedScreenshots(screenshotDoneDirectory, screenshotArchiveDirectory, fromEmail, emailPassword, emailHost, emailPort string, recipients []string) {
	files, err := ioutil.ReadDir(screenshotDoneDirectory)
	common.CheckErr(err)

	for _, file := range files {
		if string(file.Name()[0]) == "." {
			continue
		}

		print("Sending email for the following file: ", file.Name())

		alert.SendMailAttachment("", fromEmail, emailPassword,
			screenshotDoneDirectory+file.Name(), emailHost,
			emailPort, recipients)

		err = os.Rename(screenshotDoneDirectory+file.Name(), screenshotArchiveDirectory+file.Name())

		common.CheckErr(err)
	}
}

// ParseSubdomains loop through every file in the todo directory and parse / insert
// once completed, move each all to done and store the last completed file.
func ParseSubdomains(subdomainTodoDirectory, subdomainTodoResolveDirectory, subdomainDoneDirectory, subdomainLastCompleted string) {
	files, err := ioutil.ReadDir(subdomainTodoDirectory)
	common.CheckErr(err)

	for _, file := range files {
		if string(file.Name()[0]) == "." {
			continue
		}

		hosts := parser.ReadLines(subdomainTodoDirectory + file.Name())

		fileSplit := strings.Split(file.Name(), "_")
		domainName := strings.ToLower(fileSplit[0])
		dataSource := strings.ToLower(fileSplit[1])
		timeScanned := fileSplit[2]

		repo.DbInsertHostsForDomain(hosts, domainName, dataSource, timeScanned)
		err := ioutil.WriteFile(subdomainLastCompleted, []byte(file.Name()), 0644)
		common.CheckErr(err)

		err = common.CopyFile(subdomainTodoDirectory+file.Name(), subdomainTodoResolveDirectory+file.Name())
		err = os.Rename(subdomainTodoDirectory+file.Name(), subdomainDoneDirectory+file.Name())

		common.CheckErr(err)
	}
}

func CreateHostPortsFromNmapFile(parseNmapDataFlag string) {
	print("Beginning XML Parse")
	mv, err := parser.NmapGetMapFromXmlFile(parseNmapDataFlag)
	if err != nil {
		print(err)
	}

	paths := mv.PathsForKey("host")
	hosts, err := mv.ValuesForPath(paths[0])

	// Get Host Ports
	var hostPorts []dto.HostPort
	for _, host := range hosts {
		hostEndTimeString := host.(map[string]interface{})["-endtime"].(string)
		hostEndTime, err := strconv.ParseInt(hostEndTimeString, 10, 64)
		if err != nil {
			print(err)
		}
		hostAddress := host.(map[string]interface{})["address"].(map[string]interface{})["-addr"].(string)
		ports := host.(map[string]interface{})["ports"]
		portsSlice := ports.(map[string]interface{})["port"]
		switch v := portsSlice.(type) {
		case []interface{}:
			for _, port := range portsSlice.([]interface{}) {
				hostPort := parser.NmapGetPortFromXml(port, hostAddress, hostEndTime)
				hostPorts = append(hostPorts, hostPort)
			}
		case map[string]interface{}:
			hostPort := parser.NmapGetPortFromXml(portsSlice, hostAddress, hostEndTime)
			hostPorts = append(hostPorts, hostPort)
		default:
			log.Fatal("unexpected type %T", v)
		}
	}

	// Insert host ports
	for _, hostPort := range hostPorts {
		print("Inserting the following record:", hostPort)
		repo.DbCreateIpPort(hostPort.IP, hostPort.Protocol, hostPort.Port, hostPort.Service, hostPort.State, hostPort.TimeScanned)
	}
}

func AlertSubdomains() ([]repo.HostDb, []repo.HostDb) {

	domains := repo.DbViewHostUniqueDomains()
	sources := repo.DbViewHostUniqueSources()

	print("Domains to alert subdomains", domains)
	print("Domains to alert sources", sources)

	var missingHostsAll = make([]repo.HostDb, 0)
	var newHostsAll = make([]repo.HostDb, 0)
	for _, domain := range domains {

		for _, source := range sources {

			hosts := repo.DbViewHost(domain, source)
			if len(hosts) == 0 {
				print("No Hosts to process alerts. Skipping.")
				continue
			}

			print("Hosts Count", len(hosts))
			print("Hosts", hosts)

			// for every host, group by timestamp.
			var hostMapTimestamp = make(map[int64][]repo.HostDb)
			for _, host := range hosts {
				hostMapTimestamp[host.TimeScanned] = append(hostMapTimestamp[host.TimeScanned], host)
			}

			print(len(hostMapTimestamp))

			if len(hostMapTimestamp) == 1 {
				print("No other hosts to compare against. Must be first run. Skipping.")
				continue
			}

			olderTimestamp, newerTimestamp := GetHostCompareTimestamps(hostMapTimestamp)

			print("Old hosts:", hostMapTimestamp[olderTimestamp])

			newHosts, hostsToApprove := GetNewApproveHosts(hostMapTimestamp, olderTimestamp, newerTimestamp)
			missingHosts, hostsToInactivate := GetMissingInactivateHosts(hostMapTimestamp, olderTimestamp, newerTimestamp)

			print("Missing hosts:", missingHosts)
			print("New Hosts", newHosts)

			// Set status to an alert for missing hosts and new hosts
			repo.DbCreateHostAlert(missingHosts)
			repo.DbCreateHostAlert(newHosts)

			print("Hosts to approve", hostsToApprove)
			print("Hosts to inactivate", hostsToInactivate)

			repo.DbSetHostSubdomainStatus(hostsToApprove, "Approved")
			repo.DbSetHostSubdomainStatus(hostsToInactivate, "Inactive")

			// Insert alert uuid entries so can be approved
			repo.DbCreateAlertHostUuid(missingHosts)
			repo.DbCreateAlertHostUuid(newHosts)

			for _, m := range missingHosts {
				missingHostsAll = append(missingHostsAll, m)
			}

			for _, n := range newHosts {
				newHostsAll = append(newHostsAll, n)
			}
		}
	}
	return missingHostsAll, newHostsAll
}

func ApproveLatestIpPorts(v []repo.IPPortDb, alertPortFlag, screenshotTodoDirectory string) {
	var timeSlice = make(map[int64][]repo.IPPortDb)
	for _, item := range v {
		timeSlice[item.DateScanned] = append(timeSlice[item.DateScanned], item)
	}
	keys := make([]int, 0)
	for k, _ := range timeSlice {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// latest timeslice entries.
	latestTimesliceKey := keys[len(keys)-1]
	latestEntries := timeSlice[int64(latestTimesliceKey)]

	// grab the IPPortIds to save these jawns
	ipPortIdsToMarkApproved := make([]string, 0)
	for _, item := range latestEntries {
		ipPortIdsToMarkApproved = append(ipPortIdsToMarkApproved, strconv.FormatInt(item.IPPortId, 10))
	}
	alert.PrepNmapScreenshot(latestEntries, alertPortFlag, screenshotTodoDirectory)
	print("Approving the following IPPortIds", ipPortIdsToMarkApproved)
	repo.DbSetIpPortIdsStatus(ipPortIdsToMarkApproved, "Approved")
}

// CreateAlertsForMissingAddedPorts
// Set alert flag for added and missing ports
// Creates alert approval uuids so can be marked as approved for the future
func CreateAlertsForPorts(ports []repo.IPPortDb) {

	var ipPortIdsToCreateAlertsFor = make([]string, 0)
	for _, item := range ports {
		ipPortIdsToCreateAlertsFor = append(ipPortIdsToCreateAlertsFor, strconv.FormatInt(item.IPPortId, 10))
	}
	repo.DbCreateIpPortAlert(ipPortIdsToCreateAlertsFor)

	// create an alert approval entry uuid so that we have a means to approve. Will be command line but could ultimately be gui
	for _, item := range ipPortIdsToCreateAlertsFor {
		ipPortId, err := strconv.ParseInt(item, 10, 64)
		common.CheckErr(err)
		uuid, _ := common.NewUUID()
		repo.DbCreateAlertIpPortUuid(ipPortId, uuid)
	}
}

// SetPortIdsApproved sets the passed in ports to approved
func SetPortIdsApproved(ports []repo.IPPortDb) {

	newlyApprovedPortIds := make([]string, 0)
	for _, item := range ports {
		newlyApprovedPortIds = append(newlyApprovedPortIds, strconv.FormatInt(item.IPPortId, 10))
	}
	print("Setting the following newer timestamps to approved: ", newlyApprovedPortIds)
	repo.DbSetIpPortIdsStatus(newlyApprovedPortIds, "Approved")

}

// SetPortIdsInactive
func SetPortIdsInactive(ports []repo.IPPortDb) {
	oldApprovedPortIdsToInactivate := make([]string, 0)
	for _, item := range ports {
		oldApprovedPortIdsToInactivate = append(oldApprovedPortIdsToInactivate, strconv.FormatInt(item.IPPortId, 10))
	}
	print("Setting the following older timestamps to not approved: ", oldApprovedPortIdsToInactivate)
	repo.DbSetIpPortIdsStatus(oldApprovedPortIdsToInactivate, "Inactive")
}

func ContainsHost(s []repo.HostDb, e repo.HostDb) bool {
	for _, a := range s {
		if a.Host == e.Host {
			return true
		}
	}
	return false
}

func ContainsIpPort(s []repo.IPPortDb, e repo.IPPortDb) bool {
	for _, a := range s {
		if a.Port == e.Port && a.State == e.State {
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
		if !ContainsHost(hostMapTimestamp[olderTimestamp], host) {
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
		if !ContainsHost(hostMapTimestamp[newerTimestamp], host) {
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
	common.CheckErr(err)

	fileString := string(fileBytes[:])

	fileLines := strings.Split(fileString, "\n")
	// remove last line b\c new line
	fileLines = fileLines[0 : len(fileLines)-1]

	domainNameLine := fileLines[0]
	domainName := strings.TrimPrefix(domainNameLine, "www.")

	return fileLines, domainName
}
