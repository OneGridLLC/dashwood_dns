FROM golang:1.19-alpine

RUN mkdir /app

COPY go.mod go.sum main.go database.go dns.go /app
# ^ not all working directory files should be transferred

WORKDIR /app

RUN go build -v -o agent .

CMD [ "/app/agent" ]

# copied from https://hub.docker.com/_/golang