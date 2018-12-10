package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"text/template"
)

func GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}

func GenerateHTML(w http.ResponseWriter, name string, data interface{}) {
	templates := template.Must(template.ParseGlob("templates/*"))
	templates.ExecuteTemplate(w, name, data)
}