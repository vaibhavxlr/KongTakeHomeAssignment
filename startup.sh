#!/bin/bash

# start mongod server
mongod --dbpath /data/db --logpath /var/log/mongodb.log --fork

# wait for it to get started
while true; do
    if [ -S "/tmp/mongodb-27017.sock" ]; then
        echo "MongoDB is ready"
        break
    else
        echo "Waiting for MongoDB to start..."
        sleep 1
    fi
done

# import dummy data
mongoimport --host localhost --port 27017 --db services --collection serviceList --file /data/scripts/servicesData.json --type json --jsonArray
mongoimport --host localhost --port 27017 --db services --collection versions --file /data/scripts/versionsData.json --type json --jsonArray

# start the go application
/app/cmd/main