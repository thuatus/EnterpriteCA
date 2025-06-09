package db

// provides database access for the CA application

import (
	"database/sql"
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
	Expire    time.Time
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
