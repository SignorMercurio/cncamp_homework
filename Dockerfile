FROM golang:1.17.3 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN mkdir bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/httpserver

FROM alpine AS production
COPY --from=builder /app/bin/* .
CMD ["./httpserver", ":8080"]