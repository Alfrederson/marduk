FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o exe/zigurat -ldflags '-extldflags "-static"' .

FROM gcr.io/distroless/static-debian11
WORKDIR /

COPY --from=builder /app/exe /