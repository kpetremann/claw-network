FROM golang:1.18-bullseye AS builder

WORKDIR /go/src/
COPY ./ /go/src/
RUN mkdir build
RUN go build -o /go/src/build/claw-network

FROM debian:bullseye-slim AS runner
WORKDIR /app/claw-network/
RUN mkdir examples
COPY --from=builder /go/src/build/claw-network ./
CMD ["/app/claw-network/claw-network"]
