package dto

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
	ScreenshotDoneDirectory          string   `json:"screenshotDoneDirectory"`
	ScreenshotTodoDirectory          string   `json:"screenshotTodoDirectory"`
	ScreenshotArchiveDirectory       string   `json:"screenshotArchiveDirectory"`
	NmapTodoDirectory                string   `json:"nmapTodoDirectory"`
	EmailHost                        string   `json:"emailHost"`
	EmailPort                        string   `json:"emailPort"`
	FromEmail                        string   `json:"fromEmail"`
	EmailPassword                    string   `json:"emailPassword"`
	Recipients                       []string `json:"recipients"`
	ErrorRecipients                  []string `json:"errorRecipients"`
}

type Map map[string]interface{}
