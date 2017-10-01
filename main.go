package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/sea-erkin/that-shouldnt-be-there/src/alert"
	"github.com/sea-erkin/that-shouldnt-be-there/src/common"
	"github.com/sea-erkin/that-shouldnt-be-there/src/logic"
	"github.com/sea-erkin/that-shouldnt-be-there/src/parser"
	"github.com/sea-erkin/that-shouldnt-be-there/src/repo"

	mxj "github.com/clbanning/mxj"
)

var (
	configFlag         = flag.String("c", "", "Config file location")
	debugFlag          = flag.Bool("d", false, "Debug flag")
	alertPortFlag      = flag.Bool("alertPort", false, "Alert port logic flag")
	alertSubdomainFlag = flag.Bool("alertSubdomain", false, "Alert subdomain flag")
	parseNmapDataFlag  = flag.Bool("parseNmap", false, "Parse Nmap data flag")
	parseSubdomainFlag = flag.Bool("parseSubdomain", false, "Parse subdomain output. Must specify where to look for files")
	parseHostIpFlag    = flag.Bool("parseHostIp", false, "Parses host ip output.")
	approvalCodeFlag   = flag.String("approvalCode", "", "Approve UUID for alert code")
)

var (
	config Config
)

type HostPort struct {
	IP          string
	Protocol    string
	Port        string
	State       string
	Service     string
	TimeScanned int64
}

type Config struct {
	DbLocation                       string   `json:"dbLocation"`
	NmapParseFile                    string   `json:"nmapParseFile"`
	DomainFile                       string   `json:"domainFile"`
	SubdomainTodoDirectory           string   `json:"subdomainTodoDirectory"`
	SubdomainDoneDirectory           string   `json:"subdomainDoneDirectory"`
	SubdomainLastCompleted           string   `json:"subdomainLastCompleted"`
	SubdomainTodoResolveDirectory    string   `json:"subdomainTodoResolveDirectory"`
	SubdomainDoneResolveDirectory    string   `json:"subdomainDoneResolveDirectory"`
	SubdomainArchiveResolveDirectory string   `json:"subdomainArchiveResolveDirectory"`
	EmailHost                        string   `json:"emailHost"`
	EmailPort                        string   `json:"emailPort"`
	FromEmail                        string   `json:"fromEmail"`
	EmailPassword                    string   `json:"emailPassword"`
	Recipients                       []string `json:"recipients"`
	ErrorRecipients                  []string `json:"errorRecipients"`
}

type Map map[string]interface{}

func print(params ...interface{}) {
	if *debugFlag == true {
		fmt.Println(params)
	}
}

func checkFlags() {
	flag.Parse()

	if *configFlag == "" {
		log.Fatal("Must supply a config file.")
	}

	if *alertPortFlag == false && *parseNmapDataFlag == false && *approvalCodeFlag == "" &&
		*parseSubdomainFlag == false && *alertSubdomainFlag == false && *parseHostIpFlag == false {
		log.Fatal("Must specify something to do. Either alert, parse data, or approve an alert.")
	}
}

