FROM mongo:latest

COPY airports.json /airports.json
COPY users.json /users.json
COPY accounts.json /accounts.json

COPY import.sh /import.sh

RUN chmod +x import.sh

CMD ./import.sh


