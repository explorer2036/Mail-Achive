#/bin/bash

curl -X POST 'http://62.171.183.92:8000/login' \
--header 'Content-Type: application/json' \
--data '{
	"username": "tm-de",
	"password": "archiv14"
}'
