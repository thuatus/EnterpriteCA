// Implements ca basics apis
package ca

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

// struct to store Certificate Information
type CertInfo struct {
	Issuer             string
	Subject            string
	PublicKey          []byte
	Country            []string
	State              []string
	Locality           []string
	Organization       []string
	OrganizationalUnit []string
	Expire             string
}

// Creates folders and files structure
func CreatePKIFolders(caName string) (string, string, error) {

	rootPath := "/home/alvaro/srv/CA/"
	intermediateCaName := "intermediateCA"
	caPath := rootPath + caName
	intermediateCAPath := caPath + "/" + intermediateCaName
	certsFolderPath := intermediateCAPath + "/certs"
	csrFOlderPath := intermediateCAPath + "/csr"
	crlFolderPath := intermediateCAPath + "/crl"
	fmt.Printf("folder path: %v\n", caPath)
	infoDir, err := os.Stat(caPath)

	if err == nil {

		if infoDir.IsDir() {
			log.Fatalf("CA folder already exists %v: %d", caPath, err)
			return "", "", err
		}

	} else {

		fmt.Printf("Creating CA Folder %s ...\n", caPath)

		err = os.MkdirAll(caPath, 0700)
		if err != nil {
			log.Fatalf("Can't create ca folder %v: %d", caPath, err)
			return "", "", err
		}
		err = os.MkdirAll(intermediateCAPath, 0700)
		if err != nil {
			log.Fatalf("Can't intermediate ca folder %v: %d", intermediateCAPath, err)
			return "", "", err
		}
		err = os.MkdirAll(intermediateCAPath, 0700)
		if err != nil {
			log.Fatalf("Can't intermediate ca folder %v: %d", intermediateCAPath, err)
			return "", "", err
		}
		err = os.MkdirAll(certsFolderPath, 0700)
		if err != nil {
			log.Fatalf("Can't certs ca folder %v: %d", certsFolderPath, err)
			return "", "", err
		}
		err = os.MkdirAll(csrFOlderPath, 0700)
		if err != nil {
			log.Fatalf("Can't crs ca folder %v: %d", csrFOlderPath, err)
			return "", "", err
		}
		err = os.MkdirAll(crlFolderPath, 0700)
		if err != nil {
			log.Fatalf("Can't crl folder %v: %d", crlFolderPath, err)
			return "", "", err
		}

	}
	//caPathName = fmt.Sprintf("CA path: %s | Intermediate CA path: %s", caPath, intermediateCAPath)
	return caPath, intermediateCAPath, nil
}

