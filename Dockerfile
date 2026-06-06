FROM node:24-alpine AS ui-builder

WORKDIR /ui
COPY ./frontend/package.json ./frontend/package-lock.json ./
RUN npm ci

COPY ./frontend/ ./
RUN npm run build

FROM golang:1.26.3-trixie AS app-builder

WORKDIR /app
COPY go.mod ./

RUN go mod download

COPY . .
COPY --from=ui-builder /ui/dist ./frontend/dist
RUN go generate ./...
RUN CGO_ENABLED=0 go build -v -o ./bin/oaas

FROM scratch

COPY --from=app-builder /app/bin/oaas .

ENTRYPOINT ["/oaas"]
