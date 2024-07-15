
FROM golang:1.21.4 AS build


WORKDIR /app


COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN go build -o main ./cmd/main.go

FROM golang:1.21.4

WORKDIR /app

COPY --from=build /app/main .


EXPOSE 8181

CMD ["./main"]
