#/bin/bash

curl -X GET 'http://62.171.183.92:8000/health' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJ0bS1kZSIsImV4cCI6MTU4Njg4OTYzM30.DwCmOai1evA9BvX149vqkRHB8cvzWpIjfbbxK1YS0vY'
