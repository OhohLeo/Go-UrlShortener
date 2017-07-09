package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUrlShortener(t *testing.T) {

	// Aucun serveur
	_, err := http.Post("/encode", "", strings.NewReader(""))
	if assert.NotNil(t, err) == false {
		t.Fail()
	}

	server := httptest.NewServer(new(UrlShortener).Init())

	// Requêtes invalides
	for _, route := range []string{"encode", "decode", "redirect"} {
		urlError(t, http.StatusBadRequest,
			server.URL+"/"+route, "",
			"Invalid url parse : empty url\n")
		urlError(t, http.StatusBadRequest,
			server.URL+"/"+route, "www.test.bzh",
			"Invalid url parse www.test.bzh: invalid URI for request\n")
	}

	// Encodage
	shortUrl := urlOk(t, server.URL+"/encode",
		"http://www.test.bzh", http.StatusCreated)

	// Décodage
	longUrl := urlOk(t, server.URL+"/decode",
		shortUrl, http.StatusOK)

	assert.Equal(t, "http://www.test.bzh", longUrl)

	server.Close()
}

func urlOk(t *testing.T, url string, input string, status int) (location string) {

	rsp, err := http.Post(url, "", strings.NewReader(input))
	if assert.Nil(t, err) == false {
		return
	}

	assert.Equal(t, status, rsp.StatusCode)

	location = rsp.Header.Get("Location")
	assert.NotNil(t, location)
	return
}

func urlError(t *testing.T, status int, url string, location string, error string) {

	rsp, err := http.Post(url, "", strings.NewReader(location))
	if assert.Nil(t, err) == false {
		return
	}

	assert.Equal(t, status, rsp.StatusCode)
	assert.Equal(t, "", rsp.Header.Get("Location"))

	// Récupération de l'erreur
	body, err := ioutil.ReadAll(rsp.Body)
	if assert.Nil(t, err) == false {
		return
	}

	// Validation de l'url
	assert.Equal(t, error, string(body))
}
