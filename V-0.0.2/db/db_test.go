package db

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Example DBCertInfo for testing
var dataPubKey = []byte("TestPublicKey")
var hexString = hex.EncodeToString(dataPubKey)
var testCert = DBCertInfo{
	Issuer:    "Test Issuer",
	Subject:   "Test Subject",
	PublicKey: hexString,
	Created:   time.Now(),
	Expire:    time.Now().Add(365 * 24 * time.Hour),
}

func TestConnectDB(t *testing.T) {
	// Connect to the database
	db, err := ConnectDB()
	assert.NoError(t, err, "ConnectDB should not return an error")
	assert.NotNil(t, db, "Database connection should not be nil")

	// Close the database connection
	err = db.Close()
	assert.NoError(t, err, "Closing the database connection should not return an error")
}
func TestAddAndGetCertificate(t *testing.T) {
	// Add certificate
	db, _ := ConnectDB()
	assert.NotNil(t, db, "Database connection should not be nil")
	defer db.Close()

	err := AddCertificate(testCert)
	assert.NoError(t, err, "AddCertificate should not return an error")

	// Retrieve certificate by subject
	certs, err := GetServerCertificates()
	assert.NoError(t, err, "GetCertificatesBySubject should not return an error")
	assert.NotEmpty(t, certs, "Should retrieve at least one certificate")

	// Check if the retrieved certificate matches the inserted one
	found := false
	for _, cert := range certs {
		if cert.Subject == testCert.Subject && cert.Issuer == testCert.Issuer {
			found = true
			break
		}
	}
	assert.True(t, found, "Inserted certificate should be found in the database")
}

func TestRevokeCertificate(t *testing.T) {
	// Revoke the test certificate
	err := UpdateValidCertificate(testCert.Subject)
	assert.NoError(t, err, "RevokeCertificate should not return an error")
}

// teste add user
func TestAddUser(t *testing.T) {
	// Connect to the database
	db, err := ConnectDB()
	assert.NoError(t, err, "ConnectDB should not return an error")
	assert.NotNil(t, db, "Database connection should not be nil")
	defer db.Close()

	// Add a user
	err = AddUser("testuser", "testpasswordtyrytry")
	assert.NoError(t, err, "AddUser should not return an error")

}
