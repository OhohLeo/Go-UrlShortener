# Go-UrlShortener

Implémentation en Go d'une API Rest permettant la gestion d'un URL shortener.

L'option choisie est d'utiliser Esther-PLM pour paramétrer le port
d'écoute de l'URL shortener.

Temps de recherches, implémentation & documentation : ~10h
Testé sur RaspberryPi 2 (en utilisant bin/url_shortener_armhfv7)

## Installation
```bash
glide up
```
## Génération des binaires pour 386, amd64 & armhfv7
```bash
chmod +x script.sh
./script.sh
ls bin
```

## Lancement du programme
```bash
go build
./Go-UrlShortener -h
./Go-UrlShortener
./Go-UrlShortener -ip 127.0.0.1 -port 1234
```

## Lancement du programme en utilisant PLM (option)
```bash
esther-plm -c config.json
./Go-UrlShortener -plm http://localhost:9000
```

## Exemples d'utilisation

Encodage :
```bash
curl -i http://localhost:8080/encode --data "https://www.quai-des-apps.com"
```

```http
HTTP/1.1 201 Created
Date: Mon, 10 Jul 2017 11:13:40 GMT
Content-Length: 40
Content-Type: text/plain; charset=utf-8

http://localhost:8080/redirect?id=efghCb
```

Décodage :
```bash
curl -i http://localhost:8080/decode?id=efghCb
```

```http
HTTP/1.1 200 OK
Date: Mon, 10 Jul 2017 11:14:16 GMT
Content-Length: 29
Content-Type: text/plain; charset=utf-8

https://www.quai-des-apps.com
```

Redirection :
```bash
curl -i http://localhost:8080/redirect?id=efghCb
```

```http
HTTP/1.1 303 See Other
Location: https://www.quai-des-apps.com
Date: Mon, 10 Jul 2017 11:14:52 GMT
Content-Length: 29
Content-Type: text/plain; charset=utf-8
```

## Commentaires

Limites de la conception actuelle :
 - Génération d'une url courte à partir d'une clé composée de 6
   caractères avec 62 possibilités chacunes (/a-zA-Z0-9/) soit 56 800
   235 584 possibilités (62^6)

Piste d'améliorations :
 - Stockage des données en base
 - Fichier de configuration pour spécifier le domaine, port d'écoute, longueur des clés & emplacement des logs
 - Timestamp associé à chaque url pour définir un délai de péremption (en vue d'un nettoyage)

Fait en plus :
 - Tests unitaires sur url_shortener.go
 - Envoi d'une requête /status au PLM (non fonctionnel)