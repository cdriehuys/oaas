FROM golang:1.26.3-trixie AS builder

WORKDIR /app
COPY go.mod ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -v -o ./bin/oaas

FROM scratch

COPY --from=builder /app/bin/oaas .

ENTRYPOINT ["/oaas"]
