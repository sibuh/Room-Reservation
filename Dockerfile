FROM golang:1.23.3
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod verify
ADD . .
RUN go build -o binary cmd/main.go
CMD ["binary"]
