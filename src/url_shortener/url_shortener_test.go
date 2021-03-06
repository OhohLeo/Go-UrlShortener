// This file is part of Go-UrlShortener.
//
// Go-UrlShortener is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Go-UrlShortener is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Go-UrlShortener.  If not, see <http://www.gnu.org/licenses/>.
//
// Authored by OhohLeo
package main

import (
	"fmt"
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
		"https://www.quai-des-apps.com", false)
	if urlDst == nil {
		t.Fail()
		return
	}

	// Récupération de l'identifiant
	id := urlDst.Query().Get("id")

	// Décodage
	longUrl := urlOk(t, http.StatusOK,
		server.URL+"/decode?id="+id, "", false)
	if longUrl == nil {
		t.Fail()
		return
	}

	assert.Equal(t, "https://www.quai-des-apps.com", longUrl.String())

	// Redirection
	longUrl = urlOk(t, http.StatusOK,
		server.URL+"/redirect?id="+id, "", true)
	if longUrl == nil {
		t.Fail()
		return
	}

	assert.Equal(t, "https://www.quai-des-apps.com", longUrl.String())
}

func urlOk(t *testing.T, status int, dst string, data string, isRedirect bool) (redirect *url.URL) {

	rsp, err := http.Post(dst, "", strings.NewReader(data))
	if assert.Nil(t, err) == false {
		fmt.Println(err.Error())
		return
	}

	if assert.Equal(t, status, rsp.StatusCode) == false {
		return
	}

	var urlStr string
	if isRedirect {
		urlStr = rsp.Request.URL.String()
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
