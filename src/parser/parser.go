package parser

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/sea-erkin/that-shouldnt-be-there/src/common"
	"github.com/sea-erkin/that-shouldnt-be-there/src/dto"
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

func NmapGetMapFromXmlFile(fileName string) (mxj.Map, error) {
	xmlValue := common.ReadFile(fileName)
	mv, err := mxj.NewMapXml(xmlValue)
	if err != nil {
		print(err)
	}
	return mv, nil
}

func NmapGetPortFromXml(portsSlice interface{}, hostAddress string, hostEndTime int64) dto.HostPort {
	protocol := portsSlice.(map[string]interface{})["-protocol"].(string)
	portNumber := portsSlice.(map[string]interface{})["-portid"].(string)
	portService := portsSlice.(map[string]interface{})["service"].(map[string]interface{})["-name"].(string)
	portStatus := portsSlice.(map[string]interface{})["state"].(map[string]interface{})["-state"].(string)
	hostPort := dto.HostPort{
		IP:          hostAddress,
		Protocol:    protocol,
		Port:        portNumber,
		State:       portStatus,
		TimeScanned: hostEndTime,
		Service:     portService,
	}
	return hostPort
}
