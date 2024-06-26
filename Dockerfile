FROM golang:1.20-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git openssh
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/users/server.go
CMD ["./server"]