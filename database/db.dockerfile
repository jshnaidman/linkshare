FROM mongo:5.0.9

USER mongodb
WORKDIR /docker-entrypoint-initdb.d
RUN echo "alias ll='ls -l --color=auto'" >> ~/.bashrc
COPY env_init_mongo.sh env_init_mongo.sh

WORKDIR /writing
COPY mongo_init.js mongo_init.js

WORKDIR /db/data