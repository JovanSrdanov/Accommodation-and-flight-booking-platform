DB_URI="mongodb://root:pass@mongo:27017 "

mongoimport --uri $DB_URI --db flightDb --collection airports  --file ./airports.json --authenticationDatabase=admin  --jsonArray --drop
mongoimport --uri $DB_URI --db flightDb --collection users  --file ./users.json --authenticationDatabase=admin  --jsonArray --drop
mongoimport --uri $DB_URI --db flightDb --collection accounts  --file ./accounts.json --authenticationDatabase=admin  --jsonArray --drop

echo "Database seeded"
