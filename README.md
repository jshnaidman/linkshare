# A site to share links with others

## HOW TO RUN

Configure a root-level .secrets file from the .secrets.sample file and fill out the missing fields.

To fill out the .secrets file, you will need a [Google Client ID](https://developers.google.com/workspace/guides/create-credentials) for the google login button on the frontend.

Since the nextjs image is mounted in development, you need to run npm install from the nextjs/ directory. You can remove the mount if you want.

```
cd nextjs
npm install
```

Then, execute `docker-compose up -d` (must have docker installed)

## Tech Stack

- Everything is containerized through docker
- Frontend is powered by a nextjs docker container
- The backend API is provided through a graphql server container written in golang
- The database used is mongoDB
