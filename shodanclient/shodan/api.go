package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIInfo struct {
	QueryCredits int    `json:"query_credits"`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet"`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
	MonitoredIPs int    `json:"monitored_ips"`
	UsageLimts   Limits `json:"usage_limits"`
}

type Limits struct {
	ScanCredits  int `json:"scan_credits"`
	QueryCredits int `json:"query_credits"`
	MonitoredIPs int `json:"monitored_ips"`
}

func (s *Client) APIInfo() (*APIInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", BaseURL, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.Status != "200" {
		fmt.Printf("An Error occured: %s\n", res.Status)
	}

	var ret APIInfo

	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}
	return &ret, err
}
