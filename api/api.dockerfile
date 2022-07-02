FROM golang:1.18-alpine as build-target
# RUN apk add --no-cache libc6-compat
ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0
RUN apk update \
    && apk add --no-cache \
    ca-certificates \
    git \
    && update-ca-certificates
RUN addgroup -S goGroup && adduser -S goUser -G goGroup
RUN mkdir /app_bin
RUN chown goUser:goGroup /app_bin 
USER goUser
WORKDIR /app
COPY --chown=goUser:goGroup . .

RUN CGO_ENABLED=0 go mod download

# keep binary outside /app to prevent getting overshadowed by mounted volume
# Not currently using this binary in dev though and prod will not mount volume, but good practice
RUN CGO_ENABLED=0 go build -a -o /app_bin/linkshare_api

FROM scratch as production-target
WORKDIR /app

COPY --from=build-target --chown=goUser:goGroup /app_bin/linkshare_api .

CMD ["/linkshare_api"]

FROM  build-target as dev-target
# for dev it doesn't matter if we're root, keep it simple
USER root
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/cosmtrek/air@latest
EXPOSE 8080
EXPOSE 2345
CMD ["air", "-c", ".air.toml"]