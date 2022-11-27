FROM golang:1.19-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/home-exporter .


FROM alpine:3.9

RUN apk add ca-certificates

WORKDIR /app

COPY --from=build /app/home-exporter /app/home-exporter

EXPOSE 80 2112

CMD ["/app/home-exporter"]
