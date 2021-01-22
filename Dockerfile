FROM golang:alpine
RUN apk update
RUN apk upgrade
RUN apk add --no-cache git make build-base

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV API_IP=127.0.0.1
ENV API_PORT=8000

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN PARCELLO_RESOURCE_DIR=./static go generate ./...
RUN go build -o dascr-board .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/dascr-board .

# Command to run when starting the container
CMD ["/dist/dascr-board"]