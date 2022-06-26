FROM golang:1.18-alpine as build-target

WORKDIR /app
USER goUser

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY ./graph/ ./graph/
COPY database/ database/
COPY utils/ utils/

RUN go build -o ./linkshare_api

FROM scratch as production-target

WORKDIR /app

COPY --from=build-target /app/linkshare_api .

EXPOSE 5000

CMD ["/app/linkshare_api"]

FROM build-target as dev-target

EXPOSE 5000
CMD ["/app/linkshare_api"]