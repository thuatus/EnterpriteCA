package ca

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var now = time.Now()
var strDate = now.Format("020106150405")

var serverName = "test-server-" + strDate + ".com"

func TestGenerateCSR(t *testing.T) {
	//serverName := "test-server.com"
	csrPath, err := GenerateCSR(serverName)
	assert.NoError(t, err, "GenerateCSR should not return an error")
	assert.FileExists(t, csrPath, "CSR file should be created")
	// Clean up
	_ = os.Remove(csrPath)
}

func TestIssueServerCertificate(t *testing.T) {
	//serverName := "test-server.com"
	csrPath, err := GenerateCSR(serverName)
	assert.NoError(t, err, "GenerateCSR should not return an error")
	defer os.Remove(csrPath)

	csrFile, err := os.Open(csrPath)
	assert.NoError(t, err, "Should open CSR file")
	defer csrFile.Close()

	certInfo, err := IssueServerCertificate(serverName, csrFile)
	assert.NoError(t, err, "IssueServerCertificate should not return an error")
	assert.NotEmpty(t, certInfo, "Certificate file should be created")
	// Clean up

}

func TestGetServerCertificateInfo(t *testing.T) {
	//serverName := "test1-server.com"
	_, err := GetServerCertificateInfo(serverName)
	assert.NoError(t, err, "GetServerCertificateInfo should not return an error")
	assert.NotNil(t, "Certificate info should not be nil")
}

func TestRevokeServerCertificate(t *testing.T) {
	//serverName := "test-server.com"
	err := RevokeServerCertificate(serverName)
	assert.NoError(t, err, "RevokeServerCertificate should not return an error")
}
