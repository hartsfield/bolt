#!/bin/bash
bolt
touch model.json
echo '{
    "textarea": ["Message"],
    "file": ["FileElement"],
    "text": ["Title", "Email"]
}' >> model.json

bolt --streamable model.json
