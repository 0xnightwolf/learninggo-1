package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	porterrormsg = "Invalid Port specifications"
)

func dashSplit(sp string, ports *[]int) error {
	dp := strings.Split(sp, "-")
	if len(dp) != 2 {
		return errors.New(porterrormsg)
	}
	start, err := strconv.Atoi(dp[0])
	if err != nil {
		return errors.New(porterrormsg)
	}
	end, err := strconv.Atoi(dp[1])
	if err != nil {
		return errors.New(porterrormsg)
	}
	if start > end || start < 1 || end > 65535 {
		return errors.New(porterrormsg)
	}
	for ; start <= end; start++ {
		*ports = append(*ports, start)
	}
	return nil
}

func convertAndAddPort(p string, ports *[]int) error {
	i, err := strconv.Atoi(p)
	if err != nil {
		return errors.New(porterrormsg)
	}
	if i < 1 || i > 65535 {
		return errors.New(porterrormsg)
	}
	*ports = append(*ports, i)
	return nil
}

// Parse strings seperated by '-' or ',' and return a slice of Ints
func Parse(s string) ([]int, error) {
	ports := []int{}
	if strings.Contains(s, ",") && strings.Contains(s, "-") {
		sp := strings.Split(s, ",")
		for _, p := range sp {
			if strings.Contains(p, "-") {
				if err := dashSplit(p, &ports); err != nil {
					return ports, err
				}
			} else {
				if err := convertAndAddPort(p, &ports); err != nil {
					return ports, err
				}
			}
		}
	} else if strings.Contains(s, ",") {
		sp := strings.Split(s, ",")
		for _, p := range sp {
			convertAndAddPort(p, &ports)
		}
	} else if strings.Contains(s, "-") {
		if err := dashSplit(s, &ports); err != nil {
			return ports, err
		}
	}
	return ports, nil
}

func worker(target string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", target, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			// Write result as nothing.
			results <- 0
			continue
		}
		conn.Close()
		// Write open port to result
		results <- p
	}
}

func main() {
	// Define variables
	var scan_target string
	var openports []int
	var threads int
	var scan_ports string

	// Define flags
	flag.StringVar(&scan_target, "ip", "", "Specify target to scan.")
	flag.IntVar(&threads, "t", 100, "Number of threads.")
	flag.StringVar(&scan_ports, "p", "", "Ports to scan (Default 1-1024")

	// Parse flags
	flag.Parse()

	// Define port and results channels.
	ports := make(chan int, threads)
	results := make(chan int, 65535)

	parsed_ports, error := Parse(scan_ports)
	if error != nil {
		fmt.Printf("Error: %s", error)
	}

	if scan_target == "" {
		fmt.Printf("No target specified to scan")
		os.Exit(1)
	}

	// Start workers up to cap on number of threads
	for i := 0; i < cap(ports); i++ {
		go worker(scan_target, ports, results)
	}

	// Create list of ports to scan.
	if scan_ports == "" {
		go func() {
			for i := 0; i < 1024; i++ {
				ports <- i
			}
		}()
	} else {
		go func() {
			for _, i := range parsed_ports {
				ports <- i
			}
		}()
	}

	// Collect the results from the scan.
	if scan_ports != "" {
		for i := range parsed_ports {
			// Another garbage assignment. The loop needs it but we don't acually use it's value anywhere.
			_ = i
			port := <-results
			if port != 0 {
				openports = append(openports, port)
			}
		}
	} else {
		for i := 0; i < 1024; i++ {
			ports := <-results
			if ports != 0 {
				openports = append(openports, ports)
			}
		}
	}

	// Close channels
	close(ports)
	close(results)
	// Sort the port list. (They probalby didn't come back in order and this meakes it more human readable.)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}

// Future expansion ideas, add parsing for lists and rangs of IP addresses, ping sweeps, host os detection
