# A site to share links with others

## HOW TO RUN

Create the following two root level files:

```
# Mongo db root username goes in here
.mongodb_root_username
# Mongo db root password goes in here
.mongodb_root_password
```

Configure a root-level .secrets file from the .secrets.sample file and fill out the missing fields. 


To fill out the .secrets file, you will need a [Google Client ID](https://developers.google.com/workspace/guides/create-credentials) for the google login button on the frontend.

Then, execute docker-compose up -d (must have docker installed)


## TODO
- probably don't need two username / password files if we're storing google client_id in the .env file... can probably have one .secrets file and create a script to load envrionment variables with them (and then use MONGO_INITDB_ROOT_PASSWORD instead of MONGO_INITDB_ROOT_PASSWORD_FILE)