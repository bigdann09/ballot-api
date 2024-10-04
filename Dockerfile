FROM golang:1.23-alpine 

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ballot .

EXPOSE 8003

CMD ["/app/ballot"]