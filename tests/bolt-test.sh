#!/bin/bash
mkdir btest
cd btest
bolt
bolt --autonav about,contact,footer
touch model.json
echo '{
    "file": ["FileElement"],
    "text": ["Title","Year","Price"],
    "textarea": ["About"]
}' >> model.json
bolt --streamable model.json
go build -o btestbin
./btestbin &
