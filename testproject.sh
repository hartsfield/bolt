#!/bin/bash

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
    "textarea": ["Message"],
    "file": ["FileElement"],
    "text": ["Title", "Email"]
}' >> model.json
bolt --streamable model.json

# run the app
./autoload.sh
