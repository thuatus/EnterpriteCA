package db

// provides database access for the CA application

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

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

/*func AddCertificate(cert string) error {
	// AddCertificate adds a certificate to the database
	query := "INSERT INTO certificates (cert) VALUES (?)"
	_, err := db.Exec(query, cert)
	if err != nil {
		return fmt.Errorf("failed to add certificate: %w", err)
	}
	return nil

}*/
