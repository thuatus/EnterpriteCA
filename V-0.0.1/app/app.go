package main

// Web interface application to manager PKI

import (
	"fmt"
	"html/template"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/ca"
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
	_, err := os.Stat("/home/alvaro/srv/CA/")
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

func ApplyInitialSettings(w http.ResponseWriter, r *http.Request) {
	// Apply the initial settings
	// If not, return an error
	// If yes, return nil
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Here you would process the form data and apply the initial settings
	// For now, we will just log the form data
	CaName := r.FormValue("ca_name")
	CaFolderName := strings.ReplaceAll(CaName, " ", "_")
	Country := r.FormValue("ca_country")
	State := r.FormValue("ca_state")
	Local := r.FormValue("ca_location")
	Organization := r.FormValue("ca_organization")
	OrganizationalUnit := r.FormValue("ca_ou")
	RootDomain := r.FormValue("ca_root_domain")
	log.Printf("Applying initial settings: CaName=%s, Country=%s, State=%s, Local=%s, Organization=%s, OrganizationalUnit=%s, RootDomain=%s",
		CaName, Country, State, Local, Organization, OrganizationalUnit, RootDomain)

	// Create the CA folders
	// This is where you would create the necessary directories for the CA
	caPathName, intermediateCaPathName, err := ca.CreatePKIFolders(CaFolderName)
	if err != nil {
		log.Println("Error creating PKI folders:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Folders created: CA Path=%s, Intermediate CA Path=%s", caPathName, intermediateCaPathName)

	// Create the CA configuration files
	_, err = ca.CreateConfigFiles(caPathName, Country, State, Local, Organization, OrganizationalUnit)
	if err != nil {
		log.Println("Error creating config files:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Config files created successfully")

	// Create CA certificate
	_, err = ca.CreateCACertificate(caPathName)
	if err != nil {
		log.Println("Error creating CA certificate:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Initial settings applied successfully")

	// Redirect to the main page or show a success message
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

		fmt.Println("Initial settings found, serving main page")
		// http.ServeFile(w, r, "templates/index.html")
	}))

	// Register the form for initial settings
	mux.Handle("/save_init_cfg/", Logging()(ApplyInitialSettings))

	// Start the server
	err := http.ListenAndServe(":8080", mux)
	// Log the error if any
	if err != nil {
		log.Fatal(err)
	}
}
