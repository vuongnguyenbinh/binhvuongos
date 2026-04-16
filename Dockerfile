FROM golang:1.22-alpine AS builder

RUN apk add --no-cache nodejs npm git curl
RUN go install github.com/a-h/templ/cmd/templ@v0.2.793

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

WORKDIR /app
COPY . .

# Resolve dependencies
RUN go mod tidy

# Generate templ
RUN templ generate

# Build Tailwind CSS
RUN npm install -D tailwindcss@3
RUN npx tailwindcss -i web/static/css/input.css -o web/static/css/output.css --minify

# Build Go binary
RUN CGO_ENABLED=0 go build -o bin/server cmd/server/main.go

# Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates curl
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/web/static ./web/static
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
EXPOSE 3000
ENV PORT=3000
ENTRYPOINT ["./entrypoint.sh"]
