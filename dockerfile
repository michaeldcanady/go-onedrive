FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /out/onedrive ./cmd/odc

FROM docker.io/library/golang:1.25-alpine

COPY --from=builder /out/onedrive /onedrive

ENV ODC_LOG_OUTPUT="stdout"

WORKDIR /

ENTRYPOINT ["/onedrive"]
