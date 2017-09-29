package repo

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/sea-erkin/that-shouldnt-be-there/src/common"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db    *sql.DB
	print = fmt.Println
)

type IPDb struct {
	IPId        string
	IP          string
	DateScanned int64
}

type IPPortDb struct {
	IPPortId    int64
	Protocol    string
	IP          string
	Port        string
	Service     string
	IsApproved  bool
	Alert       bool
	DateScanned int64
}

type HostDb struct {
	HostId      int64
	Host        string
	Domain      string
	Status      bool
	Alert       bool
	TimeScanned int64
}

type IpHostDb struct {
	HostId      int64
	Host        string
	Domain      string
	IP          string
	Status      bool
	Alert       bool
	TimeScanned int64
}

func InitializeDb(dbfileName string) {
	conn, _ := sql.Open("sqlite3", dbfileName)
	db = conn
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func DbGetApprovedPorts(ipParam string) []IPPortDb {
	rows, err := db.Query(`SELECT IpPortId, Ip, Protocol, Port, Service, IsApproved, TimeScanned
			     		   FROM IPPort
						   WHERE IP = ?
						   AND IsApproved = 1`, ipParam)
	checkErr(err)

	var ipPorts []IPPortDb

	var ipPortId int64 = -1
	var ip string
	var protocol string
	var port string
	var service string
	var isApproved bool
	var timeScanned int64
	for rows.Next() {
		err = rows.Scan(&ipPortId, &ip, &protocol, &port, &service, &isApproved, &timeScanned)
		checkErr(err)

		ipPortDb := IPPortDb{
			IPPortId:    ipPortId,
			IP:          ip,
			Protocol:    protocol,
			Port:        port,
			Service:     service,
			IsApproved:  isApproved,
			DateScanned: timeScanned,
		}

		ipPorts = append(ipPorts, ipPortDb)
	}
	rows.Close()
	return ipPorts
}

func DbCreateHostAlert(hosts []HostDb) bool {

	var ids = make([]string, 0)
	for _, host := range hosts {
		ids = append(ids, strconv.FormatInt(host.HostId, 10))
	}

	comma := "'"
	updateString := strings.Join(ids, "','")
	finalString := comma + updateString + comma

	sqlString := "UPDATE Host SET Alert=1 WHERE HostId IN " + "(" + finalString + ")"

	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect > 0 {
		print("Rows approved: " + strconv.FormatInt(affect, 10))
	}

	return true
}

func DbCreateIpPortAlert(ipPortIds []string) bool {
	comma := "'"
	updateString := strings.Join(ipPortIds, "','")
	finalString := comma + updateString + comma

	sqlString := "UPDATE IPPort SET Alert=1 WHERE IPPortId IN " + "(" + finalString + ")"

	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect > 0 {
		print("Rows approved: " + strconv.FormatInt(affect, 10))
	}

	return true
}

func DbSetIpPortAlertStatusByIpAddress(ip string) bool {

	sqlString := "UPDATE IPPort SET IsApproved = -1 WHERE IP = ? AND Alert = 1"
	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec(ip)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect > 0 {
		print("Rows removed alert: " + strconv.FormatInt(affect, 10))
	}

	return true
}

func DbSetIpPortAlertApprovedByUuid(uuid string) bool {

	sqlString := `UPDATE ipport
	SET IsApproved = 1
	WHERE IpPortId = (
		SELECT ip.ipportid
		FROM IpPort ip
		JOIN AlertIPPort alert
			ON ip.IpPortId = alert.ipPortId
		WHERE UUID = ?
	)`

	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec(uuid)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect == 1 {
		print("Approved alert")
		return true
	} else {
		return false
	}

}

func DbSetHostAlertApprovedByUuid(uuid string) bool {

	sqlString := `UPDATE host
		SET Status = 1
		WHERE HostId = (
			SELECT host.HostId
			FROM Host host
			JOIN AlertHost alert
				ON host.HostId = alert.HostId
			WHERE UUID = ?
		)`

	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec(uuid)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect == 1 {
		print("Approved alert")
		return true
	} else {
		return false
	}

}

func DbCreateAlertIpPortUuid(ipPortId int64, uuid string) int64 {
	stmt, err := db.Prepare("INSERT INTO AlertIPPort(ipPortId, uuid) values(?,?)")
	checkErr(err)

	print(stmt)

	res, err := stmt.Exec(ipPortId, uuid)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

func DbCreateAlertHostUuid(hosts []HostDb) {

	for _, host := range hosts {

		uuid, err := common.NewUUID()
		checkErr(err)

		stmt, err := db.Prepare("INSERT INTO AlertHost(HostId, uuid) values(?,?)")
		checkErr(err)

		print(stmt)

		res, err := stmt.Exec(host.HostId, uuid)
		checkErr(err)

		id, err := res.LastInsertId()
		print(id)
		checkErr(err)

	}
}

func DbSetHostSubdomainStatus(hosts []HostDb, status string) bool {

	var ids = make([]string, 0)
	for _, host := range hosts {
		ids = append(ids, strconv.FormatInt(host.HostId, 10))
	}

	statuses := map[string]int{
		"Inactive": -1,
		"Approved": 1,
		"New":      0,
	}

	var statusCode int
	if _statusCode, ok := statuses[status]; ok {
		statusCode = _statusCode
	} else {
		log.Fatal("DB: Invalid status")
	}

	comma := "'"
	updateString := strings.Join(ids, "','")
	finalString := comma + updateString + comma

	var sqlString string
	switch statusCode {
	case -1:
		sqlString = "UPDATE Host SET Status=-1 WHERE HostId IN " + "(" + finalString + ")"
	case 0:
		sqlString = "UPDATE Host SET Status=0 WHERE HostId IN " + "(" + finalString + ")"
	case 1:
		sqlString = "UPDATE Host SET Status=1 WHERE HostId IN " + "(" + finalString + ")"
	}

	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect > 0 {
		print("Rows approved: " + strconv.FormatInt(affect, 10))
	}

	return true

}

func DbSetIpPortIdsStatus(ipPortIds []string, status string) bool {

	statuses := map[string]int{
		"Inactive": -1,
		"Approved": 1,
		"New":      0,
	}

	var statusCode int
	if _statusCode, ok := statuses[status]; ok {
		statusCode = _statusCode
	} else {
		log.Fatal("DB: Invalid status")
	}

	comma := "'"
	updateString := strings.Join(ipPortIds, "','")
	finalString := comma + updateString + comma

	var sqlString string
	switch statusCode {
	case -1:
		sqlString = "UPDATE IPPort SET IsApproved=-1 WHERE IPPortId IN " + "(" + finalString + ")"
	case 0:
		sqlString = "UPDATE IPPort SET IsApproved=0 WHERE IPPortId IN " + "(" + finalString + ")"
	case 1:
		sqlString = "UPDATE IPPort SET IsApproved=1 WHERE IPPortId IN " + "(" + finalString + ")"
	}

	print("SQL String:", sqlString)

	stmt, err := db.Prepare(sqlString)

	checkErr(err)

	res, err := stmt.Exec()
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	if affect > 0 {
		print("Rows approved: " + strconv.FormatInt(affect, 10))
	}

	return true
}

func DbViewHost(domainParam string) []HostDb {
	print(domainParam, "DOMAIN")
	rows, err := db.Query(`SELECT HostId, Host, Domain, Status, Alert, TimeScanned
		FROM Host
		WHERE Status IN (0,1)
		AND Alert = 0
		AND Domain = ?`, domainParam)
	checkErr(err)

	var hosts []HostDb

	var hostId int64 = -1
	var host string
	var domain string
	var status bool
	var alert bool
	var timeScanned int64
	for rows.Next() {
		err = rows.Scan(&hostId, &host, &domain, &status, &alert, &timeScanned)
		checkErr(err)

		hostDb := HostDb{
			HostId:      hostId,
			Host:        host,
			Domain:      domain,
			Status:      status,
			Alert:       alert,
			TimeScanned: timeScanned,
		}

		hosts = append(hosts, hostDb)
	}
	rows.Close()

	return hosts
}

func DbViewIpPort() []IPPortDb {
	rows, err := db.Query(`SELECT IpPortId, Ip, Protocol, Port, Service, IsApproved, TimeScanned
						   FROM IPPort
						   WHERE IsApproved = 0
						   AND Alert = 0`)
	checkErr(err)

	var ipPorts []IPPortDb

	var ipPortId int64 = -1
	var ip string
	var protocol string
	var port string
	var service string
	var isApproved bool
	var timeScanned int64
	for rows.Next() {
		err = rows.Scan(&ipPortId, &ip, &protocol, &port, &service, &isApproved, &timeScanned)
		checkErr(err)

		ipPortDb := IPPortDb{
			IPPortId:    ipPortId,
			IP:          ip,
			Protocol:    protocol,
			Port:        port,
			Service:     service,
			IsApproved:  isApproved,
			DateScanned: timeScanned,
		}

		ipPorts = append(ipPorts, ipPortDb)
	}
	rows.Close()
	return ipPorts
}

func DbCreateIpPort(ip string, protocol string, port string, service string, timeScanned int64) int64 {
	stmt, err := db.Prepare("INSERT INTO IPPort(IP, Protocol, Port, Service, TimeScanned) values(?,?,?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(ip, protocol, port, service, timeScanned)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

func DbInsertHostsForDomain(hosts []string, domainName string, timeString string) {

	timeScanned, err := strconv.ParseInt(timeString, 10, 64)

	checkErr(err)

	for _, host := range hosts {

		stmt, err := db.Prepare("INSERT INTO HOST(Domain, Host, TimeScanned) VALUES (?, ?, ?)")
		checkErr(err)

		res, err := stmt.Exec(domainName, host, timeScanned)
		checkErr(err)

		print("Host Inserted:", host)
		print("Rows Inserted:", res.RowsAffected)
	}
}

func DbInsertIpHostsForDomain(hosts []string, domainName string, timeString string) {

	timeScanned, err := strconv.ParseInt(timeString, 10, 64)

	checkErr(err)

	for _, host := range hosts {

		print(host)

		hostSplit := strings.Split(host, "_")

		print(hostSplit)

		stmt, err := db.Prepare("INSERT INTO IPHOST(Domain, Ip, Host, TimeScanned) VALUES (?, ?, ?, ?)")
		checkErr(err)

		res, err := stmt.Exec(domainName, hostSplit[0], hostSplit[1], timeScanned)
		checkErr(err)

		print("Host Inserted:", host)
		print("Rows Inserted:", res.RowsAffected)
	}
}
