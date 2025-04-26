FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/service-user

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY service-user/go.mod . 
COPY service-user/go.sum .

# copy proto packages
COPY pkg ../pkg

RUN go mod download

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# COPY . .
COPY service-user/. .

EXPOSE 8080

# Default command for development
CMD ["air"]