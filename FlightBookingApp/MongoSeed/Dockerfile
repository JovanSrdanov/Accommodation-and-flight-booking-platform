FROM mongo:latest

COPY json /json
COPY import.sh /import.sh

RUN chmod +x import.sh

CMD ./import.sh


