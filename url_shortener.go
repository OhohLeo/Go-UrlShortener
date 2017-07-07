package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	ENCODE = iota
	DECODE
	REDIRECT
)

type UrlShortener struct {
	urls map[string]string
}

// Encode permet d'obtenir l'url réduite
func (u *UrlShortener) Encode(host string, url *url.URL) (status int, shortUrl string, err error) {

	// Algorithme utilisé pour réduire l'url [A AMELIORER]
	hash := md5.New()
	io.WriteString(hash, url.EscapedPath())
	key := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	shortUrl = "http://" + host + "/" + key

	// Vérification que la shortUrl n'est pas déjà utilisée
	originalUrl, ok := u.urls[key]
	if ok {

		// Rejet en cas de détection de collisions
		if originalUrl != url.String() {
			err = fmt.Errorf(
				"similar short url found for '%s' and already stored '%s'",
				url, originalUrl)
			return
		}

		// Déjà existant
		status = http.StatusNotModified
		return
	}

	// Stockage de l'url
	status = http.StatusCreated
	u.urls[key] = url.String()

	return
}

// Decode permet d'obtenir l'url originale à partir de l'url réduite
func (u *UrlShortener) Decode(shortUrl *url.URL) (url string, err error) {

	var ok bool

	key := shortUrl.RequestURI()

	// Remove 1st '/'
	if key[:1] == "/" {
		key = key[1:]
	}

	// Vérification de la présence de l'adresse courte
	url, ok = u.urls[key]
	if ok == false {
		err = fmt.Errorf("no short url '%s' found", shortUrl)
		return
	}

	return
}

// onError est appelé en cas d'erreur et retourne une erreur de type HTTP BadRequest
func (u *UrlShortener) onError(w http.ResponseWriter, msg string, err error) {

	if err != nil {
		msg += " " + err.Error()
	}

	// Log du message d'erreur
	log.Println(msg)

	// Renvoie de l'erreur
	http.Error(w, msg, http.StatusBadRequest)
}

// handle permet une gestion centralisée des requêtes de l'API
func (u *UrlShortener) handle(requestType int) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// Vérification de la présence du body
		if r.Body == nil {
			u.onError(w, "No body found", nil)
			return
		}

		// Récupération vers l'url spécifiée dans le body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			u.onError(w, "Invalid body", err)
			return
		}

		// Validation de l'url
		rcvUrl, err := url.ParseRequestURI(string(body))
		if err != nil {
			u.onError(w, "Invalid url", err)
			return
		}

		var status int
		var dst string

		// Gestion du type de requête
		switch requestType {
		case ENCODE:
			status, dst, err = u.Encode(r.Host, rcvUrl)
		case DECODE:
			status = http.StatusOK
			dst, err = u.Decode(rcvUrl)
		case REDIRECT:
			status = http.StatusMovedPermanently
			dst, err = u.Decode(rcvUrl)
		}

		// Gestion des cas d'erreurs
		if err != nil {
			u.onError(w, "Invalid", err)
			return
		}

		w.Header().Set("Location", dst)
		w.WriteHeader(status)
		w.Write([]byte(dst))
	}
}

// Start déclare les routes et démarre le serveur web
func (u *UrlShortener) Start(address string) error {

	// Initialisation des urls
	u.urls = make(map[string]string)

	// Gestion des routes
	http.HandleFunc("/encode", u.handle(ENCODE))
	http.HandleFunc("/decode", u.handle(DECODE))
	http.HandleFunc("/redirect", u.handle(REDIRECT))

	// Lancement du serveur Web
	return http.ListenAndServe(address, nil)
}

func main() {

	// Démarrage du serveur
	log.Fatal(new(UrlShortener).Start(":8080"))
}
