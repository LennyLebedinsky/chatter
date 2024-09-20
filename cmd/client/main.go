package main

import (
	"log"
	"net"
	"net/http"
)

type config struct {
	Host string
	Port string
}

func serveClientPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "cmd/client/client.html")
}

func main() {
	config := &config{
		Host: "localhost",
		Port: "8081",
	}

	http.HandleFunc("/", serveClientPage)
	err := http.ListenAndServe(net.JoinHostPort(config.Host, config.Port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
