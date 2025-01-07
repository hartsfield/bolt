#!/bin/bash
if [ -d "testApp" ]; then          # check for an old testApp dir and delete it
  rm -rf testApp
fi

mkdir testApp && cd testApp        # make an empty testApp dir and cd into it

bolt                               # initiate bolt app
bolt --new-component=test          # create a component called test
bolt --insert-component=test,main  # insert test into the main page
bolt --autonav about,merch,contact # automatically create a nav bar with components

echo '{
    "textarea": ["MyText"],
    "file": ["Media"],
    "text": ["Title", "Email"]
}' >> model.json                   # echo a data structure into a file called "model.json"
bolt --streamable model.json       # use model.json to create the form and data stream
./autoload.sh                      # run the app
cd .. && go run main.go            # cd back up to the test dir and run the test
