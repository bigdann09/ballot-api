FROM golang:1.23-alpine 

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ballot .

EXPOSE 8002

CMD ["/app/ballot"]