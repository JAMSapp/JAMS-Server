#!/bin/sh

curl --verbose -X POST -H "Content-Type: application/json" -d \
"{\"username\":\"Berdy\", \"password\":\"farts\"}" \
https://jams.howardisaslut.com/api/user

curl --verbose -X GET https://jams.howardisaslut.com/api/user\?username\=Berdy
