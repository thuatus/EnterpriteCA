// Implements ca basics apis
package main

import (
	"fmt"
	"log"
	"os"
)

// Creates folders and files structure
func CreatePKIFolders(caName string) (string, error) {

	var caPathName string
	rootPath := "/home/alvaro/srv/CA/"
	intermediateCaName := "intermediate-" + caName
	caPath := rootPath + caName
	intermediateCAPath := caPath + "/" + intermediateCaName
	fmt.Printf("folder path: %v\n", caPath)
	infoDir, err := os.Stat(caPath)

	if err == nil {

		if infoDir.IsDir() {
			log.Fatalf("CA folder already exists %v: %d", caPath, err)
			return "", err
		}

	} else {

		fmt.Printf("Creating CA Folder %s ...\n", caPath)

		err = os.MkdirAll(caPath, 0700)
		if err != nil {
			log.Fatalf("Can't create ca folder %v: %d", caPath, err)
			return "", err
		}
		err = os.MkdirAll(intermediateCAPath, 0700)
		if err != nil {
			log.Fatalf("Can't intermediate ca folder %v: %d", intermediateCAPath, err)
			return "", err
		}

	}
	caPathName = fmt.Sprintf("CA path: %s | Intermediate CA path: %s", caPath, intermediateCAPath)
	return caPathName, nil
}

// main to validade
func main() {
	caName := "CA-6"
	caPathName, err := CreatePKIFolders(caName)
	if err != nil {
		log.Fatalf("Error creating CA: %v", err)
	}
	log.Printf("CA path: %s", caPathName)
}
