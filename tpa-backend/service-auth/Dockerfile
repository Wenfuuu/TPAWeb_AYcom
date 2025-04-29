FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/service-auth

COPY service-auth/go.mod . 
COPY service-auth/go.sum .

# copy proto packages
COPY pkg ../pkg

RUN go mod download

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# COPY . .
COPY service-auth/. .

EXPOSE 50052

# Default command for development
CMD ["air"]
# CMD ["go", "run", "main.go"]