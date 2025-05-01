FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/service-follow

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY service-follow/go.mod . 
COPY service-follow/go.sum .

# copy proto packages
COPY pkg ../pkg

RUN go mod download

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# COPY . .
COPY service-follow/. .

EXPOSE 50055

# Default command for development
CMD ["air"]
# CMD ["go", "run", "main.go"]