#Getting service name
read -p "Enter service name (snake case): " serviceName
mkdir ../$serviceName
cd ../$serviceName

#Generating directories
mkdir domain
cd domain
mkdir domain service repository error
touch domain/destroyMe service/destroyMe repository/destroyMe error/destroyMe
cd ..
mkdir persistence
cd persistence
mkdir error repository
touch error/destroyMe repository/destroyMe
cd ..
mkdir configuration
touch configuration/destroyMe
mkdir communication
mkdir -p communication/handler
touch communication/handler/destroyMe

#Empty dockerfile
touch Dockerfile

#Initializing go module
go mod init $serviceName

# Generating empty proto directory with empty proto file
mkdir ../common/proto/$serviceName
cd ../common/proto/$serviceName
mkdir generated
touch generated/destroyMe
touch $serviceName.proto

