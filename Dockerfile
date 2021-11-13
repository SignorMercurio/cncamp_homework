FROM golang:1.17.1 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN mkdir bin && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/httpserver

FROM alpine AS production
COPY --from=builder /app/bin/* .
CMD ["./httpserver", ":8080"]