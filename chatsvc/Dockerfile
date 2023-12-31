# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/chatsvc

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/* ./
COPY *.html /app/

EXPOSE 8080

# Must run as root to enable '-addr :80'
# On GCE, the container runs with --network=host, meaning container ports map directly to host ports
USER root:root

ENTRYPOINT ["/app/chatsvc"]