package main

// Web interface application to manager PKI

import (
	"html/template"
	"log"
	"log/syslog"
	"net/http"
	"os"
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

func checkInitialSettings() (bool, error) {
	// Check if the initial settings are set
	// If not, return false and an error
	// If yes, return true and nil
	_, err := os.Stat("/home/alvaro/srv/CA/CA-02/index.txt")
	if os.IsNotExist(err) {
		return false, err
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func FormInitialSettings(w http.ResponseWriter, _ *http.Request) {
	// Create the initial CA settings - form  - PKI and admin web console user
	// If not, return an error
	// If yes, return nil

	template, err := template.ParseFiles("templates/frm_init_ca.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	FrmInitialSettings := struct {
		Title              string
		CaName             string
		Country            string
		State              string
		Local              string
		Organization       string
		OrganizationalUnit string
		RootDomain         string
	}{
		Title:              "Initial Settings",
		CaName:             "Enterprite Root CA",
		Country:            "ES",
		State:              "Madrid",
		Local:              "Madrid",
		Organization:       "Enterprite",
		OrganizationalUnit: "Enterprite CA",
		RootDomain:         "enterprite.com",
	}
	if err := template.Execute(w, FrmInitialSettings); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux()
	// Register the logging middleware
	mux.HandleFunc("/", Logging()(func(w http.ResponseWriter, r *http.Request) {
		_, err := checkInitialSettings()
		if err != nil {
			log.Println("Initial settings not found, redirecting to form:", err)
			FormInitialSettings(w, r)
			return
		}
	}))
	// Start the server
	err := http.ListenAndServe(":8080", mux)
	// Log the error if any
	if err != nil {
		log.Fatal(err)
	}
}
