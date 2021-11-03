FROM golang:alpine

LABEL maintainer="IDerr <ibrahim@derraz.fr>"

WORKDIR /app

COPY . .

RUN go build -v -o bin/app app.go

EXPOSE 9142

CMD ["./bin/app"]