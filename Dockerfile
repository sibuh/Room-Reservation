FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod verify
ADD . .
RUN go build -o ./bin cmd/main.go
CMD ["/app/bin"]
