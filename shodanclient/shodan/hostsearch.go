package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IPInfo struct {
	IPString   string   `json:"ip_string"`
	OS         string   `json:"os"`
	ISP        string   `json:"isp"`
	Ports      []int    `json:"ports"`
	HostNames2 []string `json:"hostnames"`
	Org2       string   `json:"org"`
	Domains2   []string `json:"domains"`
	ASN        string   `json:"asn"`
}

func (s *Client) IPSearch(q string) (*IPInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/%s?key=%s", BaseURL, q, s.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.Status != "200" {
		fmt.Printf("An Error occured: %s\n", res.Status)
	}

	var ret IPInfo
	if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