// create cnf files to PKI
func CreateConfigFiles(caPath string, country string, state string, local string, organization string, organizationalUnit string) (string, error) {

	// structure that store cnf config parameters
	type config struct {
		C     string
		St    string
		L     string
		O     string
		Ou    string
		CADir string
	}

	//create serial files
	serialFile, err := os.Create(caPath + "/serial")
	if err != nil {
		log.Fatalf("Error creating serial file: %v", err)
		return "", err
	}

	if _, err := serialFile.WriteString("1000"); err != nil {
		log.Fatalf("Error writing serialCA file: %v", err)
		return "", err
	}
	defer serialFile.Close()

	//create crl files
	crlFile, err := os.Create(caPath + "/crlnumber")
	if err != nil {
		log.Fatalf("Error creating crlCA file: %v", err)
		return "", err
	}

	if _, err := crlFile.WriteString("100"); err != nil {
		log.Fatalf("Error writing crlCA file: %v", err)
		return "", err
	}
	defer crlFile.Close()

	//create index file
	indexFile, err := os.Create(caPath + "/index.txt")
	if err != nil {
		log.Fatalf("Error creating index file: %v", err)
		return "", err
	}
	defer indexFile.Close()

	//create index Intermediate CA file

	serialIntFile, err := os.Create(caPath + "/intermediateCA/serial")
	if err != nil {
		log.Fatalf("Error creating serial file: %v", err)
		return "", err
	}

	if _, err := serialIntFile.WriteString("1000"); err != nil {
		log.Fatalf("Error writing serialCA file: %v", err)
		return "", err
	}
	defer serialIntFile.Close()

	//create crl files
	crlIntCAFile, err := os.Create(caPath + "/intermediateCA/crlnumber")
	if err != nil {
		log.Fatalf("Error creating crlCA file: %v", err)
		return "", err
	}

	if _, err := crlIntCAFile.WriteString("100"); err != nil {
		log.Fatalf("Error writing crlCA file: %v", err)
		return "", err
	}
	defer crlIntCAFile.Close()

	//create index file
	indexintCAFile, err := os.Create(caPath + "/intermediateCA/index.txt")
	if err != nil {
		log.Fatalf("Error creating index file: %v", err)
		return "", err
	}
	defer indexintCAFile.Close()

	//create config file
	// Parse the templateCA file
	wd, _ := os.Getwd()
	log.Printf("Current working directory: %s\n", wd)

	tmplCAConfig, err := template.ParseFiles("../ca/templates/ca-cnf.txt")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}

	// Parse the template IntermediateCA file
	tmplIntCAConfig, err := template.ParseFiles("../ca/templates/intermediateCA-cnf.txt")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}
	caconfig := config{
		C:     country,
		St:    state,
		L:     local,
		O:     organization,
		Ou:    organizationalUnit,
		CADir: caPath,
	}
	fileCAConfig, err := os.Create(caPath + "/ca.cnf")
	if err != nil {
		log.Fatalf("Error creating config CA file: %v", err)
		return "", err
	}
	defer fileCAConfig.Close()

	err = tmplCAConfig.Execute(fileCAConfig, caconfig)
	if err != nil {
		log.Fatalf("Error configuring CA template file: %v", err)
		return "", err
	}

	//create config file for intermediate CA
	fileIntConfig, err := os.Create(caPath + "/intermediateCA/intermediateCA.cnf")
	if err != nil {
		log.Fatalf("Error creating config IntermediateCA file: %v", err)
		return "", err
	}
	defer fileIntConfig.Close()

	err = tmplIntCAConfig.Execute(fileIntConfig, caconfig)
	if err != nil {
		log.Fatalf("Error configuring Intermediate CA template file: %v", err)
		return "", err
	}

	defFilePath := "/home/alvaro/srv/CA/definitions.txt"
	defFile, err := os.Create(defFilePath)
	if err != nil {
		log.Fatalf("Error creating definition file: %v", err)
		return "", err
	}

	if _, err := defFile.WriteString(caPath + "|" + caPath + "/intermediateCA"); err != nil {
		log.Fatalf("Error writing definition folder file: %v", err)
		return "", err
	}
	defer defFile.Close()
	return fileCAConfig.Name() + "\n" + fileIntConfig.Name(), nil
}

// get the CA folder info
func getCAInfo() (string, string, error) {

	// Lê o conteúdo do arquivo
	caInfo, err := os.ReadFile("/home/alvaro/srv/CA/definitions.txt")
	if err != nil {
		fmt.Printf("can't read file %v\n", err)
		return "", "", err
	}

	str := string(caInfo)

	paths := strings.Split(str, "|")

	return paths[0], paths[1], nil
}

// returns the Informnation about the Intermediate CA
func GetIntermediateCAInfo() (CertInfo, error) {
	// Lê o conteúdo do arquivo
	IntCAInfo, err := os.ReadFile("/home/alvaro/srv/CA/definitions.txt")
	if err != nil {
		fmt.Printf("can't read file %v\n", err)
		return CertInfo{}, fmt.Errorf("can not read the intca file: %v", err)
	}

	str := string(IntCAInfo)

	paths := strings.Split(str, "|")

	if len(paths) < 2 {
		return CertInfo{}, fmt.Errorf("invalid CA definitions file format")
	}

	certPath := paths[1] + "/intermediateCA.crt"
	certFile, err := os.ReadFile(certPath)
	if err != nil {
		fmt.Println("error reading intermediate CA certificate file:", err)
		return CertInfo{}, fmt.Errorf("error reading intermediate CA certificate file: %v", err)
	}

	// Decodificando o certificado PEM
	block, _ := pem.Decode(certFile)
	if block == nil {
		fmt.Println("error decoding intermediate CA certificate PEM block")
		return CertInfo{}, fmt.Errorf("error decoding intermediate CA certificate PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Erro ao parsear o certificado:", err)
		return CertInfo{}, fmt.Errorf("error parsing intermediate CA certificate: %v", err)
	}

	certInfo := CertInfo{
		Issuer:             cert.Issuer.String(),
		Subject:            cert.Subject.String(),
		Country:            cert.Issuer.Country,
		State:              cert.Issuer.Province,
		Locality:           cert.Issuer.Locality,
		Organization:       cert.Issuer.Organization,
		OrganizationalUnit: cert.Issuer.OrganizationalUnit,
	}

	return certInfo, nil
}

