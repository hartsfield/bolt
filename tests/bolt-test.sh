#!/bin/bash
rm -rf btest
pkill btest
mkdir btest
cd btest
bolt
bolt --autonav about,contact
touch model.json
echo '{
    "textarea": ["Message"],
    "file": ["FileElement"],
    "text": ["Title", "Email"]
}' >> model.json
bolt --streamable model.json
go build -o btestbin
./btestbin &
cd ..
go run main.go localhost:9125/uploadItem/ me@me.com "my thing" "a picture of a thing"
