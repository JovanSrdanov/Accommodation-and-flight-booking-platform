@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout accommodation-service-key.pem -out accommodation-service-req.pem -subj "/C=BU/ST=Bugarska/L=Sofija/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in accommodation-service-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out accommodation-service-cert.pem -extfile accommodation-service-ext.cnf

echo "Server's signed certificate"
openssl x509 -in accommodation-service-cert.pem -noout -text
