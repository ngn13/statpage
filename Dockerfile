FROM golang:1.21.4

WORKDIR /app

COPY *.go ./
COPY *.mod ./
COPY *.sum ./
COPY lib ./lib
COPY public ./public
COPY views ./views

RUN go build . 

ENTRYPOINT ["/app/statpage"]
