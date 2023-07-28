@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout auth-service-key.pem -out auth-service-req.pem -subj "/C=FR/ST=Ile de France/L=Paris/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in auth-service-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out auth-service-cert.pem -extfile auth-service-ext.cnf

echo "Server's signed certificate"
openssl x509 -in auth-service-cert.pem -noout -text
