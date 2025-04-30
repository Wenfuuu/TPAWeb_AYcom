FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/service-thread

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY service-thread/go.mod . 
COPY service-thread/go.sum .

# copy proto packages
COPY pkg ../pkg

RUN go mod download

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# COPY . .
COPY service-thread/. .

EXPOSE 50054

# Default command for development
CMD ["air"]
# CMD ["go", "run", "main.go"]