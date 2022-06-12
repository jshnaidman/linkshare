# WARNING: docker-entrypoint-initdb.d runs in alphabetical order, so keep this named alphabetically prior to mongo_init.js
sed "s/SCHEMA_VERSION/$SCHEMA_VERSION/g" -i mongo_init.js
sed "s/DB_NAME/${MONGO_INITDB_DATABASE}/g" -i mongo_init.js