FROM golang:1.17.1 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN make build

FROM alpine AS production
COPY --from=builder /app/bin/* .
CMD ["./httpserver", ":8080"]