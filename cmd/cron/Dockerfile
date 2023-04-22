##
## Build
##

FROM golang:1.18 AS build

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY .env .
COPY . .

RUN CGO_ENABLED=0 go build ./cmd/cron/cron.go

##
## Deploy
##

FROM alpine

WORKDIR /app

COPY --from=build /build/cron .
COPY --from=build /build/.env .

RUN chmod +x ./cron

EXPOSE 9999/tcp

ENTRYPOINT ["./cron"]