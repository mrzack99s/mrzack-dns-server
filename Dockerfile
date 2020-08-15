from golang:1.14.7-alpine

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN go build -o bin/mrzack-dns-server main.go

EXPOSE 53/udp

ENTRYPOINT [ "./bin/mrzack-dns-server" ]
