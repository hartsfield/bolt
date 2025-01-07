#!/bin/bash

mkdir testApp && cd testApp        # make an empty directoy and cd into

bolt                               # initiate app
bolt --new-component=test          # create a component called test
bolt --insert-component=test,main  # insert test into the main page
bolt --autonav about,merch,contact # create a nav bar and insert the necessary components

echo '{
    "textarea": ["MyText"],
    "file": ["Media"],
    "text": ["Title", "Email"]
}' >> model.json                   # echo the data into a file "model.json"
bolt --streamable model.json       # use model.json to create the form and data stream

./autoload.sh                      # run the app
cd .. && go run main.go            # cd back up to the test dir and run the test
