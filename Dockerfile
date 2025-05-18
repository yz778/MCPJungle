# 1️. Build stage
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache git && go mod download && CGO_ENABLED=0 go build -o /mcp-registry ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static
COPY --from=build /mcp-registry /mcp-registry
EXPOSE 8080
ENTRYPOINT ["/mcp-registry"]