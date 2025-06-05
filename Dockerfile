# Use debug version until MCPJungle docker images are more mature and production-ready.
FROM gcr.io/distroless/base-debian12:debug-nonroot

# The build is handled by goreleaser
# Copy the binary from the build stage
COPY mcpjungle /mcpjungle

EXPOSE 8080
ENTRYPOINT ["/mcpjungle"]

# Run the Registry Server by default
CMD ["start"]