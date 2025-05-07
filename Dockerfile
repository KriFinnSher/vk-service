FROM golang:1.23

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o /build ./task2/cmd/server

EXPOSE 8080
CMD ["/build"]