FROM golang:1.26

COPY . /tests
WORKDIR /tests

RUN go mod download

ENTRYPOINT ["go", "test", "-v", "./..."]
