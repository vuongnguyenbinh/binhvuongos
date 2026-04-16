FROM golang:1.22-alpine AS builder

RUN apk add --no-cache nodejs npm git
RUN go install github.com/a-h/templ/cmd/templ@v0.2.793

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
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/web/static ./web/static
EXPOSE 3000
ENV PORT=3000
CMD ["./server"]
