#!/bin/bash

mongoimport --host localhost --port 27017 --db services --collection serviceList --file ./servicesData.json --type json --jsonArray


mongoimport --host localhost --port 27017 --db services --collection versions --file ./versionsData.json --type json --jsonArray