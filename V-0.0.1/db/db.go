// Package db provides database access and management functions for the Certificate Authority (CA) application.
// It includes functionality to connect to a MySQL database, add and retrieve certificate information, and update
// certificate validity status. The package defines the DBCertInfo type to represent certificate records and
// handles conversion between database representations and application data structures.
package db

// provides database access for the CA application

import (
	"database/sql"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

// dbCertInfo is a type that represents the certificate information stored in the database
type DBCertInfo struct {
	Issuer    string
	Subject   string
	PublicKey string
	Created   time.Time
	Expire    time.Time
	Status    int // 1 for valid, 0 for invalid
}

// ConnectDB connects to the CA database and returns a pointer to the sql.DB instance
func ConnectDB() (*sql.DB, error) {

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("MYSQL_USER")
	cfg.Passwd = os.Getenv("MYSQL_PASSWORD")
	if cfg.Passwd == "" {
		return nil, fmt.Errorf("MYSQL_PASSWORD environment variable is not set:%s", cfg.Passwd)
	}

	cfg.Net = "tcp"
	cfg.Addr = "localhost:3306"
	cfg.DBName = "ca"

	//Cria o tratamento de erro de conexão do DB
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	return db, nil

}

// Add user to the database
func AddUser(username, password string) error {
	log.Println("Conecting to the database...")
	dbConn, err := ConnectDB()

	if err != nil {
		log.Println("Error connecting to the database:", err)
		return err
	}
	defer dbConn.Close()

	//  create the admin user
	_, err = dbConn.Exec("INSERT INTO ca.users (name, passwd, active ) VALUES (?, ?, ?)", username, password, 1)
	if err != nil {
		log.Println("Error inserting admin user into the database:", err)
		return err
	}

	return nil
}

// Add certificate information to the database
func AddCertificate(cert DBCertInfo) error {
	// AddCertificate adds a certificate to the database
	if db == nil {
		return fmt.Errorf("database connection is not established")
	}

	// Prepare the SQL query to insert the certificate
	query := "INSERT INTO ca (issuer,subject,public_key,created,expiration,is_valid) VALUES (?,?,?,?,?,?)"
	_, err := db.Exec(query, cert.Issuer, cert.Subject, cert.PublicKey, time.Now(), cert.Expire, 1)
	if err != nil {
		return fmt.Errorf("failed to add certificate: %w", err)
	}
	return nil

}

// Retrieve Certificates from the database by its subject
func GetServerCertificates() ([]DBCertInfo, error) {

	if db == nil {
		return nil, fmt.Errorf("database connection is not established")
	}
	query := "SELECT issuer, subject, public_key, created, expiration, is_valid FROM ca where is_valid = 1"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve certificates: %v", err)
	}
	defer rows.Close()

	var certs []DBCertInfo
	for rows.Next() {
		var cert DBCertInfo
		var date_create string
		var date_exp string
		if err := rows.Scan(&cert.Issuer, &cert.Subject, &cert.PublicKey, &date_create, &date_exp, &cert.Status); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		cert.Created, err = time.Parse("2006-01-02 15:04:05", date_create)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %v", err)
		}

		cert.Expire, err = time.Parse("2006-01-02 15:04:05", date_exp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expiration date: %v", err)
		}

		// convert public key from hex to pem data
		pubKeyBytes, err := hex.DecodeString(cert.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode public key: %v", err)
		}

		pemBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: []byte(pubKeyBytes),
		}
		pemData := pem.EncodeToMemory(pemBlock)
		// append the certificate to the slice
		cert.PublicKey = string(pemData)

		certs = append(certs, cert)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %v", err)
	}
	return certs, nil

}

// Alter information of a certificate to Rovoked in the database
func UpdateValidCertificate(subject string) error {

	db, err := ConnectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Prepare the SQL query to update the certificate status
	query := "UPDATE ca SET is_valid = 0 WHERE subject = ?"
	_, err = db.Exec(query, subject)
	if err != nil {
		return fmt.Errorf("failed to update revoke certificate information: %w", err)
	}
	return nil
}
