<h1 align="center">
  :deciduous_tree: MCPJungle :deciduous_tree:
</h1>
<p align="center">
  Self-hosted MCP Server registry for your private AI agents
</p>

MCPJungle is a single source-of-truth registry for all [Model Context Protocol](https://modelcontextprotocol.io/introduction) based Servers running in your Organisation.

🧑‍💻 Developers use it to register & manage MCP servers and the tools they provide from a central place.

🤖 AI Agents use it to discover and consume tools across all MCP servers.

![diagram](./assets/mcpjungle-diagram/mcpjungle-diagram.png)

<p align="center">MCPJungle is the only MCP Server your agents need to connect to!</p>

## Installation

> [!WARNING]
> MCPJungle is **BETA** software.
>
> We're actively working to make it production-ready.
> You can provide your feedback by [creating an Issue](https://github.com/duaraghav8/MCPJungle/issues) in this repository.


MPCJungle is shipped as a stand-alone binary.

You can either download it from the [Releases](https://github.com/duaraghav8/MCPJungle/releases) Page or use [Homebrew](https://brew.sh/) to install it:

```bash
$ brew install duaraghav8/tap/mcpjungle
```

Verify your installation by running

```bash
$ mcpjungle version
```

> [!IMPORTANT]
> On MacOS, you will have to use homebrew because the compiled binary is not [Notarized](https://developer.apple.com/documentation/security/notarizing-macos-software-before-distribution).


## Usage

MCPJungle has a Client-Server architecture and the binary allows you to run both.

### Server
First, start a server locally:

```bash
$ mcpjungle start
```

This starts the main registry server responsible for managing all MCP servers.

The server also exposes its own MCP server at `/mcp` for AI Agents to discover and call Tools provided by the registered MCP Servers.

It relies on a database and by default, creates a SQLite DB in the current working directory.
Alternatively, you can supply a DSN for a Postgresql database to the server:

```bash
$ export DATABASE_URL=postgres://admin:root@localhost:5432/mcpjungle_db
$ mcpjungle start
```

### Client
Once the server is up, you can use the CLI to interact with it.

Register a running MCP server
```bash
$ mcpjungle register --name my_mcp_server --description "Sample mcp server providing some tools" --url http://127.0.0.1:8000/mcp
```

The registry will now start tracking this MCP server and load its tools.
All tools provided by this server are now accessible via MCPJungle:

```bash
$ mcpjungle list tools

# Check tool usage
$ mcpjungle usage my_mcp_server/calculator/multiply

# Call a tool
$ mcpjungle invoke my_mcp_server/calculator/multiply --input '{"a": 100, "b": 50}'

```

> [!NOTE]
> A tool in MCPJungle must be referred to by its canonical name which follows the pattern `<mcp-server-name>/<tool-name>`.
>
> eg- If I registered a MCP server `github` which provides a tool called `git_commit`, you can invoke it in MCPJungle using the name `github/git_commit`.


Finally, you can remove a MCP server from the registry:
```bash
$ mcpjungle deregister my_mcp_server
```

After running this, the registry will stop tracking this server and its tools will no longer be available to use.

## Development