// Creates CA certificate, private key and public key
func CreateCACertificate(caPathName string) (string, error) {
	caCertificate := caPathName + "/ca.crt"
	caPrivateKey := caPathName + "/ca.key"
	caConfigFIle := caPathName + "/ca.cnf"

	cmd := exec.Command("openssl", "genrsa", "-out", caPrivateKey, "4096")

	// Captura a saída do comando
	caKeyinfo, err := cmd.Output()
	if err != nil {
		fmt.Println("Error creating CA key:", err)
		return "", err
	}

	cmd = exec.Command("openssl", "req", "-new", "-x509", "-days", "3650", "-key", caPrivateKey, "-out", caCertificate, "-config", caConfigFIle)
	caCertinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating CA Certificate %v: %d", caCertinfo, err)
		return "", err
	}

	// create intermediate certificate based on intermediateCA.cnf
	intermediateCACertificate := caPathName + "/intermediateCA/intermediateCA.crt"
	intermediateCAKey := caPathName + "/intermediateCA/intermediateCA.key"
	intermediateCARequest := caPathName + "/intermediateCA/intermediateCA.csr"

	cmd = exec.Command("openssl", "genrsa", "-out", intermediateCAKey, "4096")
	intermediateCAKeyinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating intermediate CA key %v: %d", intermediateCAKeyinfo, err)
		return "", err
	}

	//create intermediate certificate request
	cmd = exec.Command("openssl", "req", "-key", intermediateCAKey, "-new", "-sha256", "-out", intermediateCARequest, "-config", caConfigFIle, "-batch")

	intermediateCAReqinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating intermediate CA request %v: %d", intermediateCAReqinfo, err)
		return "", err
	}

	// create intermediate certificate based on intermediateCA.cnf
	cmd = exec.Command("openssl", "ca", "-config", caConfigFIle, "-extensions", "v3_intermediate_ca", "-notext", "-in", intermediateCARequest, "-md", "sha256", "-out", intermediateCACertificate, "-days", "1825", "-batch")

	intermediateCACertinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating intermediate CA certificate %v: %d", intermediateCACertinfo, err)
		return "", err
	}

	return string(caKeyinfo) + "|" + string(caCertinfo), nil
}

// Generates a private key for a server certificate file
func GenerateServerKey(serverName string) (string, error) {

	_, keyPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA structure definitions: %v", err)
		return "", err
	}
	serverKeyPath := keyPath + "/private/" + serverName + ".key"
	if _, err := os.Stat(serverKeyPath); err == nil {
		log.Printf("Server key %s already exists, skipping generation.\n", serverKeyPath)
		return serverKeyPath, nil
	}
	fmt.Printf("Generating server private key %s ...\n", serverName+".key")
	// Execute the openssl command to generate the server key
	cmd := exec.Command("openssl", "genrsa", "-out", serverKeyPath, "4096")
	serverKeyinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating server key %v: %d", serverKeyinfo, err)
		return "", err
	}
	fmt.Printf("Server private key generated: %s\n", serverKeyPath)
	return serverKeyPath, nil
}

// Generates a Certificate Signing Request (CSR) for a server certificate file
func GenerateCSR(serverName string) (string, error) {

	serverCertinfo, err := GetIntermediateCAInfo()
	if err != nil {
		log.Fatalf("Error getting Intermediate CA information: %v", err)
		return "", err
	}

	// Generate the server key if it does not exist
	serverKeyPath, err := GenerateServerKey(serverName)
	if err != nil {
		log.Fatalf("Error generating server key: %v", serverKeyPath)
		return "", err
	}
	fmt.Printf("Server private key generated: %s\n", serverKeyPath)

	// Use the Intermediate CA path to store the CSR
	_, intCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
		return "", err
	}

	serverCsrPath := intCAPath + "/csr/" + serverName + ".csr"
	if _, err := os.Stat(serverCsrPath); err == nil {
		log.Printf("Server CSR %s already exists, skipping generation.\n", serverCsrPath)
		return "", nil
	}

	//fmt.Printf("Generating server certificate request %s ...\n", serverName+".csr")

	subject := fmt.Sprintf("/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s", serverCertinfo.Country[0], serverCertinfo.State[0], serverCertinfo.Locality[0], serverCertinfo.Organization[0], serverCertinfo.OrganizationalUnit[0], serverName)

	fmt.Println("Subject for server certificate request:", subject)
	// Execute the openssl command to generate the server certificate request

	// generating server certificate request

	cmd := exec.Command("openssl", "req", "-new", "-key", serverKeyPath, "-out", serverCsrPath, "-subj", subject, "-addext", fmt.Sprintf("subjectAltName = DNS:%s", serverName))
	serverCsrinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating server certificate request %v: %d", serverCsrinfo, err)
		return "", err
	}

	fmt.Printf("Server certificate request generated: %s\n", serverCsrPath)
	return serverCsrPath, nil

}

