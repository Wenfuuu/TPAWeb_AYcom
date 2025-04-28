FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/service-email

COPY service-email/go.mod . 
COPY service-email/go.sum .

# copy proto packages
COPY pkg ../pkg

RUN go mod download

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# COPY . .
COPY service-email/. .

EXPOSE 50053

# Default command for development
CMD ["air"]