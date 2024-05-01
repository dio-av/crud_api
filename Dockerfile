FROM golang:1.22-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD [ "app/crud_api" ]