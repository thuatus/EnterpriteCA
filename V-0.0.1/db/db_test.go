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
