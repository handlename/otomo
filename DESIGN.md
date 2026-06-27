# Otomo - Design Document

## Overview

Otomo is a Slack bot powered by generative AI that provides intelligent responses to user interactions. The bot can respond to mentions, summarize threads, and perform various AI-powered tasks within Slack channels.

## Architecture

For the core architectural principles, layer structure, and patterns, please refer to [ARCHITECTURE.md](./ARCHITECTURE.md).

## Technology Stack

- **Language**: Go 1.24
- **AI Provider**: AWS Bedrock (Claude models)
- **Chat Platform**: Slack API
- **Distributed Tracing**: OpenTelemetry Go SDK (Stdout and OTLP HTTP exporters)
- **Deployment**: AWS Lambda + Function URLs
- **Infrastructure**: Terraform
- **HTTP Framework**: chi router with Ridge (AWS Lambda HTTP adapter)
- **Configuration**: TOML-based with environment variable templating
- **Logging**: zerolog
- **Testing**: Go standard testing + testify

## Directory Structure

```
otomo/
├── .github/                    # GitHub workflows and configuration
│   └── workflows/
├── cli/                        # Command-line interface
│   ├── cli.go                 # Main CLI entry point
│   └── command/               # CLI subcommands
├── cmd/                        # Application entry points
│   └── otomo/
│       └── main.go            # Main application entry
├── config/                     # Configuration management
├── internal/                   # Internal application code
│   ├── app/                   # Application layer
│   │   ├── service/           # Application services (event_publisher.go, messenger.go)
│   │   └── usecase/           # Use case implementations (tool_loop.go)
│   ├── domain/                # Domain layer (business logic)
│   │   ├── core/              # System-wide core concepts (Prompt, Event, Message)
│   │   ├── reasoning/         # Reasoning and AI model context (Brain, Context, Answer)
│   │   └── chat/              # Chat conversation context (Otomo, Thread, Reply)
│   ├── errorcode/             # Error code definitions
│   ├── infra/                 # Infrastructure layer
│   │   ├── app/               # Main application container logic
│   │   ├── brain/             # AI brain implementations
│   │   ├── repository/        # Repository implementations
│   │   ├── service/           # External service integrations
│   │   ├── tool/              # AI-powered tools (web_fetch.go, web_search.go)
│   │   └── ui/                # User interface (HTTP handlers, CLI terminal interface, MCP server)
│   │       ├── http/
│   │       │   ├── middleware/
│   │       │   └── slack/
│   │       ├── mcp/           # MCP server implementation
│   │       └── terminal/
│   └── testutil/              # Testing utilities
├── lambda/                     # AWS Lambda deployment files
│   ├── function.jsonnet       # Lambda function configuration
│   └── bootstrap              # Lambda bootstrap binary
├── terraform/                  # Infrastructure as Code
│   ├── main.tf               # Main Terraform configuration
│   ├── lambda.tf             # Lambda-specific resources
│   ├── iam.tf                # IAM roles and policies
│   └── modules/              # Terraform modules
├── tools/                      # Development and code generation tools
│   └── gen-vo/                # AST-based Value Object code generator
├── config.toml               # Configuration template
├── go.mod                    # Go module dependencies
├── logger.go                 # Logging configuration
├── version.go                # Version information
├── .mise.toml                # Tool, build, and deployment automation
└── README.md                 # Project documentation
```

## Core Components

For the details of each layer, packages, and their responsibilities, please refer to [ARCHITECTURE.md#core-components--layer-responsibilities](./ARCHITECTURE.md#core-components--layer-responsibilities).


## Configuration

The application uses TOML configuration with environment variable templating:

```toml
port = 8080

[slack]
signing_secret = "{{ must_env `SLACK_SIGNING_SECRET` }}"
bot_user_id = "@U08K30DRHRP"
bot_token = "{{ must_env `SLACK_BOT_TOKEN` }}"
app_token = "{{ must_env `SLACK_APP_TOKEN` }}"

[slack.error_feedback]
enable_reaction = true
reaction_emoji = "warning"
enable_post_snippet = false

[llm]
model_type = "claude"
model_id = "{{ must_env `BEDROCK_MODEL_ID` }}"

[tool]
[tool.web_search]
tavily_api_key = "{{ env `TAVILY_API_KEY` }}"

[tool.web_fetch]
whitelist_patterns = [
  '^https://.*'
]

[mcp]
port = 8000

[otel]
enabled = false
exporter = "otlp"
service_name = "otomo"
```

### Configuration Details

- **`[tool.web_search]`**:
  - `tavily_api_key`: API key for Tavily Search API.
- **`[tool.web_fetch]`**:
  - `whitelist_patterns`: Regular expression patterns specifying allowed URLs to fetch.
- **`[mcp]`**:
  - `port`: Port for the local MCP server to listen on.
- **`[otel]`**:
  - `enabled`: Set to `true` to enable distributed tracing.
  - `exporter`: The trace exporter to use. Options are `"stdout"` (for pretty-printed JSON output to standard output) or `"otlp"` (for sending OTLP traces over HTTP).
  - `service_name`: The service name registered in trace resource metadata. Default is `"otomo"`.

## Deployment Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│                 │    │                 │    │                 │
│   Slack API     │───▶│  AWS Lambda     │───▶│  AWS Bedrock    │
│                 │    │  (Function URL) │    │  (Claude AI)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │                 │
                       │  Lambda Logs    │
                       │  (CloudWatch)   │
                       └─────────────────┘
```

### Lambda Deployment
- Built as a single binary using Go's AWS Lambda runtime
- Deployed via Terraform with Function URLs for HTTP access
- Uses Ridge framework for HTTP request handling in Lambda environment

### Infrastructure Management
- Terraform for infrastructure provisioning
- Automated builds with GoReleaser
- GitHub Actions for CI/CD

## Key Features

1. **Slack Integration**
   - Event-driven architecture for Slack webhooks
   - Message processing and response generation
   - Thread summarization capabilities

2. **AI-Powered Responses**
   - Integration with AWS Bedrock (Claude models)
   - Context-aware conversation handling
   - Customizable AI brain implementations

3. **Event-Driven Architecture**
   - Domain events for loose coupling
   - Event publishing and subscription patterns
   - Asynchronous processing capabilities

4. **Clean Architecture**
   - Clear separation of concerns
   - Dependency inversion
   - Testable components with mock implementations

5. **Production Ready**
   - Comprehensive error handling
   - Structured logging
   - Health checks and monitoring
   - Retry mechanisms for external services

6. **External Tool Integrations**
   - Web search capability via Tavily Search API
   - Secure web page fetching with URL whitelist validation patterns

7. **Distributed Tracing**
   - End-to-end execution flow tracing using OpenTelemetry (OTel)
   - Parent spans started at Slack HTTP webhooks, Terminal TUI chat, and MCP server requests
   - Sub-spans instrumented for use cases, Bedrock model invocations, and Slack API services
   - Decoupled design keeping Domain layer entities entirely clean of OTel SDK/API dependencies

## Testing Strategy

- Unit tests for domain logic
- Integration tests for external service interactions
- Mock implementations for testing isolation
- Test utilities for common testing patterns

## Build and Development

- **Build and Task Automation**: Task automation and tool management via `mise`
- **Dependency Management**: Go modules
- **Code Generation**: AST-based code generator for value objects (`go generate ./...`)
- **Release Management**: GoReleaser for automated releases
- **Code Quality**: Static analysis and linting integration