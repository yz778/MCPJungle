# MCPJungle
Open-source, Selft-hosted MCP tool Catalogue for your AI agents

Self‑hosted registry for [Model Context Protocol](https://github.com/modelcontextprotocol/spec) compliant tools.

## Quickstart (Docker Compose)

```bash
git clone https://github.com/your-org/mcp-registry.git
cd mcp-registry
docker compose up --build
```

The API is now on `http://localhost:8080`.

### CLI usage

```bash
# Start server
mcp serve --port 8080

# Client ops
mcp register --name weather --url https://weather.example.com --type rest_api
mcp tools
mcp invoke weather --input '{"city":"Amsterdam"}'
```

---

## Local Dev

```bash
export DATABASE_URL="postgres://mcp:mcp@localhost:5432/mcp?sslmode=disable"
```