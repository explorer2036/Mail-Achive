#/bin/bash

curl -X POST 'http://192.168.0.1:8000/login' \
--header 'Content-Type: application/json' \
--data '{
	"username": "tm-de",
	"password": "archiv14"
}'
