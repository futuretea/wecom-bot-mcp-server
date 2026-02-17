# WeCom Bot MCP Server

[![Build](https://github.com/futuretea/wecom-bot-mcp-server/actions/workflows/build.yaml/badge.svg)](https://github.com/futuretea/wecom-bot-mcp-server/actions/workflows/build.yaml)
[![GitHub License](https://img.shields.io/github/license/futuretea/wecom-bot-mcp-server)](https://github.com/futuretea/wecom-bot-mcp-server/blob/main/LICENSE)
[![npm](https://img.shields.io/npm/v/@futuretea/wecom-bot-mcp-server)](https://www.npmjs.com/package/@futuretea/wecom-bot-mcp-server)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/futuretea/wecom-bot-mcp-server?sort=semver)](https://github.com/futuretea/wecom-bot-mcp-server/releases/latest)

[Features](#features) | [Getting Started](#getting-started) | [Configuration](#configuration) | [Tools](#tools) | [Development](#development)

## Features <a id="features"></a>

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for [WeCom (WeChat Work)](https://work.weixin.qq.com/) bot webhooks.

- **Text Messages**: Send plain text with @mention support (by user ID or mobile number)
- **Markdown Messages**: Send rich formatted messages using Markdown syntax
- **Image Messages**: Send base64-encoded images (JPG/PNG, up to 2MB)
- **News Messages**: Send article list cards (1–8 articles with title, description, URL, cover image)
- **Template Cards**: Send text notice and news notice template cards with emphasis content, key-value lists, jump links, and card actions
- **File Upload**: Upload files to WeCom server (up to 20MB) and obtain `media_id`
- **Dual Transport**: Stdio mode for MCP client integration or HTTP/SSE mode for network access
- **Cross-platform**: Native binaries for Linux, macOS, Windows (amd64/arm64), npm package, and Docker images

## Getting Started <a id="getting-started"></a>

### Requirements

- A WeCom bot webhook key (from the webhook URL: `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY`)

### Claude Code

```shell
claude mcp add wecom-bot -- npx @futuretea/wecom-bot-mcp-server@latest \
  --wecom-bot-key YOUR_KEY
```

### VS Code / Cursor

Add to `.vscode/mcp.json` or `~/.cursor/mcp.json`:

```json
{
  "servers": {
    "wecom-bot": {
      "command": "npx",
      "args": [
        "-y",
        "@futuretea/wecom-bot-mcp-server@latest",
        "--wecom-bot-key",
        "YOUR_KEY"
      ]
    }
  }
}
```

### Docker

Stdio mode:

```shell
docker run --rm -i ghcr.io/futuretea/wecom-bot-mcp-server:latest \
  --wecom-bot-key YOUR_KEY
```

HTTP/SSE mode:

```shell
docker run --rm -p 8080:8080 ghcr.io/futuretea/wecom-bot-mcp-server:latest \
  --port 8080 --wecom-bot-key YOUR_KEY
```

## Configuration <a id="configuration"></a>

Configuration can be set via CLI flags, environment variables, or a config file.

**Priority (highest to lowest):**
1. Command-line flags
2. Environment variables (prefix: `WECOM_MCP_`)
3. Configuration file
4. Default values

### CLI Options

```shell
npx @futuretea/wecom-bot-mcp-server@latest --help
```

| Option | Description | Default |
|--------|-------------|---------|
| `--config` | Config file path (YAML) | |
| `--port` | Port for HTTP/SSE mode (0 = stdio mode) | `0` |
| `--sse-base-url` | Public base URL for SSE endpoint | |
| `--log-level` | Log level (0-9) | `5` |
| `--wecom-bot-key` | WeCom bot webhook key (**required**) | |
| `--enabled-tools` | Specific tools to enable | |
| `--disabled-tools` | Specific tools to disable | |

### Configuration File

Create `config.yaml`:

```yaml
port: 0  # 0 for stdio, or set a port like 8080 for HTTP/SSE

log_level: 5

# Get the key from your WeCom bot webhook URL:
# https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY_HERE
wecom_bot_key: your-bot-key-here

# enabled_tools: []
# disabled_tools: []
```

### Environment Variables

Use `WECOM_MCP_` prefix with underscores:

```shell
WECOM_MCP_PORT=8080
WECOM_MCP_WECOM_BOT_KEY=your-key
WECOM_MCP_LOG_LEVEL=5
```

### HTTP/SSE Mode

Run with a port number for network access:

```shell
wecom-bot-mcp-server --port 8080 --wecom-bot-key YOUR_KEY
```

Endpoints:
- `/healthz` - Health check
- `/mcp` - Streamable HTTP endpoint
- `/sse` - Server-Sent Events endpoint
- `/message` - Message endpoint for SSE clients

With a public URL behind a proxy:

```shell
wecom-bot-mcp-server --port 8080 \
  --sse-base-url https://your-domain.com:8080 \
  --wecom-bot-key YOUR_KEY
```

## Tools <a id="tools"></a>

Use `--enabled-tools` / `--disabled-tools` for fine-grained control.

<details>
<summary>send_text</summary>

Send a text message via WeCom bot webhook. Supports @mentioning users by user ID or mobile number.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `content` | string | Yes | The text content to send. Maximum 2048 bytes. |
| `mentioned_list` | string[] | No | List of user IDs to @mention. Use `"@all"` to mention everyone. |
| `mentioned_mobile_list` | string[] | No | List of mobile numbers to @mention. Use `"@all"` to mention everyone. |

**Example:**

```json
{
  "content": "Hello, this is a test message.",
  "mentioned_list": ["user1", "user2"],
  "mentioned_mobile_list": ["13800138000"]
}
```

</details>

<details>
<summary>send_markdown</summary>

Send a markdown message via WeCom bot webhook. Supports headings, bold, links, quotes, and more.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `content` | string | Yes | The markdown content to send. Maximum 4096 bytes. |

**Example:**

```json
{
  "content": "# Heading\n**Bold text**\n> Quote\n[Link](https://example.com)"
}
```

</details>

<details>
<summary>send_image</summary>

Send an image message via WeCom bot webhook using base64-encoded image data.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `base64` | string | Yes | Base64-encoded image content. Max image size: 2MB. Supported formats: JPG, PNG. |
| `md5` | string | Yes | MD5 hash of the original image content (before base64 encoding). |

**Example:**

```json
{
  "base64": "iVBORw0KGgoAAAANS...",
  "md5": "d41d8cd98f00b204e9800998ecf8427e"
}
```

</details>

<details>
<summary>send_news</summary>

Send a news (article list) message via WeCom bot webhook. Supports 1–8 articles.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `articles` | object[] | Yes | Array of news articles (1–8 items). |
| `articles[].title` | string | Yes | Article title. |
| `articles[].description` | string | No | Article description. |
| `articles[].url` | string | Yes | Article link URL. |
| `articles[].picurl` | string | No | Article cover image URL. |

**Example:**

```json
{
  "articles": [
    {
      "title": "Breaking News",
      "description": "Something important happened.",
      "url": "https://example.com/news/1",
      "picurl": "https://example.com/cover.jpg"
    }
  ]
}
```

</details>

<details>
<summary>send_text_notice_card</summary>

Send a text notice template card via WeCom bot webhook. Rich card layout with emphasis content, key-value list, and jump links.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `main_title` | string | Yes | Main title of the card. |
| `main_title_desc` | string | No | Description text below the main title. |
| `sub_title` | string | No | Subtitle text displayed in the card body. |
| `source` | object | No | Source information displayed at the top of the card. |
| `source.icon_url` | string | No | URL of the source icon. |
| `source.desc` | string | No | Source description text. |
| `emphasis_content` | object | No | Emphasized content area (large text). |
| `emphasis_content.title` | string | No | Emphasis title (displayed in large font). |
| `emphasis_content.desc` | string | No | Emphasis description. |
| `horizontal_content_list` | object[] | No | Key-value pairs displayed horizontally. |
| `horizontal_content_list[].keyname` | string | Yes | Key name (label). |
| `horizontal_content_list[].value` | string | No | Value text. |
| `jump_list` | object[] | No | Jump links displayed at the bottom of the card. |
| `jump_list[].title` | string | Yes | Jump link title. |
| `jump_list[].url` | string | Yes | Jump link URL. |
| `card_action` | object | Yes | Card click action. |
| `card_action.url` | string | Yes | URL to open when the card is clicked. |

**Example:**

```json
{
  "main_title": "Deployment Notification",
  "main_title_desc": "Production environment",
  "source": {
    "icon_url": "https://example.com/icon.png",
    "desc": "CI/CD Pipeline"
  },
  "emphasis_content": {
    "title": "SUCCESS",
    "desc": "Build #1234"
  },
  "horizontal_content_list": [
    { "keyname": "Branch", "value": "main" },
    { "keyname": "Commit", "value": "abc1234" }
  ],
  "jump_list": [
    { "title": "View Details", "url": "https://example.com/build/1234" }
  ],
  "card_action": {
    "url": "https://example.com/build/1234"
  }
}
```

</details>

<details>
<summary>send_news_notice_card</summary>

Send a news notice template card via WeCom bot webhook. Card layout with a large cover image.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `main_title` | string | Yes | Main title of the card. |
| `main_title_desc` | string | No | Description text below the main title. |
| `card_image_url` | string | Yes | URL of the card cover image. |
| `source` | object | No | Source information displayed at the top of the card. |
| `source.icon_url` | string | No | URL of the source icon. |
| `source.desc` | string | No | Source description text. |
| `card_action` | object | Yes | Card click action. |
| `card_action.url` | string | Yes | URL to open when the card is clicked. |

**Example:**

```json
{
  "main_title": "New Feature Released",
  "main_title_desc": "v2.0.0 is now available",
  "card_image_url": "https://example.com/banner.jpg",
  "source": {
    "icon_url": "https://example.com/icon.png",
    "desc": "Product Team"
  },
  "card_action": {
    "url": "https://example.com/changelog"
  }
}
```

</details>

<details>
<summary>upload_file</summary>

Upload a file to WeCom server via bot webhook. Returns a `media_id` that can be used to send file messages.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `filename` | string | Yes | Name of the file to upload. |
| `base64_data` | string | Yes | Base64-encoded file content. Max file size: 20MB. |

**Example:**

```json
{
  "filename": "report.pdf",
  "base64_data": "JVBERi0xLjQK..."
}
```

</details>

## Development <a id="development"></a>

### Build

```shell
make build
```

### Test

```shell
go test ./...
```

### Lint

```shell
make lint
```

### Run with mcp-inspector

```shell
npx @modelcontextprotocol/inspector@latest -- npx @futuretea/wecom-bot-mcp-server@latest
```

## Contributing

Contributions are welcome! Please open an issue or pull request on [GitHub](https://github.com/futuretea/wecom-bot-mcp-server).

## License

[Apache-2.0](LICENSE)
