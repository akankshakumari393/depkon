##
## Build binary
##
FROM golang:1.18-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

COPY ./controller ./controller

COPY ./pkg ./pkg

RUN go mod vendor

RUN CGO_ENABLED=0 go build

##
## RUN the binary
##

FROM alpine

COPY --from=build /app/depkon /usr/local/bin

#USER root:root

ENTRYPOINT ["depkon"]
