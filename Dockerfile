FROM golang:1.24

WORKDIR /app

COPY dnsilly/go.mod dnsilly/go.sum /app
RUN go mod download && go mod verify

COPY dnsilly /app/src
RUN go build -C src -v -o /app/dnsilly

CMD [ "/app/dnsilly" ]
