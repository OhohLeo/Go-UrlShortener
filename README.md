# Go-UrlShortener

Implémentation en Go d'une API Rest permettant la gestion d'un URL shortener.

Temps d'implémentation : ~4h

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
go run url_shortener.go
```

## Exemples d'utilisation

Encodage POST /encode
```bash
curl -i http://localhost:8080/encode --data "https://www.quai-des-apps.com"
```

```http
HTTP/1.1 201 Created
Date: Mon, 10 Jul 2017 09:31:05 GMT
Content-Length: 28
Content-Type: text/plain; charset=utf-8

http://localhost:8080/0L3VWH
```

Décodage POST /decode
```bash
curl -i http://localhost:8080/decode?id=0L3VWH"
```

```http
HTTP/1.1 200 OK
Date: Mon, 10 Jul 2017 09:37:46 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8

https://www.quai-des-apps.com
```

Redirection POST /redirect
```bash
curl -i http://localhost:8080/redirect?id=0L3VWH"
```

```http
HTTP/1.1 303 See Other
Location: https://www.quai-des-apps.com
Date: Mon, 10 Jul 2017 09:47:55 GMT
Content-Length: 18
Content-Type: text/plain; charset=utf-8

https://www.quai-des-apps.com
```
## Commentaires

Limites de la conception actuelle :
 - Génération d'une url courte à partir d'une clé composée de 6
   caractères avec 62 possibilités chacunes (/a-zA-Z0-9/) soit 56 800
   235 584 possibilités (62^6)

Piste d'améliorations :
 - Stockage des données en base
 - Fichier de configuration pour spécifier le domaine + port d'écoute + longueur des clés
 - Timestamp associé à chaque url pour définir un délai de péremption (afin de nettoyer les tables)
