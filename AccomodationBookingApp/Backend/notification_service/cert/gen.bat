@echo off
rem 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout notification-service-key.pem -out notification-service-req.pem -subj "/C=MA/ST=Makedonija/L=Skoplje/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

rem 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in notification-service-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out notification-service-cert.pem -extfile notification-service-ext.cnf

echo "Server's signed certificate"
openssl x509 -in notification-service-cert.pem -noout -text
