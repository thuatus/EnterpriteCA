// Implements ca basics apis
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Creates folders and files structure
func CreateCA(caName string) (caPathName string) {

	intermediateCA := []string{"intermediate-", caName}
	caPath := []string{"/home/alvaro/srv/CA/", caName}
	intermediateCAPath := strings.Join(caPath, "") + "/" + strings.Join(intermediateCA, "")

	err := os.MkdirAll(caPath[0], 06000)
	if err != nil {
		log.Fatalf("can't create ca folder %v: %d", caPath, err)
	}

	err = os.MkdirAll(intermediateCAPath, 07770)
	if err != nil {
		log.Fatalf("can't create intermediate ca folder %v: %d", caPath, err)
	}

	caPathName = fmt.Sprintf("CA path: %s, Intermediate CA path: %s", caPath, intermediateCAPath)

	return caPathName
}

// main to validade
func main() {
	caName := "EMPRESA"
	caPathName := CreateCA(caName)
	log.Printf("CA path: %s", caPathName)
}
