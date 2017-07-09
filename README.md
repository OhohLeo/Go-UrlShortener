# Go-UrlShortener

Simple implémentation en Go d'une API permettant la gestion d'une UrlShortener.

POST /encode
Body: adresse à raccourcir

POST /decode
Body: adresse à décoder

POST /redirect
Body: adresse à rediriger

Limites de la conception actuelle :
 - Génération d'une url courte à partir d'une clé composée de 6 caractères avec 62 possibilités chacunes (/a-zA-Z0-9/)

Points d'améliorations :
 - Stockage des données dans une base
