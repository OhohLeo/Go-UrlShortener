#!/bin/bash

dst=`dirname $0`

# Suppression du dossier contenant les binaires s'il existe
if [ -d "$dst/bin" ];then
	rm -r "$dst/bin";
fi

# Création du dossier
mkdir "$dst/bin";

# Binaire pour linux 386 & amd64
for arch in 386 amd64; do
	GOARCH=$arch go build -o "$dst/bin/url_shortener_$arch"
done

# Binaire pour armhf v7
GOARM=7	GOARCH=arm go build -o "$dst/bin/url_shortener_armhfv7"
