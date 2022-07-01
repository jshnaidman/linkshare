FROM golang:1.18-alpine as build-target
# RUN apk add --no-cache libc6-compat
RUN addgroup -S goGroup && adduser -S goUser -G goGroup
USER goUser
WORKDIR /app
COPY --chown=goUser:goGroup . .

RUN CGO_ENABLED=0 go mod download

RUN CGO_ENABLED=0 go build -gcflags "all=-N -l" -o ./linkshare_api

FROM scratch as production-target
WORKDIR /app

COPY --from=build-target --chown=goUser:goGroup /app/linkshare_api .

EXPOSE 5000

CMD ["/app/linkshare_api"]

FROM  build-target as dev-target
USER root
RUN alias ll='ls -al'
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
COPY --chown=goUser:goGroup --from=build-target /app /app
RUN chmod 777 /app/linkshare_api
RUN export CGO_ENABLED=0 
ENTRYPOINT [ "/go/bin/dlv" ]