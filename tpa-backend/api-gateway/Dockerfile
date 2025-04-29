FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/api-gateway

# COPY go.mod ./
# COPY go.sum ./
# COPY ../service-user ../service-user
# RUN go mod download
COPY api-gateway/go.mod . 
COPY api-gateway/go.sum .

# copy needed services
# COPY service-user ../service-user

# copy proto packages
COPY pkg ../pkg

RUN go mod download

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# COPY . .
COPY api-gateway/. .

EXPOSE 8080

# Default command for development
CMD ["air"]
# CMD ["go", "run", "cmd/main.go"]