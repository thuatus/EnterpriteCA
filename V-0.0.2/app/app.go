// Package main implements a web interface application to manage a Public Key Infrastructure (PKI).
//
// This application provides functionality for initializing a Certificate Authority (CA),
// managing users, issuing and revoking certificates, and viewing certificate information
// through a web interface. It uses Gorilla sessions for session management, bcrypt for
// password hashing, and interacts with a database for storing user and certificate data.
//
// v 0.0.2 - The main features include:
//
//   - Initial CA setup and configuration via web forms.
//   - Admin user creation and authentication.
//   - Session management with expiration checks.
//   - Issuing, viewing, and revoking server certificates.
//   - Secure handling of private keys and certificate data.
//   - Logging of requests and actions for auditing purposes.
//
// The application is structured around HTTP handlers, middleware for logging and session
// validation, and utility functions for PKI and database operations.
package main

// Web interface application to manager PKI

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"github.com/thuatus/EnterpriteCA/V-0.0.2/ca"
	"github.com/thuatus/EnterpriteCA/V-0.0.2/db"
	"golang.org/x/crypto/bcrypt"
)

// Structs
type Middleware func(http.HandlerFunc) http.HandlerFunc

type ServerKey struct {
	ServerKey string `json:"Server_key"`
}

type CheckSetup struct {
	InitialSettings bool `json:"initial_settings"`
}

type CertInputData struct {
	ServerName string `form:"serverName" binding:"required,min=3,max=64,fqdn"`
}

type DBCertInfo struct {
	Issuer    string
	Subject   string
	PublicKey string
	Created   time.Time
	Expire    time.Time
	Status    int // 1 for valid, 0 for invalid
}
type Certificate struct {
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	PublicKey string `json:"public_key"`
	Created   string `json:"created"`
	Expire    string `json:"expire"`
	Status    int    `json:"status"`
}

