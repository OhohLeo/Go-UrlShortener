package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {

	var ip, plm string
	var port int

	// Récupération des arguments
	flag.StringVar(&ip, "ip", "localhost", "listening ip")
	flag.IntVar(&port, "port", 8080, "listening port")
	flag.StringVar(&plm, "plm", "", "http address to connect with PLM")

	flag.Parse()

	// Inscription à PLM
	if plm != "" {

		var err error

		service := &Service{
			Name: "url_shortener",
			Dst:  plm,
		}

		ip, port, err = service.Register()
		if err != nil {
			log.Fatalf("Issue during PLM registering: %s", err.Error())
			return
		}

		// Non requis : pour test...
		// pb avec la déclaration de URI => 'Check URI NA every 1s' ?
		// service.Update(STATUS_UP)
	}

	param := ip + ":" + strconv.Itoa(port)

	log.Printf("URL shortener listening '%s' ...", param)

	// Démarrage du serveur
	log.Fatal(http.ListenAndServe(param, new(UrlShortener).Init()))
}
