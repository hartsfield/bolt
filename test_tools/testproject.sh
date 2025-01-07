#!/bin/bash

# make an empty directoy and cd into
mkdir testApp && cd testApp

# initiate app
bolt

# create a component called test and insert it into the main page
bolt --new-component=test
bolt --insert-component=test,main

# create a nav bar and insert the necessary components
bolt --autonav about,contact,footer

# create an upload form and stream of uploaded data
touch model.json
echo '{
    "textarea": ["MyText"],
    "file": ["Media"],
    "text": ["Title", "Email"]
}' >> model.json
bolt --streamable model.json

# run the app
./autoload.sh

# go back up to the test dir and move run the test
cd ..
go run . 
