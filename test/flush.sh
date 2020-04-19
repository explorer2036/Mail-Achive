#/bin/bash

curl -X GET 'http://192.168.0.1:8000/flush' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJ0bS1kZSIsImV4cCI6MTU4NzMxNDUyMn0.0b9NZnnweVWDoSLmGENnGoavV5IQdWzTKPjXfO9DI-8'
