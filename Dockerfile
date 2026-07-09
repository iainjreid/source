# Builder
FROM golang:1.25 AS builder

COPY --from=node:22-bookworm /usr/local /usr/local

WORKDIR /source

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile package.json package-lock.json ./
RUN make .frontend_deps.stamp

COPY . .
RUN make

# Runtime
FROM scratch

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /source/src /src

ENTRYPOINT ["/src"]
CMD ["start"]
