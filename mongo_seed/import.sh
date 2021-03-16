#! /bin/bash

mongoimport --username=admin --password=admin --authenticationDatabase=admin --host mongodb --db test --collection options --type json --file /mongo_seed/options.json --jsonArray
mongoimport --username=admin --password=admin --authenticationDatabase=admin --host mongodb --db test --collection phoneBook --type json --file /mongo_seed/phoneBook.json --jsonArray
