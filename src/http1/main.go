package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {

}

func process(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
	log.Println(r.Form)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	log.WithFields(log.Fields{
		"hostname": "127.0.0.1:8080",
	}).Info("server start......")

	http.HandleFunc("/process", process)
	err := server.ListenAndServe()
	if err != nil {
		log.WithFields(log.Fields{
			"hostname": "127.0.0.1:8080",
		}).Info("server stop.")
	}
}
