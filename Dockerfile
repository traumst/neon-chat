FROM golang:1.23-bullseye as neon-chat-container
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /neon-chat-bin
EXPOSE 8080
CMD [ "/neon-chat-bin" ]
