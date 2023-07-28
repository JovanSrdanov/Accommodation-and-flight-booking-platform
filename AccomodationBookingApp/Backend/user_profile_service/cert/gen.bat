@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout user-info-service-key.pem -out user-info-service-req.pem -subj "/C=HR/ST=Hrvatska/L=Zagreb/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in user-info-service-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out user-info-service-cert.pem -extfile user-info-service-ext.cnf

echo "Server's signed certificate"
openssl x509 -in user-info-service-cert.pem -noout -text
