package main

// Web interface application to manager PKI

import (
	"log"
	"log/syslog"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// register the request
			start := time.Now()
			defer func() { log.Println(r.RemoteAddr, r.Method, r.URL.Path, time.Since(start)) }()
			logWriter, err := syslog.New(syslog.LOG_INFO, "EnterpriteCA")
			if err != nil {
				log.Printf("Failed to connect to syslog: %v", err)
			} else {
				log.SetOutput(logWriter)
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux()
	// Register the logging middleware
	mux.HandleFunc("/", Logging()(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	// Start the server
	err := http.ListenAndServe(":8080", mux)
	// Log the error if any
	if err != nil {
		log.Fatal(err)
	}
}
