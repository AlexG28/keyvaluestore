FROM golang:1.22-alpine 

WORKDIR /app

COPY *.go .

COPY go.mod .

RUN go build -o myapp . 

CMD [ "./myapp" ]