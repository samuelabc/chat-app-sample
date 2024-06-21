FROM golang:1.20-alpine

WORKDIR /app

# Install required tools and dependencies for sqlite3
ENV CGO_ENABLED=1

RUN apk add --no-cache \
			gcc \
			musl-dev \
			sqlite-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o chat-server ./

EXPOSE 8080

CMD ["./chat-server"]
