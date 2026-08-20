//go:build ignore

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
)

func main() {
	keyName := "APP_KEY"
	if len(os.Args) > 1 {
		keyName = os.Args[1]
	}

	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Println("Error: .env not found. Create it first.")
		os.Exit(1)
	}

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		fmt.Println("Error generating secret:", err)
		os.Exit(1)
	}
	secret := base64.StdEncoding.EncodeToString(secretBytes)

	content, err := os.ReadFile(envPath)
	if err != nil {
		fmt.Println("Error reading .env:", err)
		os.Exit(1)
	}

	pattern := regexp.MustCompile(`(?m)^` + keyName + `=.*$`)
	newLine := keyName + "=" + secret

	var updated []byte
	if pattern.Match(content) {
		updated = pattern.ReplaceAll(content, []byte(newLine))
		fmt.Printf("Updated %s in .env\n", keyName)
	} else {
		updated = append(content, []byte("\n"+newLine+"\n")...)
		fmt.Printf("Added %s to .env\n", keyName)
	}

	if err := os.WriteFile(envPath, updated, 0644); err != nil {
		fmt.Println("Error writing .env:", err)
		os.Exit(1)
	}
}
