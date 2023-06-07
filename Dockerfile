FROM golang:1.20

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/web

EXPOSE 4000

ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=mysecretpassword
ENV DB_NAME=mydb

CMD ["./main"]