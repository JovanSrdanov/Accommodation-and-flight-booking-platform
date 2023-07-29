@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout api-gateway-key.pem -out api-gateway-req.pem -subj "/C=SR/ST=Ile de France/L=Nis/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in api-gateway-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out api-gateway-cert.pem -extfile api-gateway-ext.cnf

echo "Server's signed certificate"
openssl x509 -in api-gateway-cert.pem -noout -text
