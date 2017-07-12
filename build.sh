#!/bin/sh

SRC=./src/url_shortener/*.go
DST=./bin/url_shortener_

# Suppression du dossier contenant les binaires s'il existe
if [ -d "./bin" ];then
	rm -r "./bin";
fi

# Création du dossier
mkdir "./bin";

# Binaire pour linux 386 & amd64
for arch in 386 amd64; do
	GOARCH=$arch go build -o "$DST$arch" $SRC
done

# Binaire pour armhf v7
GOARM=7	GOARCH=arm go build -o $DST"armhfv7" $SRC
