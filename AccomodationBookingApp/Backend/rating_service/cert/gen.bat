@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout rating-service-key.pem -out rating-service-req.pem -subj "/C=HR/AL=Albanija/L=Tirana/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in rating-service-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out rating-service-cert.pem -extfile rating-service-ext.cnf

echo "Server's signed certificate"
openssl x509 -in rating-service-cert.pem -noout -text