func main() {

	checkFlags()
	config = getConfig(*configFlag)
	repo.InitializeDb(config.DbLocation)

	if *approvalCodeFlag != "" {
		print("Trying to approve with the following uuid:", *approvalCodeFlag)
		if repo.DbSetIpPortAlertApprovedByUuid(*approvalCodeFlag) {
			print("Approved with the following uuid:", *approvalCodeFlag)
		}
		if repo.DbSetHostAlertApprovedByUuid(*approvalCodeFlag) {
			print("Approved with the following uuid:", *approvalCodeFlag)
		}
	}

	// Alert nmap logic. Move to it's own shit.
	if *alertPortFlag {
		print("Starting analysis for changes in open external ports")
		// get ip ports
		ipPorts := repo.DbViewIpPort()
		print("Count of IPPorts: ", len(ipPorts))
		if len(ipPorts) == 0 {
			print("No records to analyze")
			os.Exit(0)
		}
		var ipPortIpSlice = make(map[string][]repo.IPPortDb)
		for _, ipPort := range ipPorts {
			ipPortIpSlice[ipPort.IP] = append(ipPortIpSlice[ipPort.IP], ipPort)
		}
		for ip, v := range ipPortIpSlice {
			print("Analyzing IP: ", ip)
			approvedPorts := repo.DbGetApprovedPorts(ip)
			if len(approvedPorts) > 1 {
				newPorts := ipPortIpSlice[ip]
				print("IP Has previous approved history: ", ip)
				print("Approved ports: ", approvedPorts)
				print("Ports to compare against: ", newPorts)

				// Does not account for timestamp
				if reflect.DeepEqual(approvedPorts, newPorts) {
					print("Approved ports are the same as new ports, nothing to do here.")
					continue
				}

				// check new ports are equal to approved ports
				samePorts := make([]repo.IPPortDb, 0)
				missingPorts := make([]repo.IPPortDb, 0)
				for _, item := range approvedPorts {
					if containsIpPort(newPorts, item) {
						samePorts = append(samePorts, item)
					} else {
						missingPorts = append(missingPorts, item)
					}
				}

				addedPorts := make([]repo.IPPortDb, 0)
				for _, item := range newPorts {
					if !containsIpPort(approvedPorts, item) {
						addedPorts = append(addedPorts, item)
					}
				}

				if len(samePorts) == len(approvedPorts) {
					// stop here and set the new ports to approved, old ports to not approved.
					print("Ports are the same but have a newer timestamp. Updating entries")

					newlyApprovedPortIds := make([]string, 0)
					oldApprovedPortIdsToInactivate := make([]string, 0)
					for _, item := range newPorts {
						newlyApprovedPortIds = append(newlyApprovedPortIds, strconv.FormatInt(item.IPPortId, 10))
					}
					for _, item := range approvedPorts {
						oldApprovedPortIdsToInactivate = append(oldApprovedPortIdsToInactivate, strconv.FormatInt(item.IPPortId, 10))
					}
					print("Setting the following newer timestamps to approved: ", newlyApprovedPortIds)
					print("Setting the following older timestamps to not approved: ", oldApprovedPortIdsToInactivate)
					print("Removing alerts for the following IP if they exist: ", ip)
					repo.DbSetIpPortIdsStatus(newlyApprovedPortIds, "Approved")
					repo.DbSetIpPortIdsStatus(oldApprovedPortIdsToInactivate, "Inactive")
					repo.DbSetIpPortAlertStatusByIpAddress(ip)
					continue
				}

				print("New ports found: ", addedPorts)
				print("Ports missing:", missingPorts)

				// Set alert flag for added and missing ports
				var ipPortIdsToCreateAlertsFor = make([]string, 0)
				for _, item := range addedPorts {
					ipPortIdsToCreateAlertsFor = append(ipPortIdsToCreateAlertsFor, strconv.FormatInt(item.IPPortId, 10))
				}
				for _, item := range missingPorts {
					ipPortIdsToCreateAlertsFor = append(ipPortIdsToCreateAlertsFor, strconv.FormatInt(item.IPPortId, 10))
				}
				repo.DbCreateIpPortAlert(ipPortIdsToCreateAlertsFor)

				// create an alert approval entry uuid so that we have a means to approve. Will be command line but could ultimately be gui
				for _, item := range ipPortIdsToCreateAlertsFor {
					ipPortId, err := strconv.ParseInt(item, 10, 64)
					checkErr(err)
					uuid, _ := common.NewUUID()
					repo.DbCreateAlertIpPortUuid(ipPortId, uuid)
				}

				// Prep and send E-mail code. Will eventually support multiple alert formats such as text & app push notifs
				body := createEmailBodyFromAlertablePorts(addedPorts, missingPorts)
				alert.SendMail(body, config.FromEmail, config.EmailPassword, config.EmailHost, config.EmailPort, config.Recipients)

			} else {
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

				print("Approving the following IPPortIds", ipPortIdsToMarkApproved)
				repo.DbSetIpPortIdsStatus(ipPortIdsToMarkApproved, "Approved")
			}
		}
	}

	if *alertSubdomainFlag != false {

		domains := repo.DbViewHostUniqueDomains()
		sources := repo.DbViewHostUniqueSources()

		for _, domain := range domains {

			for _, source := range sources {

				hosts := repo.DbViewHost(domain, source)
				if len(hosts) == 0 {
					print("No Hosts to process alerts. Quitting.")
					os.Exit(0)
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
					print("No other hosts to compare against. Must be first run. Exit.")
					os.Exit(0)
				}

				olderTimestamp, newerTimestamp := logic.GetHostCompareTimestamps(hostMapTimestamp)

				print("Old hosts:", hostMapTimestamp[olderTimestamp])

				newHosts, hostsToApprove := logic.GetNewApproveHosts(hostMapTimestamp, olderTimestamp, newerTimestamp)
				missingHosts, hostsToInactivate := logic.GetMissingInactivateHosts(hostMapTimestamp, olderTimestamp, newerTimestamp)

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

				if len(newHosts) > 0 || len(missingHosts) > 0 {
					body := createEmailBodyFromAlertableHosts(newHosts, missingHosts)
					alert.SendMail(body, config.FromEmail, config.EmailPassword, config.EmailHost, config.EmailPort, config.Recipients)
				}

			}

		}

	}

	if *parseNmapDataFlag {
		print("Beginning XML Parse")
		mv, err := getMapFromXmlFile(config.NmapParseFile)
		if err != nil {
			print(err)
		}

		paths := mv.PathsForKey("host")
		hosts, err := mv.ValuesForPath(paths[0])

		// get host ports
		var hostPorts []HostPort
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
					hostPort := getPortFromXml(port, hostAddress, hostEndTime)
					hostPorts = append(hostPorts, hostPort)
				}
			case map[string]interface{}:
				hostPort := getPortFromXml(portsSlice, hostAddress, hostEndTime)
				hostPorts = append(hostPorts, hostPort)
			default:
				log.Fatal("unexpected type %T", v)
			}
		}

		// insert into sql lite
		for _, hostPort := range hostPorts {
			print("Inserting the following record:", hostPort)
			repo.DbCreateIpPort(hostPort.IP, hostPort.Protocol, hostPort.Port, hostPort.Service, hostPort.TimeScanned)
		}
	}

	// loop through every file in the todo directory and parse / insert
	// once completed, move each all to done and store the last completed file.
	if *parseSubdomainFlag != false {
		files, err := ioutil.ReadDir(config.SubdomainTodoDirectory)
		checkErr(err)

		for _, file := range files {
			hosts := parser.ReadLines(config.SubdomainTodoDirectory + file.Name())

			fileSplit := strings.Split(file.Name(), "_")
			domainName := strings.ToLower(fileSplit[0])
			dataSource := strings.ToLower(fileSplit[1])
			timeScanned := fileSplit[2]

			repo.DbInsertHostsForDomain(hosts, domainName, dataSource, timeScanned)
			err := ioutil.WriteFile(config.SubdomainLastCompleted, []byte(file.Name()), 0644)
			checkErr(err)

			err = common.CopyFile(config.SubdomainTodoDirectory+file.Name(), config.SubdomainTodoResolveDirectory+file.Name())
			err = os.Rename(config.SubdomainTodoDirectory+file.Name(), config.SubdomainDoneDirectory+file.Name())

			checkErr(err)
		}
	}

	// loop through every file in the todo subdomains resolve and insert IpHost
	if *parseHostIpFlag != false {
		files, err := ioutil.ReadDir(config.SubdomainDoneResolveDirectory)
		checkErr(err)

		for _, file := range files {
			ipHosts := parser.ReadLines(config.SubdomainDoneResolveDirectory + file.Name())

			fileSplit := strings.Split(file.Name(), "_")
			domainName := strings.ToLower(fileSplit[0])
			timeScanned := fileSplit[2]

			print("Inserting IPhosts: ", ipHosts, domainName)

			// create function to insert the ip hosts as IP_host into sqlite
			// once that is done create the alerting logic.
			// once that is done, automate the nmap scans for these IPs and automate eyewitness.
			repo.DbInsertIpHostsForDomain(ipHosts, domainName, timeScanned)

			err = os.Rename(config.SubdomainDoneResolveDirectory+file.Name(), config.SubdomainArchiveResolveDirectory+file.Name())
		}
	}

}

