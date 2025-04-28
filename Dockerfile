FROM golang:alpine

WORKDIR /habitsgobackend
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api \
    && go build -o ./bin/migrate ./cmd/migrate

CMD ["/habitsgobackend/bin/api"]
EXPOSE 8080