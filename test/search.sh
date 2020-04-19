curl -X POST 'http://62.171.183.92:8000/search' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJ0bS1kZSIsImV4cCI6MTU4NzI3OTI4MH0.Tw__rbVN4uMSxbAPStuVbWNk9gcpc5u7kbAFsQhiX3o' \
--data '{
	"query": "mit",
	"skip": 0,
	"take": 100
}'
