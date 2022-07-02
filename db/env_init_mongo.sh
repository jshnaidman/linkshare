# WARNING: docker-entrypoint-initdb.d runs in alphabetical order, so keep this named alphabetically prior to mongo_init.js

mongo_init="/writing/mongo_init.js"
sed "s/SCHEMA_VERSION/$SCHEMA_VERSION/g" -i $mongo_init
sed "s/DB_NAME/${MONGO_INITDB_DATABASE}/g" -i $mongo_init
sed "s/SESSION_LIFETIME_SECONDS/${SESSION_LIFETIME_SECONDS}/g" -i $mongo_init
sed "1s/^/use ${MONGO_INITDB_DATABASE}\n/" -i $mongo_init # add to top of file to use DB_NAME as DB
mongo < $mongo_init