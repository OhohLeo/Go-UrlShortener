package main

import (
	// "fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

const (
	HTTP = "http://127.0.0.1"
	PORT = ":8080"
)

func TestUrlShortener(t *testing.T) {

	go func() {
		assert.Nil(t, new(UrlShortener).Start(PORT))
	}()

	// Requêtes invalides
	urlBadRequest(t, "/encode", "", "Invalid url parse : empty url\n")
	urlBadRequest(t, "/encode", "www.test.bzh",
		"Invalid url parse www.test.bzh: invalid URI for request\n")

	// Encodage
	urlOk(t, "/encode", "http://www.test.bzh",
		"http://127.0.0.1:8080/1B2M2Y8AsgTpgAmY7PhCfg==", http.StatusCreated)

	// Encodage d'une adresse similaire
	urlOk(t, "/encode", "http://www.test.bzh",
		"http://127.0.0.1:8080/1B2M2Y8AsgTpgAmY7PhCfg==", http.StatusNotModified)

	// Décodage
	urlOk(t, "/decode", "http://127.0.0.1:8080/1B2M2Y8AsgTpgAmY7PhCfg==",
		"http://www.test.bzh", http.StatusOK)
	urlOk(t, "/decode", "http://127.0.0.1:8080/1B2M2Y8AsgTpgAmY7PhCfg==",
		"http://www.test.bzh", http.StatusOK)
}

func urlOk(t *testing.T, url string, input string, output string, status int) {

	rsp, err := http.Post(HTTP+PORT+url, "", strings.NewReader(input))
	if assert.Nil(t, err) == false {
		return
	}

	assert.Equal(t, status, rsp.StatusCode)
	assert.Equal(t, output, rsp.Header.Get("Location"))

}

func urlBadRequest(t *testing.T, url string, location string, error string) {

	rsp, err := http.Post(HTTP+PORT+url, "", strings.NewReader(location))
	if assert.Nil(t, err) == false {
		return
	}

	assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	assert.Equal(t, "", rsp.Header.Get("Location"))

	// Récupération de l'erreur
	body, err := ioutil.ReadAll(rsp.Body)
	if assert.Nil(t, err) == false {
		return
	}

	// Validation de l'url
	assert.Equal(t, error, string(body))
}
