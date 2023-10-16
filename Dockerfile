FROM golang:1.21.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /GO-POSTGRES-ESP32

EXPOSE 8080

CMD [ "/GO-POSTGRES-ESP32" ]