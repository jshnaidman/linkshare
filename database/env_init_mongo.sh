# WARNING: docker-entrypoint-initdb.d runs in alphabetical order, so keep this named alphabetically prior to mongo_init.js

mongo_init="/writing/mongo_init.js"
sed "s/SCHEMA_VERSION/$SCHEMA_VERSION/g" -i $mongo_init
sed "s/DB_NAME/${MONGO_INITDB_DATABASE}/g" -i $mongo_init
sed "s/URL_CAPACITY/${URL_CAPACITY}/g" -i $mongo_init 
sed "s/RUN_ONCE_BULK_LOAD_PAGE_AMOUNT/${RUN_ONCE_BULK_LOAD_PAGE_AMOUNT}/g" -i $mongo_init
sed "1s/^/use ${MONGO_INITDB_DATABASE}\n/" -i $mongo_init # add to top of file
mongo < $mongo_init