// Issue a server certificate
func IssueServerCertificate(ServerName string, csr io.Reader) (string, error) {

	fmt.Printf("getting CA information ...\n")
	_, fintCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
	}

	serverCrt := fintCAPath + "/certs/" + ServerName + ".crt"
	ServerReq := fintCAPath + "/csr/" + ServerName + ".csr"
	csrInfo, err := io.ReadAll(csr)
	if err != nil {
		fmt.Println("Can't read csr file:", err)
		return "", err
	}
	// Write the contents to a file
	csrFile, err := os.Create(ServerReq)
	if err != nil {
		fmt.Println("Error creating CSR file:", err)
		return "", err
	}
	defer csrFile.Close()
	_, err = csrFile.Write(csrInfo)
	if err != nil {
		fmt.Println("Error writing CSR file:", err)
		return "", err
	}

	// create server certificate based on intermediateCA.cnf

	cmd := exec.Command("openssl", "ca", "-config", fintCAPath+"/intermediateCA.cnf", "-extensions", "server_cert", "-notext", "-md", "sha256", "-in", ServerReq, "-out", serverCrt, "-days", "365", "-batch")

	serverCertinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating server certificate: %v, %v", serverCertinfo, err)
		return "", err
	}

	// show the server certificate

	cmd = exec.Command("openssl", "x509", "-in", serverCrt, "-text", "-noout")

	serverCertinfo, err = cmd.Output()
	if err != nil {
		log.Fatalf("Error showing server certificate %v: %d", serverCertinfo, err)
		return "", err
	}
	return string(serverCertinfo), nil
}

// return server certificate information
func GetServerCertificateInfo(serverName string) (CertInfo, error) {
	_, fintCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
		return CertInfo{}, err
	}

	serverCertPath := fintCAPath + "/certs/" + serverName + ".crt"
	certFile, err := os.ReadFile(serverCertPath)
	if err != nil {
		fmt.Println("error reading server certificate file:", err)
		return CertInfo{}, fmt.Errorf("error reading server certificate file: %v", err)
	}

	// Decodificando o certificado PEM
	block, _ := pem.Decode(certFile)
	if block == nil {
		fmt.Println("error decoding server certificate PEM block")
		return CertInfo{}, fmt.Errorf("error decoding server certificate PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing the server certificate:", err)
		return CertInfo{}, fmt.Errorf("error parsing server certificate: %v", err)
	}

	dt_expire := cert.NotAfter.Format("2006-01-02 15:04:05")
	certInfo := CertInfo{
		Issuer:             cert.Issuer.String(),
		Subject:            cert.Subject.String(),
		PublicKey:          cert.Raw,
		Country:            cert.Issuer.Country,
		State:              cert.Issuer.Province,
		Locality:           cert.Issuer.Locality,
		Organization:       cert.Issuer.Organization,
		OrganizationalUnit: cert.Issuer.OrganizationalUnit,
		Expire:             dt_expire,
	}

	return certInfo, nil
}

// Retirn a serverNAme of a certificate given the certificate subject
func GetServerNameFromCert(certSubject string) (string, error) {

	re := regexp.MustCompile(`CN=([^,]+)`)
	servername := re.FindStringSubmatch(certSubject)
	if len(servername) < 2 {
		return "", fmt.Errorf("server name not found in certificate subject")
	}
	return servername[1], nil
}

// Return Server Certificate Private Key
func GetServerPrivateKey(serverName string) (string, error) {
	_, fintCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
		return "", err
	}

	serverKeyPath := fintCAPath + "/private/" + serverName + ".key"
	if _, err := os.Stat(serverKeyPath); os.IsNotExist(err) {
		log.Printf("Server private key %s does not exist.\n", serverKeyPath)
		return "", fmt.Errorf("server private key does not exist")
	}

	// set the content of the server key file on keyContent
	keyContent, err := os.ReadFile(serverKeyPath)
	if err != nil {
		log.Printf("Error reading server private key file %s: %v\n", serverKeyPath, err)
		return "", fmt.Errorf("error reading server private key file: %v", err)
	}

	// If the key is in PEM format, return the path
	return string(keyContent), nil
}

