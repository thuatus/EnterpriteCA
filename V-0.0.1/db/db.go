package db

// provides database access for the CA application

import (
	"database/sql"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
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
	cfg.User = "root"
	cfg.Passwd = "passwd@G0"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
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

// add certificate information to the database
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

// Retrieves Certificates from the database by its subject

func GetServerCertificates() ([]DBCertInfo, error) {

	if db == nil {
		return nil, fmt.Errorf("database connection is not established")
	}
	query := "SELECT issuer, subject, public_key, created, expiration, is_valid FROM ca"
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

	if db == nil {
		return fmt.Errorf("database connection is not established")
	}

	// Prepare the SQL query to update the certificate status
	query := "UPDATE ca SET is_valid = 0 WHERE subject = ?"
	_, err := db.Exec(query, subject)
	if err != nil {
		return fmt.Errorf("failed to update revoke certificate information: %w", err)
	}
	return nil
}