var (
	key        = []byte("super-secret-key")
	store      = sessions.NewCookieStore(key)
	rootCAPath = "/home/alvaro/srv/CA"
	//rootCAPath = "/srv/CA"
	appCertPath = "/srv/ssl/app.crt"
	appKeyPath  = "/srv/ssl/app.key"
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
			log.SetOutput(os.Stdout)

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Check if the initial configuration was made before
func checkInitialSettings() bool {
	// Check if the initial settings are set
	// If not, return false and an error
	// If yes, return true and nil
	var folderConfig CheckSetup
	folderConfig.InitialSettings = true
	_, err := os.Stat(rootCAPath)
	if os.IsNotExist(err) {
		log.Println("CA folder does not exist, initial settings not applied")
		folderConfig.InitialSettings = false
	}
	if err != nil {
		log.Println("Error checking CA folder:", err)
		folderConfig.InitialSettings = false
	}
	// set tlsEnabled to true if the CA folders exist

	return folderConfig.InitialSettings
}

// Handle initiation CA process, based on user input
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

// Handle the initial setting of PKI configuration
func ApplyInitialSettings(w http.ResponseWriter, r *http.Request) error {
	// Apply the initial settings
	// If not, return an error
	// If yes, return nil

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
	caPathName, intermediateCaPathName, err := ca.CreatePKIFolders(rootCAPath, CaFolderName)
	if err != nil {
		log.Println("Error creating PKI folders:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	log.Printf("Folders created: CA Path=%s, Intermediate CA Path=%s", caPathName, intermediateCaPathName)

	// Create the CA configuration files
	_, err = ca.CreateConfigFiles(caPathName, Country, State, Local, Organization, OrganizationalUnit)
	if err != nil {
		log.Println("Error creating config files:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	log.Println("Config files created successfully")

	// Create CA certificate
	_, err = ca.CreateCACertificate(caPathName)
	if err != nil {
		log.Println("Error creating CA certificate:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	fmt.Println("Initial settings applied successfully")

	return nil

	// Redirect to the main page or show a success message
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	//	fmt.Println("Redirecting to first login page")
	//	http.Redirect(w, r, "/static/first_login.html", http.StatusSeeOther)

}

// Handle the user input for create a admin account
func handleAdminUser(w http.ResponseWriter, r *http.Request) error {
	// Handle the admin user creatio
	// Get the form values
	username := r.FormValue("user")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")
	// Validate the form values
	if username == "" || password == "" || confirmPassword == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		log.Println("All fields are required for admin user creation")

	}

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		log.Println("Passwords do not match for admin user creation")

	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	// Create the admin user
	err = db.AddUser(username, string(hash))
	if err != nil {
		log.Println("Error creating admin user:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err

	}

	log.Println("Admin user created successfully")
	//http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

// Handle Login request
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the login form
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed for login")
		return
	}

	// Serve the login page
	http.ServeFile(w, r, "./static/login.html")
	log.Println("Serving login page")
}

// Authenticate user credentials
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

	// Create a new session based on JWT
	bytes := make([]byte, 16)
	_, err = rand.Read(bytes) // Fill the slice with cryptographically secure random bytes
	if err != nil {
		log.Println("Error generating random UUID:", err)
		return
	}

	claims := jwt.MapClaims{
		"UUID":          hex.EncodeToString(bytes), // Generate a UUID for the user session
		"Authenticated": true,
		"LoginTime":     time.Now().Format(time.RFC3339),
		"UserAgent":     r.UserAgent(),
		"Secure":        false, // Set to true if using HTTPS
		"HttpOnly":      true,
	}

	// Set the JWT secret in an environment variable before running the application
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println("Error signing JWT token:", err)

		return
	}
	// Set the JWT in a session cookie
	session, err := store.Get(r, "session-token")
	if err != nil {
		log.Println("Error getting session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session.Values["Auth_Token"] = tokenString // Store the JWT in the session
	session.Values["Authenticated"] = true
	session.Values["LoginTime"] = time.Now().Format(time.RFC3339)
	session.Values["UserAgent"] = r.UserAgent()
	session.Values["Secure"] = false // Set to true if using HTTPS
	session.Values["HttpOnly"] = true

	// Save the session
	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Session created for user %s with UUID %s", username, claims["UUID"])
	// Set the JWT in the response header
	w.Header().Set("Authorization", "Bearer "+tokenString)
	log.Println("JWT token created and stored in session for user:", username)

}

// Verif if the initial settings were made
func checkInitialSettingsMiddleware() Middleware {
	// Middleware to check if the initial settings have been applied
	return func(f http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			// Check if the initial settings have been applied
			//check if the initial settings were made
			checkInitSet := checkInitialSettings()
			if !checkInitSet {
				FormInitialSettings(w, r)
				log.Println("Redirecting to first login page")
				fmt.Println("Redirecting to first login page")
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// CheckSession checks if the user is authenticated and the session is valid
func checkSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//get session
		session, err := store.Get(c.Request, "session-token")
		if err != nil {
			log.Println("Error getting session:", err)
			c.Redirect(http.StatusSeeOther, "/login/")
			c.Abort()
			return
		}

		// Check if the jwt token exists in the session
		tokenString, ok := session.Values["Auth_Token"].(string)
		if !ok || tokenString == "" {
			log.Println("No Authorization token found in session, redirecting to login")
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort()
			return
		}
		// Strip "Bearer " prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}
		// Check if the session exists
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Println("Unexpected signing method:", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// set the JWT scretin a environment variable before running the application
			return []byte(os.Getenv("JWT_SECRET")), nil // Return the secret key used for signing
		})

		if err != nil || !token.Valid {
			log.Println("Error parsing JWT token:", err)
			c.Redirect(http.StatusSeeOther, "/login/")
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("UUID", claims["UUID"])                   // Set the UUID in the context
		c.Set("Authenticated", claims["Authenticated"]) // Set the authentication status in the context
		c.Set("LoginTime", claims["LoginTime"])         // Set the login time in the
		c.Set("UserAgent", claims["UserAgent"])         // Set the user agent in the context
		log.Println("Session is valid, user authenticated")

		c.Next() // Call the next handler in chain
	}
}

// Apply Midware chain of functions to enforce log, check auth and the execution of especified function
func ChainMiddleware(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// Handle the certificate issuance request form
func IssueCertHandler(w http.ResponseWriter, r *http.Request) error {

	// RETURN TEMPLATE FOR THE FORM
	template, err := template.ParseFiles("templates/frm_add_cert.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	IntCAInfo, err := ca.GetIntermediateCAInfo()
	if err != nil {
		log.Println("Error getting intermediate CA info:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	if err := template.Execute(w, IntCAInfo); err != nil {
		log.Println("Error executing frm_add_cert template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	return nil

}

// Issue server certificate and save it to the database
func CreateServerCert(c *gin.Context) error {
	// Handle the server certificate creation request

	// Get the form values
	//bind servername from the form
	var issueCertData CertInputData

	// Bind the form data to the issueCertData struct
	err := c.ShouldBind(&issueCertData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	// get the server name from the request data
	serverName := issueCertData.ServerName

	/*
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
	*/

	//generate a csr file for the server certificate
	csr, err := ca.GenerateCSR(serverName, rootCAPath)
	if err != nil {
		log.Println("Error generating CSR:", err)

		return err
	}

	csrFile, err := os.Open(csr)
	if err != nil {
		log.Println("Error opening CSR file:", err)

		return err
	}
	defer csrFile.Close()
	// Create the server certificate
	// cert, err := ca.IssueServerCertificate(serverName, csrFile)
	_, err = ca.IssueServerCertificate(serverName, csrFile)
	if err != nil {
		log.Println("Error creating server certificate:", err)

		return err
	}
	log.Println("Server certificate created successfully for:", serverName)

	//Save the certificate to the database

	cert, err := ca.GetServerCertificateInfo(serverName)
	if err != nil {
		log.Println("Error getting server certificate info:", err)

		return err
	}
	log.Println("Cert Objec:", cert.Subject)

	// Connect to the database
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)

		return err
	}
	defer dbConn.Close()

	// Save the certificate to the database
	err = dbConn.Ping()
	if err != nil {
		log.Println("Error pinging the database:", err)

		return err
	}

	// Create a dbCertInfo struct to hold the certificate information
	strPublicKey := hex.EncodeToString(cert.PublicKey)
	if strPublicKey == "" {
		log.Println("Error encoding public key to hex string")

		return err
	}

	expire := cert.Expire

	//parse string to date type
	expireTime, err := time.Parse("2006-01-02 15:04:05", expire)
	log.Println("Parsed Expire Time:", expireTime)
	if err != nil {
		log.Println("Error parsing expiration time:", err)

		return err
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

		return err
	}

	log.Println("Server certificate created and saved successfully")
	// ****>change to view the certificate page

	// Redirect to the view certificate page

	return nil

}

// Handle the view ca form, set up a template view
func ViewCertificateForm(c *gin.Context) error {
	// Serve the view certificate form
	/*
		template, err := template.ParseFiles("templates/frm_view_ca.html")

		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}
	*/
	// get database info about servers certificates
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Println("Error connecting to the database:", err)

		return err
	}
	defer dbConn.Close()
	// Get the list of server certificates from the database
	ServerCerts, err := db.GetServerCertificates()
	if err != nil {
		log.Println("Error getting server certificates from database:", err)
		return err
	}

	// Set JSON response content
	var certs []Certificate
	for _, cert := range ServerCerts {
		certs = append(certs, Certificate{
			Issuer:    cert.Issuer,
			Subject:   cert.Subject,
			PublicKey: cert.PublicKey,
			Created:   cert.Created.Format("2006-01-02 15:04:05"),
			Expire:    cert.Expire.Format("2006-01-02 15:04:05"),
			Status:    cert.Status,
		})
	}

	// Return the certificates as a JSON response.
	c.JSON(http.StatusOK, gin.H{"certificates": certs})

	return nil
}

// Handle the view pserver private key request
func ViewServerPrivateKey(w http.ResponseWriter, r *http.Request) (string, error) {
	// Handle the request to view the server private key

	// Extract the server name from the form

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {

		log.Println("Error decoding request body:", err)
		return "", err
	}

	// get the subject value from the request data
	subject := requestData["subject"]
	if subject == "" {

		log.Println("Subject is empty in the request data")
		return "", err
	}

	serverName, err := ca.GetServerNameFromCert(subject)
	if err != nil {
		log.Println("Error getting server name from certificate:", err)
		return "", err
	}

	// Get the private key for the server certificate
	privateKey, err := ca.GetServerPrivateKey(serverName)
	if err != nil {
		log.Println("Error getting server private key:", err)
		return "", err
	}

	//Serve the private key file
	/*
		responseData := ServerKey{ServerKey: privateKey}
		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(responseData)
	*/
	return privateKey, nil

}

// Handle CErtificate Revocation Request
func HandleRevokeCertificate(w http.ResponseWriter, r *http.Request) error {
	// Handle the certificate revocation request

	// Extract the certificate subject from the form
	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error on process request", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return err
	}

	subject := requestData["subject"]
	if subject == "" {
		http.Error(w, "Subject cant be empty", http.StatusBadRequest)
		log.Println("Subject is empty in the request data")
		return err
	}

	serverName, err := ca.GetServerNameFromCert(subject)
	if err != nil {
		log.Println("Error getting server name from certificate:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	err = ca.RevokeServerCertificate(serverName)
	if err != nil {
		log.Println("Error revoking certificate:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	log.Println("Certificate revoked successfully for subject:", subject)

	err = db.UpdateValidCertificate(subject)
	if err != nil {
		log.Println("Error updating valid certificate in database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	/*
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Certificate revoked successfully")

		log.Println("Certificate revoked successfully for subject:", subject)
		http.Redirect(w, r, "/view_cert/", http.StatusSeeOther)
	*/
	return nil
}

// main function to start the web server
// v0.0.2 Use gin framework to handle the web server
func main() {

	// call gin framework to handle the web server
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// log customization
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		return fmt.Sprintf("%s - %s %s %d %s %s %s %s %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.StatusCode,
			param.Path,
			param.Request.UserAgent(),
			param.Request.Proto,
			param.Latency,
			param.ErrorMessage,
		)
	}))

	// Set up the routes
	// check if the initial settings were made and handle the first login or main page
	// Register middleware for "/" route only if initial settings are present
	if checkInitialSettings() {
		r.GET("/", checkSessionMiddleware(), func(c *gin.Context) {
			c.File("index.html")
			fmt.Println("Serving main page")
		})
	} else {
		r.GET("/", func(c *gin.Context) {
			FormInitialSettings(c.Writer, c.Request)
		})
	}

	// Serve static content
	r.GET("/static/*filepath", func(c *gin.Context) {
		// Serve static files
		filepath := c.Param("filepath")
		if strings.HasPrefix(filepath, "/") {
			filepath = filepath[1:] // Remove leading slash if present
		}
		c.File("./static/" + filepath)
		log.Println("Serving static content:", filepath)

	})

	// Handle the initial settings form submission
	r.POST("/save_init_cfg/", func(c *gin.Context) {
		// Apply the initial settings
		err := ApplyInitialSettings(c.Writer, c.Request)
		if err != nil {
			log.Println("Error applying initial settings:", err)
			//return json error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		log.Println("Initial settings applied successfully, redirecting to first login page")
		c.Redirect(http.StatusSeeOther, "/static/first_login.html")

	})

	// Handle the admin user creation form submisison
	r.POST("/add_user/", func(c *gin.Context) {
		err := handleAdminUser(c.Writer, c.Request)
		if err != nil {
			log.Println("Error handling admin user creation:", err)
			//return json error response
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		// Redirect to the main page
		log.Println("Admin user created successfully, redirecting to main page")
		c.Redirect(http.StatusSeeOther, "/")

	})

	r.GET("/login/", func(c *gin.Context) {
		// Serve the login form

		// Serve the login page
		c.File("./static/login.html")

		log.Println("Serving login page")
	})

	r.POST("/authUser/", func(c *gin.Context) {
		// Authenticate the user
		authenticateUser(c.Writer, c.Request)
		log.Println("User authenticated successfully, redirecting to main page")
		c.Redirect(http.StatusSeeOther, "/")
	})

	r.GET("/issue_cert/", func(c *gin.Context) {
		// Handle the certificate issuance request

		err := IssueCertHandler(c.Writer, c.Request)
		if err != nil {
			log.Println("Error handling certificate issuance request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/view_cert/")

	})

	r.POST("/add_server_cert/", func(c *gin.Context) {
		// Handle the server certificate creation request
		err := CreateServerCert(c)
		if err != nil {
			log.Println("Error creating server certificate:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		log.Println("Server certificate created successfully, redirecting to main page")
		c.Redirect(http.StatusSeeOther, "/view_cert/")
	})

	r.GET("view_cert/", checkSessionMiddleware(), func(c *gin.Context) {
		// Handle the view certificate request
		c.File("templates/frm_view_ca.html")
	})

	r.GET("/view_cert_info/", func(c *gin.Context) {
		// Handle the view certificate request
		err := ViewCertificateForm(c)
		if err != nil {
			log.Println("Error handling view certificate request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		log.Println("Serving view certificate form")
	})

	r.POST("/view_server_key/", func(c *gin.Context) {
		// Handle the request to view the server private key
		serverkey, err := ViewServerPrivateKey(c.Writer, c.Request)
		if err != nil {
			log.Println("Error viewing server private key:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		log.Println("Server private key viewed successfully")
		// return reposnse with server key in json format, using a struct
		serverKeyObj := ServerKey{ServerKey: serverkey}
		c.JSON(http.StatusOK, serverKeyObj)

	})

	r.POST("/revoke_cert/", func(c *gin.Context) {
		// Handle the certificate revocation request
		err := HandleRevokeCertificate(c.Writer, c.Request)
		if err != nil {
			log.Println("Error handling certificate revocation request:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Certificate revoked successfully"})
		log.Println("Certificate revoked successfully, redirecting to main page")
		c.Redirect(http.StatusSeeOther, "/view_cert/")
	})

	r.Run(":8080")

}

/*
	// Create a new http.ServeMux
	mux := http.NewServeMux()
	// Register the logging middleware
	mux.HandleFunc("/", ChainMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Initial settings found, serving main page")
			http.ServeFile(w, r, "index.html")
		},
		Logging(),
		checkSession(),
		checkInitialSettingsMiddleware(),
	))

	// file Server to serve static files
	mux.HandleFunc("/static/", Logging()(func(w http.ResponseWriter, r *http.Request) {

		http.StripPrefix("/static", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
		fmt.Println("Serving static content")

	}))

	// Register the form for initial settings
	mux.Handle("/save_init_cfg/", Logging()(ApplyInitialSettings))

	// Register the form for creating the admin user
	mux.Handle("/add_user/", Logging()(handleAdminUser))

	// REgister the form for login
	//*
		mux.Handle("/login/", Logging()(func(w http.ResponseWriter, r *http.Request) {

			// Serve the login form
			http.ServeFile(w, r, "./static/login.html")

		}))

	//*
	mux.Handle("/login/", (ChainMiddleware(LoginHandler, Logging(), checkInitialSettingsMiddleware())))

	// Register the login handler
	mux.Handle("/authUser/", Logging()(authenticateUser))

	// Register the form for issuing certificates
	mux.Handle("/issue_cert/", (ChainMiddleware(IssueCertHandler, Logging(), checkSession())))

	mux.Handle("/add_server_cert/", (ChainMiddleware(CreateServerCert, Logging(), checkSession())))

	// Register the form for viewing certificates
	mux.Handle("/view_cert/", (ChainMiddleware(ViewCertificateForm, Logging(), checkSession())))

	// Register the handler for viewing server private keys
	mux.Handle("/view_server_key/", (ChainMiddleware(ViewServerPrivateKey, Logging(), checkSession())))

	//Register the handler for revoking certificates
	mux.Handle("/revoke_cert/", (ChainMiddleware(HandleRevokeCertificate, Logging(), checkSession())))

	// Start the server
	//*err := http.ListenAndServe(":8080", mux)
	// Log the error if any
	if err != nil {
		log.Fatal("Failed to start HTTP Server: %v", err)
	}
	//*

	err := http.ListenAndServeTLS(":8443", appCertPath, appKeyPath, mux)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}

	log.Println("Server started on 8443 with TLS enabled")
	log.Println("Visit  https://localhost:8443 to access")
}
*/
