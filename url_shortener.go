package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const (
	ENCODE = iota
	DECODE
	REDIRECT

	KEY_LENGTH = 6
)

type UrlShortener struct {
	urls map[string]string
}

var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GetRandomKey retourne une clé aléatoire composé de 6 lettres
func (u *UrlShortener) GetRandomKey() string {

	result := make([]byte, KEY_LENGTH)

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
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
func (u *UrlShortener) Decode(shortUrl *url.URL) (longUrl string, err error) {

	var ok bool

	key := shortUrl.RequestURI()

	// Remove 1st '/'
	if key[:1] == "/" {
		key = key[1:]
	}

	// Vérification de la présence de l'adresse courte
	longUrl, ok = u.urls[key]
	if ok == false {
		err = fmt.Errorf("short url '%s' not found", shortUrl)
		return
	}

	return
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

// handle permet une gestion centralisée des requêtes de l'API
func (u *UrlShortener) handle(requestType int) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// Vérification de la présence du body
		if r.Body == nil {
			u.onError(w, http.StatusBadRequest, "No body found", nil)
			return
		}

		// Récupération vers l'url spécifiée dans le body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			u.onError(w, http.StatusBadRequest, "Invalid body", err)
			return
		}

		// Validation de l'url
		rcvUrl, err := url.ParseRequestURI(string(body))
		if err != nil {
			u.onError(w, http.StatusBadRequest, "Invalid url", err)
			return
		}

		var status int
		var dst string

		// Gestion du type de requête
		switch requestType {
		case ENCODE:
			status = http.StatusCreated

			// Génération de l'url
			dst = "http"
			if r.TLS != nil {
				dst += "s"
			}

			dst += "://" + r.Host + "/" + u.Encode(rcvUrl)
		case DECODE:
			status = http.StatusOK
			dst, err = u.Decode(rcvUrl)
		case REDIRECT:
			status = http.StatusSeeOther
			dst, err = u.Decode(rcvUrl)
		}

		// Gestion des cas d'erreurs
		if err != nil {
			u.onError(w, http.StatusNotFound, "Invalid", err)
			return
		}

		w.Header().Set("Location", dst)
		w.WriteHeader(status)
		w.Write([]byte(dst))
	}
}

// routes initialise les routes
func (u *UrlShortener) Init() http.Handler {

	// Initialisation du random
	rand.Seed(time.Now().Unix())

	// Initialisation des urls
	u.urls = make(map[string]string)

	// Initialisation du multiplexer
	mux := http.NewServeMux()

	mux.HandleFunc("/encode", u.handle(ENCODE))
	mux.HandleFunc("/decode", u.handle(DECODE))
	mux.HandleFunc("/redirect", u.handle(REDIRECT))

	return mux
}

func main() {

	// Démarrage du serveur
	log.Fatal(http.ListenAndServe(":8080",
		new(UrlShortener).Init()))
}
