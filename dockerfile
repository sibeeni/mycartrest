FROM golang:latest
RUN mkdir /app
WORKDIR /app

ADD . /app
RUN go get github.com/gin-gonic/gin
RUN go get github.com/gin-contrib/cors
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/stretchr/testify/assert
RUN go build ./base.go

EXPOSE 3001

CMD ["./base"]
