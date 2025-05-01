// Implements ca basics apis
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

// Creates folders and files structure
func createPKIFolders(caName string) (string, string, error) {

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
func createConfigFiles(caPath string, country string, state string, local string, organization string, organizationalUnit string) (string, error) {

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
		log.Fatalf("Error writing serial file: %v", err)
		return "", err
	}
	defer serialFile.Close()

	//create cerl files
	crlFile, err := os.Create(caPath + "/crlnumber")
	if err != nil {
		log.Fatalf("Error creating crl file: %v", err)
		return "", err
	}

	if _, err := crlFile.WriteString("100"); err != nil {
		log.Fatalf("Error writing serial file: %v", err)
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

	//create config file
	// Parse the templateCA file
	tmplCAConfig, err := template.ParseFiles("templates/ca-cnf.txt")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
		return "", err
	}

	// Parse the template IntermediateCA file
	tmplIntCAConfig, err := template.ParseFiles("templates/intermediateCA-cnf.txt")
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

	return fileCAConfig.Name() + "\n" + fileIntConfig.Name(), nil
}

// Creates CA certificate, private key and public key
func createCACertificate(caPathName string) (string, error) {
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
	fmt.Printf("--->%v", cmd)
	intermediateCACertinfo, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error creating intermediate CA certificate %v: %d", intermediateCACertinfo, err)
		return "", err
	}

	return string(caKeyinfo) + "|" + string(caCertinfo), nil
}

// main to validade
func main() {
	caName := "CA-10"
	fmt.Printf("Creating CA ...\n")
	caPathName, intermediateCaPathName, err := createPKIFolders(caName)
	if err != nil {
		log.Fatalf("Error creating CA Folders: %v", err)
	}
	log.Printf("Folder structure created!\nCA path: %s \n IntermediateCAPAth: %s", caPathName, intermediateCaPathName)

	fmt.Printf("Creating CA config files ...\n")
	country := "BR"
	state := "Distrito Federal"
	local := "Brasilia"
	organization := "CA-10"
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
	log.Printf("CA certificate created! %s", cacertificate)
}
