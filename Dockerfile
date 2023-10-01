# build backend service
FROM golang:1.21 AS build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app .


# build runtime image
FROM debian:bookworm-slim
COPY --from=build /usr/local/bin/app /usr/local/bin/app
EXPOSE 8000
CMD ["app"]
