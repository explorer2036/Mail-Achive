#/bin/bash

curl -X GET 'http://192.168.0.1:8000/refresh' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJ0bS1kZSIsImV4cCI6MTU4Njk2NDU1MX0.tG_wTsQzyHo6VHvaUgP5rtZ1AcZSG74Em4QhTMFNSAE'
