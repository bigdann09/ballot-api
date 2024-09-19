FROM golang:1.23-alpine 

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o ballot .

EXPOSE 8002

CMD ["/app/ballot"]