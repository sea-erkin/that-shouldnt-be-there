package main

import (
	"flag"
	"log"
	"os"
	"reflect"

	"github.com/sea-erkin/that-shouldnt-be-there/src/alert"
	"github.com/sea-erkin/that-shouldnt-be-there/src/common"
	"github.com/sea-erkin/that-shouldnt-be-there/src/dto"
	"github.com/sea-erkin/that-shouldnt-be-there/src/logic"
	"github.com/sea-erkin/that-shouldnt-be-there/src/repo"
)

var (
	configFlag         = flag.String("c", "", "Config file location")
	debugFlag          = flag.Bool("d", false, "Debug flag")
	alertPortFlag      = flag.String("alertPort", "", "Alert port logic flag")
	alertSubdomainFlag = flag.Bool("alertSubdomain", false, "Alert subdomain flag")
	parseNmapDataFlag  = flag.String("parseNmap", "", "Parse Nmap data flag")
	parseSubdomainFlag = flag.Bool("parseSubdomain", false, "Parse subdomain output. Must specify where to look for files")
	parseHostIpFlag    = flag.Bool("parseHostIp", false, "Parses host ip output.")
	approvalCodeFlag   = flag.String("approvalCode", "", "Approve UUID for alert code")
	sendScreenshotFlag = flag.Bool("sendScreenshot", false, "Sends screenshots of identified web hosts")
)

var (
	config dto.Config
)

func checkFlags() {
	flag.Parse()

	if *configFlag == "" {
		log.Fatal("Must supply a config file.")
	}

	if *alertPortFlag == "" && *parseNmapDataFlag == "" && *approvalCodeFlag == "" &&
		*parseSubdomainFlag == false && *alertSubdomainFlag == false && *parseHostIpFlag == false &&
		*sendScreenshotFlag == false {
		log.Fatal("Must specify something to do. Either alert, parse data, or approve an alert.")
	}
}

func main() {

	checkFlags()
	config = common.GetConfig(*configFlag)
	repo.InitializeDb(config.DbLocation)

	if *debugFlag {
		logic.Debug = true
	}

	if *approvalCodeFlag != "" {
		logic.ApproveApprovalCode(*approvalCodeFlag)
	}

	if *alertPortFlag != "" {
		print("Starting analysis for changes in open external ports")
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

		// loop through each IP address
		for ip, v := range ipPortIpSlice {

			// get the approved ports that exist for this IP address
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
					if logic.ContainsIpPort(newPorts, item) {
						samePorts = append(samePorts, item)
					} else {
						missingPorts = append(missingPorts, item)
					}
				}

				addedPorts := make([]repo.IPPortDb, 0)
				for _, item := range newPorts {
					if !logic.ContainsIpPort(approvedPorts, item) {
						addedPorts = append(addedPorts, item)
					}
				}

				// check the status too
				print("same ports", samePorts)
				print("approved ports", approvedPorts)
				if len(samePorts) == len(approvedPorts) {

					print("Ports are the same but have a newer timestamp. Updating entries")
					logic.SetPortIdsApproved(approvedPorts)
					logic.SetPortIdsInactive(samePorts)

					print("Removing alerts for the following IP if they exist: ", ip)
					repo.DbSetIpPortAlertStatusByIpAddress(ip)
					continue
				}

				print("New ports found: ", addedPorts)
				print("Ports missing:", missingPorts)
				logic.CreateAlertsForMissingAddedPorts(addedPorts, missingPorts)

				// Prep identified ports to take screenshot of new host ports
				alert.PrepNmapScreenshot(addedPorts, *alertPortFlag, config.ScreenshotTodoDirectory)

				// Prep and send E-mail code. Will eventually support multiple alert formats such as text & app push notifs
				body := alert.CreateEmailBodyFromAlertablePorts(addedPorts, missingPorts)
				alert.SendMail(body, config.FromEmail, config.EmailPassword, config.EmailHost, config.EmailPort, config.Recipients)

			} else {
				// no approved ports, mark the latest timestamp entries approved.
				logic.ApproveLatestIpPorts(v, *alertPortFlag, config.ScreenshotTodoDirectory)
			}
		}
	}

	if *alertSubdomainFlag != false {
		newHosts, missingHosts := logic.AlertSubdomains()

		if len(newHosts) > 0 || len(missingHosts) > 0 {
			body := alert.CreateEmailBodyFromAlertableHosts(newHosts, missingHosts)
			alert.SendMail(body, config.FromEmail, config.EmailPassword, config.EmailHost, config.EmailPort, config.Recipients)
		}
	}

	if *parseNmapDataFlag != "" {
		logic.CreateHostPortsFromNmapFile(*parseNmapDataFlag)
	}

	if *parseSubdomainFlag != false {
		logic.ParseSubdomains(config.SubdomainTodoDirectory, config.SubdomainTodoResolveDirectory, config.SubdomainDoneDirectory, config.SubdomainLastCompleted)
	}

	if *sendScreenshotFlag != false {
		logic.SendCompletedScreenshots(config.ScreenshotDoneDirectory, config.ScreenshotArchiveDirectory, config.FromEmail,
			config.EmailPassword, config.EmailHost, config.EmailPort, config.Recipients)
	}

	if *parseHostIpFlag != false {
		logic.CreateIpHostsFromResolvedSubdomains(config.SubdomainDoneResolveDirectory, config.SubdomainArchiveResolveDirectory)
	}

}
