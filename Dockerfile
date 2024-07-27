## Build

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.22.4 AS build

WORKDIR /cmd

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./

RUN go mod download
# copy all files
COPY . ./
# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app -v ./cmd/app

## Deploy
FROM scratch AS final

WORKDIR /

COPY --from=build /bin/app /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

#RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

EXPOSE 8080
EXPOSE 8090

ENTRYPOINT ["/app"]