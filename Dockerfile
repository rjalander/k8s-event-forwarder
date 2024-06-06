FROM golang:1.22 AS build-stage
WORKDIR /app

COPY . .
# Download Go modules
RUN go mod download

# Build the Go application
RUN go build -o /event-forwarder-main

# Deploy the application binary into a lean image
FROM alpine:3.19 AS build-release-stage

WORKDIR /

COPY --from=build-stage /event-forwarder-main /event-forwarder-main

# Expose the port
EXPOSE 8080

CMD ["/event-forwarder-main"]
