FROM mongo:latest

COPY json /json
COPY import.sh /import.sh

RUN chmod +x import.sh
RUN sed -i 's/\r//g' import.sh
# umesto sed-a moze search and replace \r ni sa cim
CMD ./import.sh
