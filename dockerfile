# ============================
# Stage 1: Build the CLI binary
# ============================
FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build fully static binary
RUN CGO_ENABLED=0 go build -o /out/onedrive ./cmd/odc

# ============================
# Stage 2: Minimal scratch image
# ============================
FROM scratch

# Create state + cache directories
# (scratch has no mkdir, so we copy empty dirs)
COPY --from=builder /out/onedrive /onedrive

ENV ODC_LOG_OUTPUT="stdout"

# Set working directory
WORKDIR /

ENTRYPOINT ["/onedrive"]
