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

RUN CGO_ENABLED=0 go build ./cmd/auth/auth.go

##
## Deploy
##

FROM alpine

WORKDIR /app

COPY --from=build /build/auth .
COPY --from=build /build/.env .

RUN chmod +x ./auth

EXPOSE 5555/tcp

ENTRYPOINT ["./auth"]