package main

// Web interface application to manager PKI

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/ca"
	"github.com/thuatus/EnterpriteCA/tree/main/V-0.0.1/db"
	"golang.org/x/crypto/bcrypt"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var (
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

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
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println("Redirecting to first login page")
	http.Redirect(w, r, "/static/first_login.html", http.StatusSeeOther)

}

func createAdminUser(usr string, passwd string) error {
	// Create the admin user for the web console
	// This is where you would create the necessary user for the web console
	// For now, we will just log the action
	log.Println("Conecting to the database...")
	dbConn, err := db.ConnectDB()

	if err != nil {
		log.Println("Error connecting to the database:", err)
		return err
	}
	defer dbConn.Close()

	//  create the admin user
	_, err = dbConn.Exec("INSERT INTO ca.users (name, passwd, active ) VALUES (?, ?, ?)", usr, passwd, 1)
	if err != nil {
		log.Println("Error inserting admin user into the database:", err)
		return err
	}

	return nil
}

func handleAdminUser(w http.ResponseWriter, r *http.Request) {
	// Handle the admin user creatio
	if r.Method != http.MethodPost {
		// Serve the form (first_login.html)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed for admin user creation")
		return
	} else {

		// Get the form values
		username := r.FormValue("user")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")
		// Validate the form values
		if username == "" || password == "" || confirmPassword == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			log.Println("All fields are required for admin user creation")
			return
		}

		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			log.Println("Passwords do not match for admin user creation")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create the admin user
		err = createAdminUser(username, string(hash))
		if err != nil {
			log.Println("Error creating admin user:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		log.Println("Admin user created successfully")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

// authenticate user credentials
func authenticateUser(w http.ResponseWriter, r *http.Request) {
	// Handle the login request
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// No need to generate a hash here; we'll compare the password with the stored hash below

	// Validate the form values
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// check lengh of password
	if len(password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	var storedHash string
	err = dbConn.QueryRow("SELECT passwd FROM ca.users WHERE name = ?", username).Scan(&storedHash)
	if err != nil {
		log.Println("Error querying user:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		log.Println("Invalid credentials:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// check if user exists and is active
	var active int
	err = dbConn.QueryRow("SELECT active FROM ca.users WHERE name = ?", username).Scan(&active)
	if err != nil {

		log.Println("Error querying user active status:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if active == 0 {
		log.Println("User is not active")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("User %s authenticated successfully", username)

	// Create a new session
	session, err := store.New(r, "session-token")
	if err != nil {
		log.Println("Error getting session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Store the user ID in the session
	bytes := make([]byte, 16)
	_, err = rand.Read(bytes) // Fill the slice with cryptographically secure random bytes
	if err != nil {
		log.Println("Error generating random UUID:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["UUID"] = hex.EncodeToString(bytes)
	session.Values["Authenticated"] = true
	session.Values["LoginTime"] = time.Now().Format(time.RFC3339)
	session.Values["UserAgent"] = r.UserAgent()

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Session created for user %s", username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// checkSession checks if the user is authenticated and the session is valid
func checkSession() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			// Check if the session exists
			session, err := store.Get(r, "session-token")
			if err != nil {
				log.Println("Error getting session:", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			_, ok := session.Values["UUID"]
			if !ok {

				log.Println("Session not found or user not authenticated")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			//check if token is expired based on max age
			loginTimeStr, ok := session.Values["LoginTime"].(string)
			if !ok {
				log.Println("Login time not found in session")
				http.Redirect(w, r, "/login", http.StatusSeeOther)

				return
			}
			loginTime, err := time.Parse(time.RFC3339, loginTimeStr)
			if err != nil {
				log.Println("Error parsing login time:", err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			maxAge := 60 * time.Minute // Set the maximum age for the session
			if time.Since(loginTime) > maxAge {
				log.Println("Session expired, redirecting to login")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

func ChainMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// handle the certificate issuance request form
func IssueCertHandler(w http.ResponseWriter, r *http.Request) {

	// RETURN TEMPLATE FOR THE FORM
	template, err := template.ParseFiles("templates/frm_add_cert.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	IntCAInfo, err := ca.GetIntermediateCAInfo()
	if err != nil {
		log.Println("Error getting intermediate CA info:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, IntCAInfo); err != nil {
		log.Println("Error executing frm_add_cert template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

// Issue server certificate and save it to the database
func CreateServerCert(w http.ResponseWriter, r *http.Request) {
	// Handle the server certificate creation request
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the form values
	serverName := r.FormValue("serverName")

	// Validate the form values

	if serverName == "" {
		http.Error(w, "Server name is required", http.StatusBadRequest)
		return
	}

	if len(serverName) < 3 {
		http.Error(w, "Server name must be at least 3 characters long", http.StatusBadRequest)
		return
	}
	if len(serverName) > 64 {
		http.Error(w, "Server name must be at most 64 characters long", http.StatusBadRequest)
		return
	}
	// verify that server name contains domains .com .org .gov .jus
	validDomain := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+(com|gov|jus|org)(\.[a-zA-Z]{2})?$`)
	if !validDomain.MatchString(serverName) {
		http.Error(w, "Server name must contain a domain", http.StatusBadRequest)
		return
	}

	//generate a csr file for the server certificate
	csr, err := ca.GenerateCSR(serverName)
	if err != nil {
		log.Println("Error generating CSR:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	csrFile, err := os.Open(csr)
	if err != nil {
		log.Println("Error opening CSR file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer csrFile.Close()
	// Create the server certificate
	// cert, err := ca.IssueServerCertificate(serverName, csrFile)
	_, err = ca.IssueServerCertificate(serverName, csrFile)
	if err != nil {
		log.Println("Error creating server certificate:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Save the certificate to the database
	cert, err := ca.GetServerCertificateInfo(serverName)
	if err != nil {
		log.Println("Error getting server certificate info:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Connect to the database
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer dbConn.Close()

	// Save the certificate to the database
	err = dbConn.Ping()
	if err != nil {
		log.Println("Error pinging the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a dbCertInfo struct to hold the certificate information
	strPublicKey := hex.EncodeToString(cert.PublicKey)
	if strPublicKey == "" {
		log.Println("Error encoding public key to hex string")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// use formatted as your string in the desired pattern

	dbCertInfo := db.DBCertInfo{
		Issuer:    cert.Issuer,
		Subject:   cert.Subject,
		PublicKey: strPublicKey,
		Expire:    expireTime, // Pass as time.Time
	}
	err = db.AddCertificate(dbCertInfo)
	if err != nil {
		log.Println("Error saving certificate to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Server certificate created and saved successfully")
	// ****>change to view the certificate page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// main function to start the web server
func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux()
	// Register the logging middleware
	mux.HandleFunc("/", Logging()(checkSession()(func(w http.ResponseWriter, r *http.Request) {
		_, err := checkInitialSettings()
		if err != nil {
			log.Println("Initial settings not found, redirecting to form:", err)
			FormInitialSettings(w, r)

			log.Println("Redirecting to first login page")

			// http.ServeFile(w, r, "./static/first_login.html")
			return
		}

		fmt.Println("Initial settings found, serving main page")
		http.ServeFile(w, r, "index.html")
	})))

	// file Server to serve static files
	mux.HandleFunc("/static/", Logging()(func(w http.ResponseWriter, r *http.Request) {

		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)

	}))

	// Register the form for initial settings
	mux.Handle("/save_init_cfg/", Logging()(ApplyInitialSettings))

	// Register the form for creating the admin user
	mux.Handle("/add_user/", Logging()(handleAdminUser))

	// REgister the form for login
	mux.Handle("/login/", Logging()(func(w http.ResponseWriter, r *http.Request) {

		// Serve the login form
		http.ServeFile(w, r, "./static/login.html")

	}))

	// Register the login handler
	mux.Handle("/authUser/", Logging()(authenticateUser))

	// Register the form for issuing certificates
	mux.Handle("/issue_cert/", (ChainMiddleware(IssueCertHandler, Logging(), checkSession())))

	mux.Handle("/add_server_cert/", (ChainMiddleware(CreateServerCert, Logging(), checkSession())))

	// Start the server
	err := http.ListenAndServe(":8080", mux)
	// Log the error if any
	if err != nil {
		log.Fatal(err)
	}
}