func containsIpPort(s []repo.IPPortDb, e repo.IPPortDb) bool {
	for _, a := range s {
		if a.Port == e.Port {
			return true
		}
	}
	return false
}

func getConfig(configLocation string) Config {
	file, e := ioutil.ReadFile(configLocation)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)
	return config
}

func getMapFromXmlFile(fileName string) (mxj.Map, error) {
	xmlValue := common.ReadFile(fileName)
	mv, err := mxj.NewMapXml(xmlValue)
	if err != nil {
		print(err)
	}
	return mv, nil
}

func getPortFromXml(portsSlice interface{}, hostAddress string, hostEndTime int64) HostPort {
	protocol := portsSlice.(map[string]interface{})["-protocol"].(string)
	portNumber := portsSlice.(map[string]interface{})["-portid"].(string)
	portService := portsSlice.(map[string]interface{})["service"].(map[string]interface{})["-name"].(string)
	portStatus := portsSlice.(map[string]interface{})["state"].(map[string]interface{})["-state"].(string)
	hostPort := HostPort{
		IP:          hostAddress,
		Protocol:    protocol,
		Port:        portNumber,
		State:       portStatus,
		TimeScanned: hostEndTime,
		Service:     portService,
	}
	return hostPort
}

func createEmailBodyFromAlertablePorts(addedPorts []repo.IPPortDb, missingPorts []repo.IPPortDb) string {
	body := "New port changes identified \n ================== \n"
	for _, port := range addedPorts {
		body += port.IP + ":" + port.Port + "\t [" + port.Protocol + "]" + "\t Added \n"
	}
	for _, port := range missingPorts {
		body += port.IP + ":" + port.Port + "\t [" + port.Protocol + "]" + "\t Missing \n"
	}
	print("Email Body", body)
	return body
}

func createEmailBodyFromAlertableHosts(newHosts []repo.HostDb, missingHosts []repo.HostDb) string {
	body := "New host changes identified \n ================== \n"
	for _, item := range newHosts {
		body += item.Host + "\t New \n"
	}
	for _, item := range missingHosts {
		body += item.Host + "\t Missing \n"
	}
	print("Email Body", body)
	return body
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
