# Build
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache git && go mod download && CGO_ENABLED=0 go build -o /mcp ./cmd/mcp

# Runtime
FROM gcr.io/distroless/static
COPY --from=build /mcp /mcp
EXPOSE 8080
ENTRYPOINT ["/mcp", "serve"]