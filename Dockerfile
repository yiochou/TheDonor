FROM golang:1.16 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY default.env ./

RUN go build -o /main

CMD [ "/main" ]