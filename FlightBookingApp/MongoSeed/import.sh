echo "DATABASE INIT START #############################################";
mongoimport --uri "mongodb://root:pass@mongo:27017" --db flightDb --collection airports  --file json/airports.json --authenticationDatabase=admin  --jsonArray --drop;
mongoimport --uri "mongodb://root:pass@mongo:27017" --db flightDb --collection users --file json/users.json --authenticationDatabase=admin  --jsonArray --drop;
mongoimport --uri "mongodb://root:pass@mongo:27017" --db flightDb --collection accounts --file json/accounts.json --authenticationDatabase=admin  --jsonArray --drop;
mongoimport --uri "mongodb://root:pass@mongo:27017" --db flightDb --collection flights --file json/flights.json --authenticationDatabase=admin  --jsonArray --drop;
echo "DATABASE INIT END #############################################";