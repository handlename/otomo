# otomo

otomo is a Slack bot powered by gen-AI.

## Abilities

- Respond to mentions
- Summarize thread
- Search the web for information using Tavily Search API
- Fetch and read web page content as Markdown or plain text

## Documentation

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Core architectural principles (e.g. Always Valid Domain Model, Clean Architecture).
- [DESIGN.md](./DESIGN.md) - Technical design, configuration, and technology stack.
- [GLOSSARY.md](./GLOSSARY.md) - Glossary of domain concepts and terminology.

## Installation

TBD


## Local Development with Tracing

To run the application with distributed tracing enabled locally:

1. Start the Jaeger container using Docker Compose:
   ```bash
   docker compose up -d
   ```
2. Ensure `config.toml` has tracing enabled:
   ```toml
   [otel]
   enabled = true
   exporter = "otlp"
   service_name = "otomo"
   ```
3. Run the application with the `OTEL_EXPORTER_OTLP_INSECURE` environment variable set to `true`:
   ```bash
   OTEL_EXPORTER_OTLP_INSECURE=true go run ./cmd/otomo/main.go server
   ```
   Alternatively, you can set `OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"`.
4. Open the Jaeger UI in your browser at http://localhost:16686 to view and analyze your traces.

## License

MIT

## Author

@handlename
