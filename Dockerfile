FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o exe/pilar -ldflags '-extldflags "-static"' pilar/pilar.go
RUN CGO_ENABLED=0 go build -o exe/viga -ldflags '-extldflags "-static"' viga/viga.go
RUN CGO_ENABLED=0 go build -o exe/sacerdote -ldflags '-extldflags "-static"' sacerdote/sacerdote.go

FROM gcr.io/distroless/static-debian11
# FROM alpine:latest
# RUN apk add libc6-compat
WORKDIR /

COPY --from=builder /app/exe /