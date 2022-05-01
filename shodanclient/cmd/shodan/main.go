package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"shodanclient/shodan"
	"strconv"
	"strings"
)

func main() {
	//if len(os.Args) != 2 {
	//	log.Fatalln("Usage: shodan searchterm")
	//}
	//apiKey := os.Getenv("SHODAN_API_KEY")

	// Declare variables.
	// It wasn't collecting from the env correctly. So I hardcoded it for the moment.
	apiKey := <API_Key>
	prompt := "shodan client> "
	reader := bufio.NewReader(os.Stdin)
	s := shodan.New(apiKey)
	stringPorts := []string{}

	info, err := s.APIInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("Query Credits: %d\nScan Credits: %d\n", info.QueryCredits, info.ScanCredits)

	for {
		fmt.Print(prompt)
		command, _ := reader.ReadString('\n')
		// Not sure if I need to add a janky os check here to work for unix systems.
		command = strings.Replace(command, "\r\n", "", -1)

		commandargv := strings.Join(strings.Split(command, " ")[0:1], "")

		switch commandargv {
		case "credits":
			info, err := s.APIInfo()
			if err != nil {
				log.Panicln(err)
			}
			fmt.Printf("Query Credits: %d\nScan Credits: %d\n", info.QueryCredits, info.ScanCredits)

		case "exit":
			fmt.Print("Exiting...")
			os.Exit(0)

		case "help":
			fmt.Print("Commands:\n")
			fmt.Print("help - Prints this help menu.\n")
			fmt.Print("credits - Returns number of query and scan credits remaining for this account.\n")
			fmt.Print("hostsearch [serach term] - Queries shodan based on the provided query. i.e. 'ftp'\n")
			fmt.Print("ipsearch [ip] - Retrieves info for a given ip adddress.\n")
			fmt.Print("exit - Exits the shodan client.\n")

		case "hostsearch":
			// Shave over first part remaing contnet to hostSearch(). Parsing remaing arguments should maybe be functionaized
			search := strings.TrimPrefix(command, "hostsearch ")
			hostSearch, err := s.HostSearch(search)
			if err != nil {
				log.Panicln(err)
			}

			//fmt.Printf("|%12s%8s%8s%8s", "IP", "Port", "|", "|\n")
			for _, host := range hostSearch.Matches {
				fmt.Printf("%18s%8d\n", host.IPString, host.Port)
			}

		case "ipsearch":
			// Shave off first part remaing contnet to hostSearch(). Parsing remaing arguments should maybe be functionaized
			search := strings.TrimPrefix(command, "ipsearch ")
			hostSearch, err := s.IPSearch(search)
			if err != nil {
				log.Panicln(err)
			}

			// Convert int slice to string.
			for i := range hostSearch.Ports {
				number := hostSearch.Ports[i]
				text := strconv.Itoa(number)
				stringPorts = append(stringPorts, text)
			}
			// Convert string slice to string
			hostname := strings.Join(hostSearch.HostNames2, " ")
			fmt.Printf("IPAddress: %s\n OS: %s\n ISP: %s\n Ports: %s\n Hostnames: %s\n Org: %s\n Domains: %s\n ASN: %s\n", hostSearch.IPString, hostSearch.OS, hostSearch.ISP, strings.Join(stringPorts, " "), hostname, hostSearch.Org2, strings.Join(hostSearch.Domains2, " "), hostSearch.ASN)

		default:
			fmt.Printf("Error: Unknown command, %s.\n", command)
		}
	}
}

// Future expansion, support more API endpoints.