func RevokeServerCertificate(serverName string) error {
	_, fintCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
		return err
	}

	serverCertPath := fintCAPath + "/certs/" + serverName + ".crt"
	if _, err := os.Stat(serverCertPath); os.IsNotExist(err) {
		log.Printf("Server certificate %s does not exist.\n", serverCertPath)
		return fmt.Errorf("server certificate does not exist")
	}

	cmd := exec.Command("openssl", "ca", "-config", fintCAPath+"/intermediateCA.cnf", "-revoke", serverCertPath, "-batch")
	revokeInfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error revoking server certificate: %v, %v", revokeInfo, err)
		return err
	}

	fmt.Printf("Server certificate revoked: %s\n", serverCertPath)

	// Generate CRL after revoking the certificate
	crlInfo, err := GenerateCRL()
	if err != nil {
		log.Fatalf("Error generating CRL after revoking certificate: %v", err)
		return err
	}
	fmt.Printf("CRL generated after revoking certificate: %s\n", crlInfo)

	return nil
}

// Generate and update a Certificate Revocation List (CRL)
func GenerateCRL() (string, error) {
	_, fintCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
		return "", err
	}

	crlFilePath := fintCAPath + "/crl/crl.pem"
	cmd := exec.Command("openssl", "ca", "-config", fintCAPath+"/intermediateCA.cnf", "-gencrl", "-out", crlFilePath, "-batch")

	crlInfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error generating CRL: %v, %v", crlInfo, err)
		return "", err
	}

	fmt.Printf("CRL generated: %s\n", crlFilePath)
	return string(crlInfo), nil
}

// GetCRLInfo returns the content of the Certificate Revocation List (CRL)
func GetCRLInfo() (string, error) {
	_, fintCAPath, err := getCAInfo()
	if err != nil {
		log.Fatalf("Error get CA definitions: %v", err)
		return "", err
	}

	crlFilePath := fintCAPath + "/crl/crl.pem"
	if _, err := os.Stat(crlFilePath); os.IsNotExist(err) {
		log.Printf("CRL file %s does not exist.\n", crlFilePath)
		return "", fmt.Errorf("CRL file does not exist")
	}

	// set the content of the CRL file on crlContent
	crlContent, err := os.ReadFile(crlFilePath)
	if err != nil {
		log.Printf("Error reading CRL file %s: %v\n", crlFilePath, err)
		return "", fmt.Errorf("error reading CRL file: %v", err)
	}

	return string(crlContent), nil
}

/*
// main to validade
func main() {
	caName := "CA-01"

	fmt.Printf("Creating CA ...\n")
	caPathName, intermediateCaPathName, err := createPKIFolders(caName)
	if err != nil {
		log.Fatalf("Error creating CA Folders: %v", err)
	}
	log.Printf("Folder structure created!\nCA path: %s \n IntermediateCAPAth: %s\n", caPathName, intermediateCaPathName)

	fmt.Printf("Creating CA config files ...\n")
	country := "BR"
	state := "Distrito Federal"
	local := "Brasilia"
	organization := "CA-01"
	organizationalUnit := "Authority"

	createConfigFiles(caPathName, country, state, local, organization, organizationalUnit)
	if err != nil {
		log.Fatalf("Error creating CA config files: %v", err)
	}
	log.Printf("CA config files created!\n")

	fmt.Printf("Creating CA certificate ...\n")
	cacertificate, err := createCACertificate(caPathName)

	if err != nil {
		log.Fatalf("Error creating CA certificate: %v", err)
	}
	log.Printf("CA certificate created! %s\n", cacertificate)

	fmt.Printf("Creating Server certificate ...\n")
	reqFile, err := os.Open("/home/alvaro/app01.csr")
	if err != nil {
		log.Fatalf("Error opening CSR file: %v", err)
	}
	defer reqFile.Close()
	serverCert, err := issueServerCertificate("app01", reqFile)
	if err != nil {
		log.Fatalf("Error creating server certificate: %v", err)
	}
	log.Printf("Server certificate created! %s\n", serverCert)

}
*/
