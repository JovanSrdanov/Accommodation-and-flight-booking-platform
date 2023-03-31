jsonLocation="json"
DB_URI="mongodb://root:pass@mongo:27017"
dbName="flightDb"


for file in $(ls $jsonLocation)
do
  filePath="$jsonLocation/$file"
  collectionName=$(basename -s .json $file)

  mongoimport --uri $DB_URI --db $dbName --collection $collectionName \
    --file $filePath --authenticationDatabase=admin  --jsonArray --drop
done

echo "Database seeded"

