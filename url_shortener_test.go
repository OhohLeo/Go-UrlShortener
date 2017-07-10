package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUrlShortener(t *testing.T) {

	// Aucun serveur
	_, err := http.Post("/encode", "", strings.NewReader(""))
	if assert.NotNil(t, err) == false {
		t.Fail()
	}

	// Lancement du serveur
	server := httptest.NewServer(new(UrlShortener).Init())
	defer server.Close()

	// Requêtes invalides
	urlError(t, http.StatusBadRequest,
		server.URL+"/encode", "",
		"Invalid url parse : empty url\n")
	urlError(t, http.StatusBadRequest,
		server.URL+"/encode", "www.test.bzh",
		"Invalid url parse www.test.bzh: invalid URI for request\n")

	for _, path := range []string{"decode", "redirect"} {
		urlError(t, http.StatusBadRequest,
			server.URL+"/"+path+"", "",
			"Invalid id ''\n")
		urlError(t, http.StatusBadRequest,
			server.URL+"/"+path+"?id=123", "",
			"Invalid id '123'\n")
		urlError(t, http.StatusNotFound,
			server.URL+"/"+path+"?id=123456", "",
			"Invalid id '123456' not found\n")
	}

	// Encodage
	urlDst := urlOk(t, http.StatusCreated, server.URL+"/encode",
		"http://www.test.bzh", false)
	if urlDst == nil {
		t.Fail()
		return
	}

	// Récupération de l'identifiant
	id := urlDst.Query().Get("id")

	// Décodage
	longUrl := urlOk(t, http.StatusOK, server.URL+"/decode?id="+id, "", false)

	assert.Equal(t, "http://www.test.bzh", longUrl.String())

}

func urlOk(t *testing.T, status int, dst string, data string, checkLocation bool) (redirect *url.URL) {

	rsp, err := http.Post(dst, "", strings.NewReader(data))
	if assert.Nil(t, err) == false {
		return
	}

	assert.Equal(t, status, rsp.StatusCode)

	var urlStr string
	if checkLocation {
		urlStr = rsp.Header.Get("Location")
	} else {
		body, err := ioutil.ReadAll(rsp.Body)
		if assert.Nil(t, err, "no error on getting body") == false {
			t.Fail()
			return
		}

		urlStr = string(body)
	}

	rcvUrl, err := url.ParseRequestURI(urlStr)
	if assert.Nil(t, err) == false {
		t.Fail()
		return
	}

	redirect = rcvUrl
	return
}

func urlError(t *testing.T, status int, dst string, location string, error string) {

	rsp, err := http.Post(dst, "", strings.NewReader(location))
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
