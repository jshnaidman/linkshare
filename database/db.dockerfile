FROM mongo:5.0.9

WORKDIR /docker-entrypoint-initdb.d
COPY env_init_mongo.sh env_init_mongo.sh
COPY mongo_init.js mongo_init.js
WORKDIR /db/data