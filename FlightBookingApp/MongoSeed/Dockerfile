FROM mongo:latest

ENV DB_URI=mongodb://root:pass@mongo:27017 

COPY airports.json /airports.json

CMD mongoimport --uri $DB_URI --db flightDb --collection airports  --file ./airports.json --authenticationDatabase=admin  --jsonArray --drop
# Trebalo bi se sve staviti u jednu CMD komandu tako da konkateniraj pomocu && \ 
# pa sledeca komanda

