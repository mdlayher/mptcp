// Command mptcphttp provides a very simple HTTP server, which detects if a
// client connection is using multipath TCP.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mdlayher/mptcp"
)

var (
	// host is the address to which the HTTP server is bound
	host string
)

func init() {
	// Set up flags
	flag.StringVar(&host, "host", ":8080", "HTTP server host")
}

func main() {
	log.SetPrefix("mptcphttp: ")

	// Parse flags
	flag.Parse()

	// Handle connections on root of HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if HTTP request is being issued from a client which is
		// connected using multipath TCP
		ok, err := mptcp.IsMPTCPHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Inform the client if they are connected with multipath TCP
		if ok {
			fmt.Fprintf(w, "YES")
		} else {
			fmt.Fprintf(w, "NO")
		}
	})

	// Bind HTTP server to host
	log.Println("binding to:", host)
	http.ListenAndServe(host, nil)
}
