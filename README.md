# A site to share links with others

## HOW TO RUN

Configure a root-level .secrets file from the .secrets.sample file and fill out the missing fields. 


To fill out the .secrets file, you will need a [Google Client ID](https://developers.google.com/workspace/guides/create-credentials) for the google login button on the frontend.

Then, execute docker-compose up -d (must have docker installed)


## Tech Stack

- Everything is containerized through docker
- Frontend is powered by a nextjs docker container
- The backend API is provided through a graphql server container written in golang
- The database used is mongoDB