#/bin/bash

curl -X POST http://62.171.183.92:8000/upload \
--header "Content-Type: multipart/form-data" \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJ0bS1kZSIsImV4cCI6MTU4NzI3OTI4MH0.Tw__rbVN4uMSxbAPStuVbWNk9gcpc5u7kbAFsQhiX3o' \
-F "file=@mails.zip" 
