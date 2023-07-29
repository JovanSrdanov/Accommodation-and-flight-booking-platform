@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout reservation-service-key.pem -out reservation-service-req.pem -subj "/C=CR/ST=CrnaGora/L=Podgorica/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in reservation-service-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out reservation-service-cert.pem -extfile reservation-service-ext.cnf

echo "Server's signed certificate"
openssl x509 -in reservation-service-cert.pem -noout -text
