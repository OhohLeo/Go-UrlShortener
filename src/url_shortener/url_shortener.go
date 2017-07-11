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
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	DECODE = iota
	REDIRECT

	KEY_LENGTH = 6
)

var LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var CHECK_ID = regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)

type UrlShortener struct {
	urls map[string]string
}

// Init initialise la fonctionnalité random, le stockage des urls et les routes
func (u *UrlShortener) Init() http.Handler {

	// Initialisation du random
	rand.Seed(time.Now().Unix())

	// Initialisation des urls
	u.urls = make(map[string]string)

	// Initialisation du multiplexer
	mux := http.NewServeMux()

	mux.HandleFunc("/encode", u.handleEncode)
	mux.HandleFunc("/decode", u.handle(DECODE))
	mux.HandleFunc("/redirect", u.handle(REDIRECT))

	return mux
}

// handleEncode gère la génération d'une clé aléatoire associé à l'url passée en paramètre
func (u *UrlShortener) handleEncode(w http.ResponseWriter, r *http.Request) {

	// Vérification de la présence du body
	if r.Body == nil {
		u.onError(w, http.StatusBadRequest, "No body found", nil)
		return
	}

	// Récupération de l'url longue spécifiée dans le body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.onError(w, http.StatusBadRequest, "Invalid body", err)
		return
	}

	// Validation de l'url longue
	rcvUrl, err := url.ParseRequestURI(string(body))
	if err != nil {
		u.onError(w, http.StatusBadRequest, "Invalid url", err)
		return
	}

	// Génération de l'url courte
	dst := "http"
	if r.TLS != nil {
		dst += "s"
	}
	dst += "://" + r.Host + "/redirect?id=" + u.Encode(rcvUrl)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(dst))
}

// handle gère les requêtes de type decode & redirect
func (u *UrlShortener) handle(requestType int) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		status, dst, err := u.Decode(r.URL)
		if err != nil {
			u.onError(w, status, "Invalid", err)
			return
		}

		var body []byte

		switch requestType {
		case DECODE:
			status = http.StatusOK
			body = []byte(dst)
		case REDIRECT:
			status = http.StatusSeeOther
			w.Header().Set("Location", dst)
		}

		w.WriteHeader(status)
		w.Write(body)
	}
}

// onError est appelé en cas d'erreur et retourne une erreur de type HTTP BadRequest
func (u *UrlShortener) onError(w http.ResponseWriter, status int, msg string, err error) {

	if err != nil {
		msg += " " + err.Error()
	}

	// Log du message d'erreur
	log.Println(msg)

	// Renvoie de l'erreur
	http.Error(w, msg, status)
}

// Encode permet d'obtenir l'url réduite
func (u *UrlShortener) Encode(longUrl *url.URL) (shortUrl string) {

	// Génération d'une clé aléatoire
	for {
		shortUrl = u.GetRandomKey()

		// Vérification que la shortUrl n'est pas déjà utilisée
		_, ok := u.urls[shortUrl]
		if ok == false {
			break
		}
	}

	// Stockage de la relation url courte => url longue
	u.urls[shortUrl] = longUrl.String()

	return
}

// Decode permet d'obtenir l'url originale à partir de l'url réduite
func (u *UrlShortener) Decode(shortUrl *url.URL) (status int, longUrl string, err error) {

	var ok bool

	key := shortUrl.Query().Get("id")

	// Validation de l'id
	if CHECK_ID.MatchString(key) == false {
		status = http.StatusBadRequest
		err = fmt.Errorf("id '%s'", key)
		return
	}

	// Récupération de l'adresse longue
	longUrl, ok = u.urls[key]
	if ok == false {
		status = http.StatusNotFound
		err = fmt.Errorf("id '%s' not found", key)
		return
	}

	return
}

// GetRandomKey retourne une clé aléatoire composé de 6 lettres
func (u *UrlShortener) GetRandomKey() string {

	result := make([]byte, KEY_LENGTH)

	for i := range result {
		result[i] = LETTERS[rand.Intn(len(LETTERS))]
	}

	return string(result)
}
