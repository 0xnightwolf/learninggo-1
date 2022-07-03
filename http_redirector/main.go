// Extension ideas. Staged payload for Meterpreter and HTTPS. Also arguments to set the proxyurls would be good. So meterpreter revers http is scuffed. Unstaged and staged TCP work but not staged http. I'll try Empire instead
// Current TODO: Add options for listening address and port. Add HTTP/HTTPS switch.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

var (
	redirect_address1  string
	redirect_address2  string
	redirect_host1     string
	redirect_host2     string
	ignore_self_signed bool
	port               string

	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

// This should probalby be replace to pull an address based on a defualt interface and or ask for one to be specified. Or maybe the default is on all interfaces.
func init() {

	flag.StringVar(&redirect_address1, "addr1", "", "First redirect address. <ipaddress:port>")
	flag.StringVar(&redirect_host1, "host1", "", "First redirect hostname")
	flag.StringVar(&redirect_address2, "addr2", "", "Second redirect address. <ipaddress:port>")
	flag.StringVar(&redirect_host2, "host2", "", "Second redirect hostname")
	flag.BoolVar(&ignore_self_signed, "verify", false, "Skip certificate validation, valse by default")
	flag.StringVar(&port, "port", "80", "Port to listen on. Defaults to 80.")
	flag.Parse()

	str1 := []string{"http://", redirect_address1}
	hostProxy[redirect_host1] = strings.Join(str1, "")
	str2 := []string{"http://", redirect_address2}
	hostProxy[redirect_host2] = strings.Join(str2, "")

	for k, v := range hostProxy {
		remote, err := url.Parse(v)
		if err != nil {
			log.Fatal("Unable to parse proxy target")
		}
		proxies[k] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	// TODO: Test to make sure flag parseing works forhostname and redirect address. Add flag to switch betwen http and https mode. Then maybe staged payloads. HTTPS code is working just commented out currently in main.

	fmt.Println(redirect_host1)
	fmt.Println(redirect_host2)

	// This can be used to allow the usage of self signed certs that would otherwise be deemed invalid.
	if ignore_self_signed == true {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	r := mux.NewRouter()
	for host, proxy := range proxies {
		r.Host(host).Handler(proxy)

		//log.Fatal(http.ListenAndServeTLS(":80", "server.pem", "server.key", r))
		log.Fatal(http.ListenAndServe(strings.Join([]string{":", port}, ""), r))

	}
}
