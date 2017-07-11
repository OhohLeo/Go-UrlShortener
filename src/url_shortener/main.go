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
