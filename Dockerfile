# ---------- Build stage ----------
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy go.mod and go.sum first (helps with caching)(acording to Lumo)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# We need CGO for sqlite
# We are on linux and the same arch as the runtime machine, so there is no need to specify compilation target
RUN CGO_ENABLED=1 go build -o jinovatka-main ./main.go

# ---------- Runtime stage ----------
FROM debian:bookworm

# Copy the compiled binary from the builder stage
COPY --from=builder /app/jinovatka-main /usr/local/bin/jinovatka-main

# Expose the port the app listens on
EXPOSE 8080

# Allow files under /mnt to be writable
RUN chmod 777 /mnt

# Drop privileges:
#   - Use a high, non‑zero UID/GID (e.g., 65532) that isn’t mapped to any real user.
#   - No /etc/passwd exists, but the kernel will still enforce the UID.
USER 65532:65532

# Path for the database file used by sqlite
# ENV DB_PATH=/mnt/storage.db

# Run the binary
ENTRYPOINT ["jinovatka-main"]