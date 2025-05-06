FROM golang:1.24.2

WORKDIR /app

RUN go mod init github.com/Wenfuuu/tpa-web/service-chat

COPY service-chat/go.mod . 
COPY service-chat/go.sum .

COPY pkg ../pkg

RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY service-chat/. .

EXPOSE 50055 50056

CMD ["air"]