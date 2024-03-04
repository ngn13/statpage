FROM golang:1.21.4 as build

WORKDIR /app

COPY *.go  ./
COPY *.mod ./
COPY *.sum ./
COPY lib   ./lib

RUN CGO_ENABLED=0 go build . 

FROM alpine as main

WORKDIR /app
COPY public                     ./public
COPY views                      ./views
COPY --from=build /app/statpage .

ENTRYPOINT ["/app/statpage"]
